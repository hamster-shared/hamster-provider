package event

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"github.com/hamster-shared/hamster-provider/log"
	"time"
)

type TheGraphHandler struct {
	AbstractHandler
	CoreContext *EventContext
}

func (h *TheGraphHandler) HandlerEvent(e *VmRequest) {

	log.GetLogger().Info("the graph register")

	if thegraph.IsServer() {
		log.GetLogger().Info("current provider is working, cannot accept this request!!")
		return
	}
	orderNo := e.OrderNo

	cm := h.CoreContext.Cm
	cfg, _ := cm.GetConfig()

	peerId := cfg.Identity.PeerID

	err := h.CoreContext.ReportClient.ProcessApplyFreeResource(orderNo, peerId)

	if err != nil {
		return
	}
	thegraph.SetIsServer(true)
	overdue := time.Hour * time.Duration(e.Duration)
	log.GetLogger().Infof("overdue is： %s", overdue)
	instanceTimer := time.NewTimer(overdue)
	h.CoreContext.TimerService.SubTimer(orderNo, instanceTimer)

	go func(t *time.Timer) {
		<-t.C
		log.GetLogger().Printf("over due time is : %d, now  terminal", overdue)
		err = thegraph.Uninstall()
		if err != nil {
			thegraph.SetIsServer(false)
		}
		_ = h.CoreContext.ReportClient.ReleaseApplyFreeResource(orderNo)
	}(instanceTimer)
}

func (h *TheGraphHandler) Name() string {
	return ResourceOrder_TheGraph
}
