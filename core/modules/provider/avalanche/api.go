package avalanche

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Avalanche struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Avalanche {
	return &Avalanche{
		composeFileName: avalancheComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Avalanche) InitParam(c *gin.Context) error {

	return nil
}

func (s *Avalanche) PullImage() error {
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Avalanche) Start() error { return s.base.Start(s.composeFileName) }

func (s *Avalanche) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Avalanche) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
