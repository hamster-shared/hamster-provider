package context

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/event"
	"github.com/hamster-shared/hamster-provider/core/modules/listener"
	"github.com/hamster-shared/hamster-provider/core/modules/p2p"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
)

// CoreContext the application context , wrapped with some bean
type CoreContext struct {
	P2pClient     *p2p.P2pClient
	VmManager     provider.Manager
	Cm            *config.ConfigManager
	ReportClient  chain.ReportClient
	SubstrateApi  *gsrpc.SubstrateAPI
	TimerService  *utils.TimerService
	EventService  event.IEventService
	ChainListener *listener.ChainListener
	EventContext  *event.EventContext
}

func (c *CoreContext) GetConfig() *config.Config {
	cf, _ := c.Cm.GetConfig()
	return cf
}
