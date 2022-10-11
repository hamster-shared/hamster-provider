package ethereum

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/model"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Ethereum struct {
	composeFileName string
	base            *provider.DockerComposeBase
	network         string
}

func New() *Ethereum {
	return &Ethereum{
		composeFileName: ethereumComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (e *Ethereum) InitParam(c *gin.Context) error {
	var param model.CommonDeployParam
	err := c.BindJSON(&param)
	if err != nil {
		return err
	}
	e.network = param.Network
	return nil
}

func (s *Ethereum) PullImage() error {
	if err := templateInstance(DeployParams{
		Network: s.network,
	}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Ethereum) Start() error { return s.base.Start(s.composeFileName) }

func (s *Ethereum) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Ethereum) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
