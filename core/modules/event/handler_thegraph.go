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

	err := h.CoreContext.ReportClient.OrderExec(orderNo)

	if err != nil {
		return
	}
	thegraph.SetIsServer(true)
	overdue := h.CoreContext.ReportClient.CalculateInstanceOverdue(orderNo)
	log.GetLogger().Infof("overdue is： %s", overdue)
	instanceTimer := time.NewTimer(overdue)
	h.CoreContext.TimerService.SubTimer(orderNo, instanceTimer)

	//cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	//inspect, err := cli.ContainerInspect(context.Background(), "index-agent")
	//if err != nil {
	//	log.GetLogger().Info("get thegraph cli port fail")
	//	thegraph.SetIsServer(false)
	//	return
	//}
	//ports := inspect.NetworkSettings.Ports["8500/tcp"]
	//ip := "127.0.0.1"
	//graphCliPort := "8500"
	//if len(ports) > 0 {
	//	graphCliPort = ports[0].HostPort
	//} else if val, isKeyExists := inspect.NetworkSettings.Networks["hamster-provider_default"]; isKeyExists {
	//	ip = val.IPAddress
	//}
	//targetListen := fmt.Sprintf("/ip4/%s/tcp/%s", ip, graphCliPort)
	//err = h.CoreContext.P2pClient.Listen("/x/graph-cli", targetListen)
	if err != nil {
		log.GetLogger().Info("setup thegraph cli p2p fail")
		thegraph.SetIsServer(false)
		return
	}

	agreementIndex := h.CoreContext.GetConfig().ChainRegInfo.AgreementIndex

	go func(t *time.Timer) {
		<-t.C
		log.GetLogger().Printf("over due time is : %d, now  terminal", overdue)
		//_, _ = h.CoreContext.P2pClient.Close(targetListen)
		err = thegraph.Uninstall()
		if err != nil {
			thegraph.SetIsServer(false)
		}
		_ = h.CoreContext.ReportClient.ChangeResourceStatus(orderNo)
		h.CoreContext.TimerService.UnSubTicker(agreementIndex)
	}(instanceTimer)

	go func() {
		time.Sleep(time.Second * 10)
		_ = h.CoreContext.ReportClient.Heartbeat(agreementIndex)
		ticker := time.NewTicker(time.Minute * 10)
		h.CoreContext.TimerService.SubTicker(agreementIndex, ticker)
		for {
			<-ticker.C
			// report heartbeat
			_ = h.CoreContext.ReportClient.Heartbeat(agreementIndex)
		}
	}()
}

func (h *TheGraphHandler) Name() string {
	return ResourceOrder_TheGraph
}
