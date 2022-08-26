package test

import (
	"github.com/hamster-shared/hamster-provider/cmd"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"testing"
)

func TestCreateVm(t *testing.T) {
	context := cmd.NewContext()
	eventContex := event.EventContext{
		ReportClient: context.ReportClient,
		VmManager:    context.VmManager,
		TimerService: context.TimerService,
		Cm:           context.Cm,
		P2pClient:    context.P2pClient,
	}
	eventService := event.NewEventService(&eventContex)
	cfg := context.GetConfig()

	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCWCDS4Io8+PFqGepqy0YNrtw3B7g7lhg7WNcH2VyJmmHlvft69N3S4EzDugEUDgPbgihiL56wyq56GtOG6+RuRuqkEU983MRC6j0yazem/KPs2nAS0NW5A8Nzxm9ixXnF9Bw6qHpO+L8ZbKdsIR+xux5QVriWTmDd/FeaovzRa/Ogr/BdShsp5H1s8aKkj2ygm16rlWAuQcoPQJPDWJVLM9cub8wj/AGrOzRDQCMnbcm69BZT7GPbodVmBIlugICuSVVvKSpZEa0QHCdQW2z2kIan7EwEI7LPYyDpCRAAI2mYEsl9WIIzae1ACK7dKwp9DKfLlKU4YRNfvR5stGgNezelz2pbN0TvK0T6NrqlKDo1eZQbHzRzvUKtDCiwSBdauJVus5Zowqy8lXr9wVosbra8z8cd+vM+e5+82fEjnE3BQm6NUHatOfxe/1MtYeem1Zlru5ISc25ceCXJd/l6qUrlIamHKgMpxvAt8g9pcPpPH2YozLkRlohcdWrhA+kk= gr@gr-Lenovo"
	orderNo := 0

	req := &event.VmRequest{
		Tag:       event.OPCreatedVm,
		Cpu:       cfg.Vm.Cpu,
		Mem:       cfg.Vm.Mem,
		Disk:      cfg.Vm.Disk,
		OrderNo:   uint64(orderNo),
		System:    cfg.Vm.System,
		PublicKey: publicKey,
		Image:     cfg.Vm.Image,
	}
	eventService.Create(req)

}
