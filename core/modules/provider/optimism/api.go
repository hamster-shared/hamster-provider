package optimism

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Optimism struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Optimism {
	return &Optimism{
		composeFileName: optimismComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Optimism) PullImage() error {
	if err := generateRequiredFiles(); err != nil {
		return err
	}
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Optimism) Start() error { return s.base.Start(s.composeFileName) }

func (s *Optimism) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Optimism) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
