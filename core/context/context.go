package context

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"github.com/hamster-shared/hamster-provider/core/modules/listener"
	"github.com/hamster-shared/hamster-provider/core/modules/p2p"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/hamster-shared/hamster-provider/log"
	"sync"
)

// CoreContext the application context , wrapped with some bean
type CoreContext struct {
	P2pClient     *p2p.P2pClient
	VmManager     provider.Manager
	Cm            *config.ConfigManager
	ReportClient  chain.ReportClient
	SubstrateApi  *gsrpc.SubstrateAPI
	TimerService  *utils.TimerService
	EventContext  *event.EventContext
	EventService  event.IEventService
	ChainListener *listener.ChainListener
}

func (c *CoreContext) GetConfig() *config.Config {
	cf, _ := c.Cm.GetConfig()
	return cf
}

func (c *CoreContext) InitSubstrate() error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if c.SubstrateApi == nil {
		return c.ResetSubstrate()
	} else {
		return nil
	}
}

func (c *CoreContext) ResetSubstrate() error {
	substrateApi, err := gsrpc.NewSubstrateAPI(c.GetConfig().ChainApi)
	if err != nil {
		log.GetLogger().Error("init substrate error : ", err.Error())
		return err
	}
	reportClient, err := chain.NewChainClient(c.Cm, substrateApi)
	if err != nil {
		return err
	}

	if c.SubstrateApi != nil {
		c.SubstrateApi = nil
		c.ReportClient = nil
	}
	c.SubstrateApi = substrateApi
	c.ReportClient = reportClient
	c.ChainListener.SetChainApi(substrateApi, reportClient)

	return nil
}

func (c *CoreContext) InitP2p() error {
	if c.P2pClient != nil {
		_ = c.P2pClient.Destroy()
		c.P2pClient = nil
	}

	p2pClient, err := p2p.NewP2pClient(34001, c.GetConfig().Identity.PrivKey, c.GetConfig().Identity.SwarmKey, c.GetConfig().Bootstraps)
	if err != nil {
		return err
	}
	err = p2pClient.Listen("/x/provider", fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", c.GetConfig().ApiPort+1))
	if err != nil {
		return err
	}
	c.P2pClient = p2pClient

	c.EventContext.P2pClient = p2pClient
	c.EventContext.ReportClient = c.ReportClient

	return nil
}
