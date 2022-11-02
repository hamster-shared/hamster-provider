package bsc

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Bsc struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Bsc {
	return &Bsc{
		composeFileName: bscComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (a *Bsc) InitParam(c *gin.Context) error {
	return nil
}

func (b *Bsc) PullImage() error {
	if err := generateRequiredFiles(); err != nil {
		return err
	}
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return b.base.PullImage(b.composeFileName)
}

func (b *Bsc) Start() error { return b.base.Start(b.composeFileName) }

func (b *Bsc) Stop() error { return b.base.Stop(b.composeFileName) }

func (b *Bsc) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return b.base.GetStatus(b.composeFileName, containerIDs...)
}
