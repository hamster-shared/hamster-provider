package events

import (
	"github.com/hamster-shared/hamster-provider/core/modules/vm"
	"github.com/sirupsen/logrus"
)

type StartVm struct {
	Cpu       uint64
	Memory    uint64
	Disk      uint64
	Name      string
	System    string
	PublicKey string
	Image     string

	completeCallback []func()
	vmManager        vm.Manager
}

func (e *StartVm) Hook() error {
	template := vm.Template{
		Cpu:       e.Cpu,
		Memory:    e.Memory,
		Dist:      e.Disk,
		Name:      e.Name,
		System:    e.System,
		PublicKey: e.PublicKey,
		Image:     e.Image,
	}
	e.vmManager.SetTemplate(template)

	// injection public key
	err := e.vmManager.CreateAndStartAndInjectionPublicKey(e.PublicKey)
	if err != nil {
		logrus.Error("failed to process order,%v", err)
		return err
	}

	logrus.Info("processing order complete")

	// The resource has been processed and a callback is sent
	if e.completeCallback != nil {
		for _, cb := range e.completeCallback {
			cb()
		}
	}
	return nil
}

func (e *StartVm) SetVmManager(manager vm.Manager) {
	e.vmManager = manager
}

func (e *StartVm) AddCompleteCallback(callback func()) {
	if e.completeCallback == nil {
		e.completeCallback = make([]func(), 0)
	}
	e.completeCallback = append(e.completeCallback, callback)
}
