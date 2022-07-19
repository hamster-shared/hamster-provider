package event

import (
	"github.com/hamster-shared/hamster-provider/log"
)

type CreateVmHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *CreateVmHandler) HandlerEvent(e *VmRequest) {

	// inject public key
	_, err := h.CoreContext.VmManager.CreateAndStartAndInjectionPublicKey(e.getName(), e.PublicKey)
	if err != nil {
		log.GetLogger().Errorf("failed to process order,%v", err)
		return
	}
	err = successDealOrder(h.CoreContext, e.OrderNo, e.getName())
	if err != nil {
		log.GetLogger().Error("failed to process order")
	} else {
		log.GetLogger().Info("processing order complete")
	}

	// notify vm is ready
	err = h.CoreContext.ReportClient.OrderExec(e.OrderNo)
	if err != nil {
		log.GetLogger().Errorf("failed to process order,%v", err)
	}

}

func (h *CreateVmHandler) Name() string {
	return ResourceOrder_CreateOrderSuccess
}
