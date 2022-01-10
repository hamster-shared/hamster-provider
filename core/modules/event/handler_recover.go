package event

import (
	"github.com/hamster-shared/hamster-provider/core/context"
	log "github.com/sirupsen/logrus"
)

type RecoverVmHandler struct {
	AbstractHandler
	CoreContext context.CoreContext
}

func (h *RecoverVmHandler) HandlerEvent(e *VmRequest) {

	// Determine the status of the virtual machine and restart if the virtual machine is stopped

	vmManager := h.CoreContext.VmManager

	status, err := vmManager.Status(e.getName())
	if err != nil {
		log.Errorf("Order %s failed to restore, reason: VM instance does not exist", e.getName())
		return
	}

	if !status.IsRunning() {
		err = vmManager.Start(e.getName())
		if err != nil {
			log.Errorf("Order %s failed to restore, reason: VM failed to start", e.getName())
			return
		}
	}

	err = successDealOrder(h.CoreContext, e.OrderNo, e.getName())

	if err != nil {
		log.Error("handling recovery failures")
	} else {
		log.Info("handling order complete")
	}

}

func (h *RecoverVmHandler) Name() string {
	return ResourceOrder_Recover
}
