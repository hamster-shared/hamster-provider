package aptos

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Aptos struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Aptos {
	return &Aptos{
		composeFileName: aptosComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (a *Aptos) PullImage() error {
	if err := generateRequiredFiles(); err != nil {
		return err
	}
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return a.base.PullImage(a.composeFileName)
}

func (a *Aptos) Start() error { return a.base.Start(a.composeFileName) }

func (a *Aptos) Stop() error { return a.base.Stop(a.composeFileName) }

func (a *Aptos) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return a.base.GetStatus(a.composeFileName, containerIDs...)
}
