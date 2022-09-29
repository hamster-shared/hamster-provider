package starkware

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Starkware struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Starkware {
	return &Starkware{
		composeFileName: starkwareComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Starkware) PullImage() error {
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Starkware) Start() error { return s.base.Start(s.composeFileName) }

func (s *Starkware) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Starkware) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
