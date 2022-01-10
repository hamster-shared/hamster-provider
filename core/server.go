package core

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/corehttp"
	chain2 "github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Server struct {
	ctx           context2.CoreContext
	instanceTimer *time.Timer
	eventService  event.IEventService
}

func NewServer(ctx context2.CoreContext) *Server {
	return &Server{
		ctx:          ctx,
		eventService: event.NewEventService(ctx),
	}
}

func (s *Server) Run() {

	cfg := s.ctx.GetConfig()

	// 1: start api

	peerId := cfg.Identity.PeerID
	if peerId == "" {
		logrus.Error("identity was not set in config (was 'ttcp init' run?)")
		os.Exit(0)
		return
	}
	if len(peerId) == 0 {
		logrus.Error("no peer ID in config! (was 'ipfs init' run?)")
		return
	}

	cpuModel := utils.GetCpuModel()

	// determine if registration is required
	if cfg.ConfigFlag == config.DONE {

		resource := chain2.ResourceInfo{
			PeerId:     peerId,
			Cpu:        cfg.Vm.Cpu,
			Memory:     cfg.Vm.Mem,
			System:     cfg.Vm.System,
			CpuModel:   cpuModel,
			Price:      1000,
			ExpireTime: time.Now().AddDate(0, 0, 10),
		}
		// 2: blockchain registration
		for {
			err := s.ctx.ReportClient.RegisterResource(resource)
			if err != nil {
				logrus.Errorf("Blockchain registration failed, the reason for the failure： %s", err.Error())
				time.Sleep(time.Second * 30)
			} else {
				break
			}
		}

		status := "running"

		// TODO ... init working vm
		if cfg.ChainRegInfo.Working == status {
			// calculate instance expiration time
			duration := s.ctx.ReportClient.CalculateInstanceOverdue(cfg.ChainRegInfo.OrderIndex)

			// not expired
			if duration > 0 {
				// rebuild p2p link
				req := &event.VmRequest{
					Tag:     event.OPRecoverVM,
					Cpu:     cfg.Vm.Cpu,
					Mem:     cfg.Vm.Mem,
					Disk:    cfg.Vm.Disk,
					OrderNo: uint64(cfg.ChainRegInfo.OrderIndex),
					System:  cfg.Vm.System,
					Image:   cfg.Vm.Image,
				}
				s.eventService.Recover(req)
			}
		}

		// set state
		cfg = s.ctx.GetConfig()
		cfg.ChainRegInfo.Working = status

		_ = s.ctx.Cm.Save(cfg)

		// event listener
		go s.WatchEvent()
	}

	err := corehttp.StartApi(&s.ctx)
	if err != nil {
		os.Exit(1)
	}
}

// WatchEvent chain event listener
func (s *Server) WatchEvent() {

	api := s.ctx.SubstrateApi

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}
			// Decode the event records
			event := chain2.MyEventRecords{}
			storageData := chng.StorageData
			meta, err := api.RPC.State.GetMetadataLatest()
			err = types.EventRecordsRaw(storageData).DecodeEventRecords(meta, &event)
			if err != nil {
				fmt.Println(err)
				logrus.Error(err)
				continue
			}
			for _, e := range event.ResourceOrder_CreateOrderSuccess {
				// order successfully created
				s.dealCreateOrderSuccess(e)
			}

			for _, e := range event.ResourceOrder_ReNewOrderSuccess {
				// order renewal successful
				s.dealReNewOrderSuccess(e)
			}

			for _, e := range event.ResourceOrder_WithdrawLockedOrderPriceSuccess {
				// order cancelled successfully
				s.dealCancelOrderSuccess(e)
			}
		}
	}
}

func (s *Server) dealCreateOrderSuccess(e chain2.EventResourceOrderCreateOrderSuccess) {
	cm := s.ctx.Cm
	cfg, err := cm.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\tResourceOrder:CreateOrderSuccess:: (phase=%#v)\n", e.Phase)

	if e.ResourceIndex == types.NewU64(cfg.ChainRegInfo.ResourceIndex) {
		// process the order
		fmt.Println("deal order", e.OrderIndex)
		// record the id of the processed order
		cfg.ChainRegInfo.OrderIndex = uint64(e.OrderIndex)
		_ = cm.Save(cfg)
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
		s.eventService.Create(evt)

	} else {
		fmt.Println("resourceIndex is not equals ")
	}
}

func (s *Server) dealReNewOrderSuccess(e chain2.EventResourceOrderReNewOrderSuccess) {
	cm := s.ctx.Cm
	cfg, err := cm.GetConfig()
	if err != nil {
		panic(err)
	}
	if e.ResourceIndex == types.NewU64(cfg.ChainRegInfo.ResourceIndex) {
		evt := &event.VmRequest{
			Tag:     event.OPRenewVM,
			OrderNo: uint64(e.OrderIndex),
		}
		s.eventService.Renew(evt)
	}
}

func (s *Server) dealCancelOrderSuccess(e chain2.EventResourceOrderWithdrawLockedOrderPriceSuccess) {
	cm := s.ctx.Cm
	cfg, err := cm.GetConfig()
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
		s.eventService.Destroy(evt)
	}
}
