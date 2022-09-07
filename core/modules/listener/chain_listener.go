package listener

import (
	"context"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain2 "github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/hamster-shared/hamster-provider/log"
	"time"
)

type ChainListener struct {
	eventService event.IEventService
	api          *gsrpc.SubstrateAPI
	cm           *config.ConfigManager
	reportClient chain2.ReportClient
	cancel       func()
	ctx          context.Context
}

func NewChainListener(eventService event.IEventService, cm *config.ConfigManager) *ChainListener {
	return &ChainListener{
		eventService: eventService,
		cm:           cm,
	}
}

func (l *ChainListener) SetChainApi(api *gsrpc.SubstrateAPI, reportClient chain2.ReportClient) {
	l.api = api
	l.reportClient = reportClient
}

func (l *ChainListener) GetState() bool {
	return l.cancel != nil
}

func (l *ChainListener) SetState(option bool) error {
	if option {
		return l.start()
	} else {
		return l.stop()
	}
}

func (l *ChainListener) start() error {
	if l.cancel != nil {
		l.cancel()
	}

	cfg, err := l.cm.GetConfig()
	if err != nil {
		return err
	}

	resource := chain2.ResourceInfo{
		PeerId:        cfg.Identity.PeerID,
		Cpu:           cfg.Vm.Cpu,
		Memory:        cfg.Vm.Mem,
		System:        cfg.Vm.System,
		CpuModel:      utils.GetCpuModel(),
		Price:         cfg.ChainRegInfo.Price,
		ExpireTime:    time.Now().AddDate(0, 0, 10),
		ResourceIndex: cfg.ChainRegInfo.ResourceIndex,
	}
	err = l.reportClient.RegisterResource(resource)

	if err != nil {
		return err
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())
	isPanic := make(chan bool)
	go l.setWatchEventState(l.ctx, isPanic)
	return nil
}

func (l *ChainListener) setWatchEventState(ctx context.Context, isPanic chan bool) {
	for {
		go l.watchEvent(ctx, isPanic)
		select {
		case <-isPanic:
			go l.watchEvent(ctx, isPanic)
		}
	}
}

func (l *ChainListener) stop() error {
	if l.cancel != nil {
		l.cancel()
		l.cancel = nil
	}
	cfg, err := l.cm.GetConfig()
	if err != nil {
		return err
	}
	thegraph.SetIsServer(false)
	return l.reportClient.RemoveResource(cfg.ChainRegInfo.ResourceIndex)
}

// WatchEvent chain event listener
func (l *ChainListener) watchEvent(ctx context.Context, channel chan bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("watchEventError:", err)
			channel <- true
		}
	}()

	meta, err := l.api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		panic(err)
	}

	sub, err := l.api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case set := <-sub.Chan():
			log.GetLogger().Info("watch ：", set.Block.Hex())
			for _, chng := range set.Changes {
				if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
					// skip, we are only interested in events with content
					continue
				}
				// Decode the event records
				evt := chain2.MyEventRecords{}
				storageData := chng.StorageData
				meta, err := l.api.RPC.State.GetMetadataLatest()
				err = types.EventRecordsRaw(storageData).DecodeEventRecords(meta, &evt)
				if err != nil {
					log.GetLogger().Error(err)
					continue
				}
				for _, e := range evt.ResourceOrder_CreateOrderSuccess {
					// order successfully created
					l.dealCreateOrderSuccess(e)
				}

				for _, e := range evt.ResourceOrder_ReNewOrderSuccess {
					// order renewal successful
					l.dealReNewOrderSuccess(e)
				}

				for _, e := range evt.ResourceOrder_WithdrawLockedOrderPriceSuccess {
					// order cancelled successfully
					log.GetLogger().Info("deal ResourceOrder_WithdrawLockedOrderPriceSuccess")
					l.dealCancelOrderSuccess(e)
				}

			}
		}
	}

}

func (l *ChainListener) dealCreateOrderSuccess(e chain2.EventResourceOrderCreateOrderSuccess) {
	cfg, err := l.cm.GetConfig()
	if err != nil {
		panic(err)
	}
	log.GetLogger().Infof("\tResourceOrder:CreateOrderSuccess:: (phase=%#v)\n", e.Phase)

	if e.ResourceIndex == types.NewU64(cfg.ChainRegInfo.ResourceIndex) {
		// process the order
		log.GetLogger().Info("deal order: ", e.OrderIndex)
		// record the id of the processed order
		cfg.ChainRegInfo.OrderIndex = uint64(e.OrderIndex)
		cfg.ChainRegInfo.AccountAddress = utils.AccountIdToAddress(e.AccountId)
		_ = l.cm.Save(cfg)
		evt := &event.VmRequest{
			Tag:       event.OPCreatedVm,
			Cpu:       cfg.Vm.Cpu,
			Mem:       cfg.Vm.Mem,
			Disk:      cfg.Vm.Disk,
			OrderNo:   uint64(e.OrderIndex),
			System:    cfg.Vm.System,
			PublicKey: e.PublicKey,
			Image:     cfg.Vm.Image,
		}
		l.eventService.Create(evt)

	} else {
		log.GetLogger().Warn("resourceIndex is not equals ")
	}
}

func (l *ChainListener) dealReNewOrderSuccess(e chain2.EventResourceOrderReNewOrderSuccess) {
	cfg, err := l.cm.GetConfig()
	if err != nil {
		panic(err)
	}
	if e.ResourceIndex == types.NewU64(cfg.ChainRegInfo.ResourceIndex) {
		evt := &event.VmRequest{
			Tag:     event.OPRenewVM,
			OrderNo: uint64(e.OrderIndex),
		}
		l.eventService.Renew(evt)
	}
}

func (l *ChainListener) dealCancelOrderSuccess(e chain2.EventResourceOrderWithdrawLockedOrderPriceSuccess) {
	cfg, err := l.cm.GetConfig()
	if err != nil {
		panic(err)
	}
	if e.OrderIndex == types.NewU64(cfg.ChainRegInfo.OrderIndex) {
		evt := &event.VmRequest{
			Tag:     event.OPDestroyVm,
			Cpu:     cfg.Vm.Cpu,
			Mem:     cfg.Vm.Mem,
			Disk:    cfg.Vm.Disk,
			OrderNo: uint64(e.OrderIndex),
			System:  cfg.Vm.System,
			Image:   cfg.Vm.Image,
		}
		l.eventService.Destroy(evt)
	}
}
