package core

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/corehttp"
	chain2 "github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/events"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	vm2 "github.com/hamster-shared/hamster-provider/core/modules/vm"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Server struct {
	ctx           context2.CoreContext
	instanceTimer *time.Timer
}

func NewServer(ctx context2.CoreContext) *Server {
	return &Server{
		ctx: ctx,
	}
}

func (s *Server) Run() {

	cfg := s.ctx.GetConfig()

	// 1: start api
	go func() {
		err := corehttp.StartApi(&s.ctx)
		if err != nil {
			os.Exit(1)
		}
	}()

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

	if s.ctx.GetConfig().ChainRegInfo.OrderIndex > 0 {
		s.ctx.VmManager.SetTemplate(vm2.Template{
			Cpu:    cfg.Vm.Cpu,
			Memory: cfg.Vm.Mem,
			Name:   "order_" + strconv.FormatUint(uint64(cfg.ChainRegInfo.OrderIndex), 10),
			System: cfg.Vm.System,
			Image:  cfg.Vm.Image,
		})

		// calculate instance expiration time
		isOverdue := s.dealOverdueOrder()

		// not expired
		if !isOverdue {
			// rebuild p2p link
			err := s.forwardSSHToP2p()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// eventListener
	s.WatchEvent()

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

func (s *Server) successDealOrder() {
	err := s.forwardSSHToP2p()
	if err != nil {
		fmt.Println(err)
		return
	}
	// notify vm is ready
	err = s.ctx.ReportClient.OrderExec(s.ctx.GetConfig().ChainRegInfo.OrderIndex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// report heartbeat
	agreementIndex := s.ctx.GetConfig().ChainRegInfo.AgreementIndex
	_ = s.ctx.ReportClient.Heartbeat(agreementIndex)

	// send timed heartbeats
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for {
			<-ticker.C
			// report heartbeat
			agreementIndex := s.ctx.GetConfig().ChainRegInfo.AgreementIndex
			_ = s.ctx.ReportClient.Heartbeat(agreementIndex)
		}
	}()

	s.dealOverdueOrder()
}

func (s *Server) dealCreateOrderSuccess(e chain2.EventResourceOrderCreateOrderSuccess) {
	cm := s.ctx.Cm
	vmManager := s.ctx.VmManager
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
		evt := &events.StartVm{
			Cpu:       cfg.Vm.Cpu,
			Memory:    cfg.Vm.Mem,
			Disk:      cfg.Vm.Disk,
			Name:      "order_" + strconv.FormatUint(uint64(e.OrderIndex), 10),
			System:    cfg.Vm.System,
			PublicKey: e.PublicKey,
			Image:     cfg.Vm.Image,
		}
		evt.SetVmManager(vmManager)
		evt.AddCompleteCallback(s.successDealOrder)
		err = evt.Hook()
		if err != nil {
			logrus.Error(err)
		}

	} else {
		fmt.Println("resourceIndex is not equals ")
	}
}

func (s *Server) dealReNewOrderSuccess(e chain2.EventResourceOrderReNewOrderSuccess) {
	cm := s.ctx.Cm
	vmManager := s.ctx.VmManager
	cfg, err := cm.GetConfig()
	if err != nil {
		panic(err)
	}
	if e.ResourceIndex == types.NewU64(cfg.ChainRegInfo.ResourceIndex) {
		cfg.ChainRegInfo.RenewOrderIndex = uint64(e.OrderIndex)
		_ = cm.Save(cfg)
		evt := &events.RenewVM{}
		evt.SetVmManager(vmManager)
		overdue := s.ctx.ReportClient.CalculateInstanceOverdue(s.ctx.GetConfig().ChainRegInfo.AgreementIndex)
		s.instanceTimer.Reset(overdue)
		orderIndex := s.ctx.GetConfig().ChainRegInfo.RenewOrderIndex
		err = s.ctx.ReportClient.OrderExec(orderIndex)
		err = evt.Hook()
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (s *Server) dealCancelOrderSuccess(e chain2.EventResourceOrderWithdrawLockedOrderPriceSuccess) {
	cm := s.ctx.Cm
	vmManager := s.ctx.VmManager
	cfg, err := cm.GetConfig()
	if err != nil {
		panic(err)
	}
	if e.OrderIndex == types.NewU64(cfg.ChainRegInfo.OrderIndex) {
		evt := &events.CancelVM{}
		evt.SetVmManager(vmManager)
		err = evt.Hook()
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (s *Server) dealOverdueOrder() bool {
	// calculate instance expiration time
	overdue := s.ctx.ReportClient.CalculateInstanceOverdue(s.ctx.GetConfig().ChainRegInfo.OrderIndex)
	s.instanceTimer = time.NewTimer(overdue)

	go func(t *time.Timer) {
		<-t.C
		cfg := s.ctx.GetConfig()

		// expires triggers close
		_ = s.ctx.VmManager.Stop()
		_ = s.ctx.VmManager.Destroy()
		// modify the resource status on the chain to unused
		_ = s.ctx.ReportClient.ChangeResourceStatus(s.ctx.GetConfig().ChainRegInfo.ResourceIndex)
		// delete agreement number
		cfg.ChainRegInfo.OrderIndex = 0
		cfg.ChainRegInfo.AgreementIndex = 0
		cfg.ChainRegInfo.RenewOrderIndex = 0
		_ = s.ctx.Cm.Save(cfg)
	}(s.instanceTimer)

	return overdue < 0
}

func (s *Server) forwardSSHToP2p() error {
	// P2P listen port exposure
	ip, err := s.ctx.VmManager.GetIp()
	if err != nil {
		logrus.Error(err)
		return err
	}

	targetOpt := fmt.Sprintf("/ip4/%s/tcp/%d", ip, s.ctx.VmManager.GetAccessPort())
	err = s.ctx.P2pClient.Listen(targetOpt)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
