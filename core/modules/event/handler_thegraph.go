package event

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
	"time"
)

type TheGraphHandler struct {
	AbstractHandler
	CoreContext EventContext
}

func (h *TheGraphHandler) HandlerEvent(e *VmRequest) {

	fmt.Println("the graph register")

	if thegraph.IsServer() {
		return
	}
	orderNo := e.OrderNo

	cm := h.CoreContext.Cm
	cfg, _ := cm.GetConfig()

	peerId := cfg.Identity.PeerID

	err := h.CoreContext.ReportClient.ProcessApplyFreeResource(orderNo, peerId)

	if err != nil {
		thegraph.SetIsServer(true)
	}

	overdue := time.Second * 6 * time.Duration(e.Duration)
	fmt.Printf("overdue is： %s", overdue)
	instanceTimer := time.NewTimer(overdue)
	h.CoreContext.TimerService.SubTimer(orderNo, instanceTimer)

	go func(t *time.Timer) {
		<-t.C
		fmt.Printf("over due time is : %d, now  terminal", overdue)
		err = thegraph.Uninstall()
		if err != nil {
			thegraph.SetIsServer(false)
		}
	}(instanceTimer)
}

func (h *TheGraphHandler) Name() string {
	return ResourceOrder_TheGraph
}
