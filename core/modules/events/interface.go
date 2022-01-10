package events

import "github.com/hamster-shared/hamster-provider/core/modules/vm"

type EventInterface interface {

	// Hook event triggered execution
	Hook() error
	// SetVmManager set up vm management
	SetVmManager(manager vm.Manager)

	// AddCompleteCallback  add event completion callback
	AddCompleteCallback(func())
}
