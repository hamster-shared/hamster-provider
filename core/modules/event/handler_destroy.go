package event

import (
	"github.com/hamster-shared/hamster-provider/core/context"
	log "github.com/sirupsen/logrus"
)

type DestroyVmHandler struct {
	AbstractHandler
	CoreContext context.CoreContext
}

func (h *DestroyVmHandler) HandlerEvent(e *VmRequest) {
	orderNo := e.OrderNo
	agreementNo, err := h.CoreContext.ReportClient.GetAgreementIndex(orderNo)
	if err != nil {
		log.Error("query agreementNo fail")
	}

	_ = h.CoreContext.VmManager.Stop(e.getName())
	_ = h.CoreContext.VmManager.Destroy(e.getName())
	h.CoreContext.TimerService.UnSubTimer(agreementNo)
	h.CoreContext.TimerService.UnSubTicker(agreementNo)
}

func (h *DestroyVmHandler) Name() string {
	return ResourceOrder_WithdrawLockedOrderPriceSuccess
}
