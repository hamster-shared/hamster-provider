package event

import (
	"github.com/hamster-shared/hamster-provider/core/modules/factory"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/log"
)

type DestroyVmHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *DestroyVmHandler) HandlerEvent(e *VmRequest) {
	log.GetLogger().Info("handler destroy order :", e.OrderNo)
	orderNo := e.OrderNo
	agreementNo, err := h.CoreContext.ReportClient.GetAgreementIndex(orderNo)
	if err != nil {
		log.GetLogger().Error("query agreementNo fail")
	} else {
		agreementNo = h.CoreContext.GetConfig().ChainRegInfo.AgreementIndex
	}

	//_ = thegraph.Uninstall()
	deployType := h.CoreContext.GetConfig().ChainRegInfo.DeployType
	chain, err := factory.GetChain(deployType)
	if err == nil {
		_ = chain.Stop()
	}
	thegraph.SetIsServer(false)
	h.CoreContext.TimerService.UnSubTimer(agreementNo)
	h.CoreContext.TimerService.UnSubTicker(agreementNo)
	cfg := h.CoreContext.GetConfig()
	cfg.ChainRegInfo.AccountAddress = ""
	cfg.ChainRegInfo.DeployType = 0
	_ = h.CoreContext.Cm.Save(cfg)
}

func (h *DestroyVmHandler) Name() string {
	return ResourceOrder_WithdrawLockedOrderPriceSuccess
}
