package event

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type CreateVmHandler struct {
	AbstractHandler
	CoreContext EventContext
}

func (h *CreateVmHandler) HandlerEvent(e *VmRequest) {

	// inject public key
	_, err := h.CoreContext.VmManager.CreateAndStartAndInjectionPublicKey(e.getName(), e.PublicKey)
	if err != nil {
		log.Error("failed to process order,%v", err)
		return
	}
	err = successDealOrder(h.CoreContext, e.OrderNo, e.getName())
	if err != nil {
		log.Error("failed to process order")
	} else {
		log.Info("processing order complete")
	}

	// notify vm is ready
	err = h.CoreContext.ReportClient.OrderExec(e.OrderNo)
	if err != nil {
		fmt.Println(err)
		log.Error("failed to process order,%v", err)
	}

}

func (h *CreateVmHandler) Name() string {
	return ResourceOrder_CreateOrderSuccess
}
