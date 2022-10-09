package ethereum

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Ethereum struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Ethereum {
	return &Ethereum{
		composeFileName: ethereumComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Ethereum) PullImage() error {
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Ethereum) Start() error { return s.base.Start(s.composeFileName) }

func (s *Ethereum) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Ethereum) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
