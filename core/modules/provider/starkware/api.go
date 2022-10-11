package starkware

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/model"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Starkware struct {
	composeFileName string
	base            *provider.DockerComposeBase
	ethereumApiUrl  string
}

func New() *Starkware {
	return &Starkware{
		composeFileName: starkwareComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Starkware) InitParam(c *gin.Context) error {

	var param model.StarkwareDeployParam
	err := c.BindJSON(&param)
	if err != nil {
		return err
	}
	s.ethereumApiUrl = param.EthereumApiUrl
	return nil
}

func (s *Starkware) PullImage() error {
	if err := templateInstance(DeployParams{
		s.ethereumApiUrl,
	}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Starkware) Start() error { return s.base.Start(s.composeFileName) }

func (s *Starkware) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Starkware) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
