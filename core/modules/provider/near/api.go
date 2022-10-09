package near

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Near struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Near {
	return &Near{
		composeFileName: nearComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Near) PullImage() error {
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Near) Start() error { return s.base.Start(s.composeFileName) }

func (s *Near) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Near) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
