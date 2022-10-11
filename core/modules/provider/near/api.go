package near

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/model"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Near struct {
	composeFileName string
	base            *provider.DockerComposeBase
	network         string
}

func New() *Near {
	return &Near{
		composeFileName: nearComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Near) InitParam(c *gin.Context) error {
	var param model.CommonDeployParam
	err := c.BindJSON(&param)
	if err != nil {
		return err
	}
	s.network = param.Network
	return nil
}

func (s *Near) PullImage() error {
	if err := templateInstance(DeployParams{
		Network: s.network,
	}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Near) Start() error { return s.base.Start(s.composeFileName) }

func (s *Near) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Near) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
