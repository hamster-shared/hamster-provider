package events

import "github.com/hamster-shared/hamster-provider/core/modules/vm"

type RenewVM struct {
	completeCallback []func()
	vmManager        vm.Manager
	OrderIndex       uint64
}

// Hook event triggered execution
func (e *RenewVM) Hook() error {

	return nil
}

func (e *RenewVM) SetVmManager(manager vm.Manager) {
	e.vmManager = manager
}

func (e *RenewVM) AddCompleteCallback(callback func()) {
	if e.completeCallback == nil {
		e.completeCallback = make([]func(), 0)
	}
	e.completeCallback = append(e.completeCallback, callback)
}
