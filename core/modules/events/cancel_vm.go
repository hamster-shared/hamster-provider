package events

import "github.com/hamster-shared/hamster-provider/core/modules/vm"

type CancelVM struct {
	completeCallback []func()
	vmManager        vm.Manager
	OrderIndex       uint64
}

// Hook event triggered execution
func (e *CancelVM) Hook() error {

	_ = e.vmManager.Shutdown()
	_ = e.vmManager.Destroy()

	return nil
}

func (e *CancelVM) SetVmManager(manager vm.Manager) {
	e.vmManager = manager
}

func (e *CancelVM) AddCompleteCallback(callback func()) {
	if e.completeCallback == nil {
		e.completeCallback = make([]func(), 0)
	}
	e.completeCallback = append(e.completeCallback, callback)
}
