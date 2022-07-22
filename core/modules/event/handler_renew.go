package event

import (
	"github.com/hamster-shared/hamster-provider/log"
)

type RenewVmHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *RenewVmHandler) HandlerEvent(e *VmRequest) {

	orderNo := e.OrderNo

	cm := h.CoreContext.Cm
	cfg, _ := cm.GetConfig()

	cfg.ChainRegInfo.RenewOrderIndex = orderNo
	_ = cm.Save(cfg)
	overdue := h.CoreContext.ReportClient.CalculateInstanceOverdue(e.AgreementNo)
	timer := h.CoreContext.TimerService.GetTimer(e.AgreementNo)
	timer.Reset(overdue)
	err := h.CoreContext.ReportClient.OrderExec(orderNo)
	if err != nil {
		log.GetLogger().Error("report order exec fail")
	}
}

func (h *RenewVmHandler) Name() string {
	return ResourceOrder_ReNewOrderSuccess
}
