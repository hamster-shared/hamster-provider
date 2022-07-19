package event

import (
	"github.com/hamster-shared/hamster-provider/log"
)

type DestroyVmHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *DestroyVmHandler) HandlerEvent(e *VmRequest) {
	orderNo := e.OrderNo
	agreementNo, err := h.CoreContext.ReportClient.GetAgreementIndex(orderNo)
	if err != nil {
		log.GetLogger().Error("query agreementNo fail")
	}

	_ = h.CoreContext.VmManager.Stop(e.getName())
	_ = h.CoreContext.VmManager.Destroy(e.getName())
	h.CoreContext.TimerService.UnSubTimer(agreementNo)
	h.CoreContext.TimerService.UnSubTicker(agreementNo)
	targetAddress := getVmTargetAddress(h.CoreContext, e.getName())
	_, _ = h.CoreContext.P2pClient.Close(targetAddress)
}

func (h *DestroyVmHandler) Name() string {
	return ResourceOrder_WithdrawLockedOrderPriceSuccess
}
