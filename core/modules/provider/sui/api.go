package sui

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Sui struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Sui {
	return &Sui{
		composeFileName: suiComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Sui) InitParam(c *gin.Context) error {

	return nil
}

func (s *Sui) PullImage() error {
	if err := generateRequiredFiles(); err != nil {
		return err
	}
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Sui) Start() error { return s.base.Start(s.composeFileName) }

func (s *Sui) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Sui) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
