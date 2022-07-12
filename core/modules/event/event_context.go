package event

import (
	"github.com/hamster-shared/hamster-provider/core/modules/chain"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/core/modules/p2p"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
)

type EventContext struct {
	ReportClient chain.ReportClient
	VmManager    provider.Manager
	TimerService *utils.TimerService
	Cm           *config.ConfigManager
	P2pClient    *p2p.P2pClient
}

func (ec *EventContext) GetConfig() *config.Config {
	cfg, _ := ec.Cm.GetConfig()
	return cfg
}
