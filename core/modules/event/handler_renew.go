package event

import (
	"github.com/hamster-shared/hamster-provider/core/context"
	log "github.com/sirupsen/logrus"
)

type RenewVmHandler struct {
	AbstractHandler
	CoreContext context.CoreContext
}

func (h *RenewVmHandler) HandlerEvent(e *VmRequest) {

	orderNo := e.OrderNo

	cfg := h.CoreContext.GetConfig()
	cm := h.CoreContext.Cm

	cfg.ChainRegInfo.RenewOrderIndex = orderNo
	_ = cm.Save(cfg)
	overdue := h.CoreContext.ReportClient.CalculateInstanceOverdue(e.AgreementNo)
	timer := h.CoreContext.TimerService.GetTimer(e.AgreementNo)
	timer.Reset(overdue)
	err := h.CoreContext.ReportClient.OrderExec(orderNo)
	if err != nil {
		log.Error("report order exec fail")
	}
}

func (h *RenewVmHandler) Name() string {
	return ResourceOrder_ReNewOrderSuccess
}
