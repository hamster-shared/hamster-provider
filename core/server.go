package core

import (
	"fmt"
	context2 "github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/corehttp"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"os"
)

type Server struct {
	ctx          context2.CoreContext
	eventService event.IEventService
}

func NewServer(ctx context2.CoreContext) *Server {
	return &Server{
		ctx:          ctx,
		eventService: ctx.EventService,
	}
}

func (s *Server) Run() {

	//cfg := s.ctx.GetConfig()
	//
	//// 1: start api
	//
	//peerId := cfg.Identity.PeerID
	//if peerId == "" {
	//	logrus.Error("identity was not set in config (was 'ttcp init' run?)")
	//	os.Exit(0)
	//	return
	//}
	//if len(peerId) == 0 {
	//	logrus.Error("no peer ID in config! (was 'ipfs init' run?)")
	//	return
	//}
	//
	//cpuModel := utils.GetCpuModel()

	// determine if registration is required
	//if cfg.ConfigFlag == config.DONE {
	//
	//	resource := chain2.ResourceInfo{
	//		PeerId:     peerId,
	//		Cpu:        cfg.Vm.Cpu,
	//		Memory:     cfg.Vm.Mem,
	//		System:     cfg.Vm.System,
	//		CpuModel:   cpuModel,
	//		Price:      1000,
	//		ExpireTime: time.Now().AddDate(0, 0, 10),
	//	}
	//	// 2: blockchain registration
	//	for {
	//		err := s.ctx.ReportClient.RegisterResource(resource)
	//		if err != nil {
	//			logrus.Errorf("Blockchain registration failed, the reason for the failure： %s", err.Error())
	//			time.Sleep(time.Second * 30)
	//		} else {
	//			break
	//		}
	//	}
	//
	//	status := "running"
	//
	//	// TODO ... init working vm
	//	if cfg.ChainRegInfo.Working == status {
	//		// calculate instance expiration time
	//		duration := s.ctx.ReportClient.CalculateInstanceOverdue(cfg.ChainRegInfo.OrderIndex)
	//
	//		// not expired
	//		if duration > 0 {
	//			// rebuild p2p link
	//			req := &event.VmRequest{
	//				Tag:     event.OPRecoverVM,
	//				Cpu:     cfg.Vm.Cpu,
	//				Mem:     cfg.Vm.Mem,
	//				Disk:    cfg.Vm.Disk,
	//				OrderNo: uint64(cfg.ChainRegInfo.OrderIndex),
	//				System:  cfg.Vm.System,
	//				Image:   cfg.Vm.Image,
	//			}
	//			s.eventService.Recover(req)
	//		}
	//	}
	//
	//	// set state
	//	cfg = s.ctx.GetConfig()
	//	cfg.ChainRegInfo.Working = status
	//
	//	_ = s.ctx.Cm.Save(cfg)
	//
	//	// event listener
	//	//go s.WatchEvent()
	//}

	err := s.ctx.P2pClient.Listen("/x/provider", fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", s.ctx.GetConfig().ApiPort))
	if err != nil {
		os.Exit(1)
	}
	err = corehttp.StartApi(&s.ctx)

	if err != nil {
		os.Exit(1)
	}
}
