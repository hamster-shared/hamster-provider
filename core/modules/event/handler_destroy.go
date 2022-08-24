package event

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/log"
)

type DestroyVmHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *DestroyVmHandler) HandlerEvent(e *VmRequest) {
	log.GetLogger().Info("handler destory order :", e.OrderNo)
	orderNo := e.OrderNo
	agreementNo, err := h.CoreContext.ReportClient.GetAgreementIndex(orderNo)
	if err != nil {
		log.GetLogger().Error("query agreementNo fail")
	}

	_ = thegraph.Uninstall()
	thegraph.SetIsServer(false)
	h.CoreContext.TimerService.UnSubTimer(agreementNo)
	h.CoreContext.TimerService.UnSubTicker(agreementNo)
	cfg := h.CoreContext.GetConfig()
	cfg.ChainRegInfo.AccountAddress = ""
	_ = h.CoreContext.Cm.Save(cfg)
}

func (h *DestroyVmHandler) Name() string {
	return ResourceOrder_WithdrawLockedOrderPriceSuccess
}
