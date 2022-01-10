package test

import (
	"github.com/hamster-shared/hamster-provider/cmd"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"testing"
)

func TestDestroyVm(t *testing.T) {
	context := cmd.NewContext()
	eventService := event.NewEventService(context)
	cfg := context.GetConfig()

	orderNo := 0

	req := &event.VmRequest{
		Tag:     event.OPDestroyVm,
		Cpu:     cfg.Vm.Cpu,
		Mem:     cfg.Vm.Mem,
		Disk:    cfg.Vm.Disk,
		OrderNo: uint64(orderNo),
		System:  cfg.Vm.System,
		Image:   cfg.Vm.Image,
	}
	eventService.Destroy(req)
}
