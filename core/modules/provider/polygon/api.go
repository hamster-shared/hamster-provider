package polygon

import (
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
)

type Polygon struct {
	composeFileName string
	base            *provider.DockerComposeBase
}

func New() *Polygon {
	return &Polygon{
		composeFileName: polygonComposeFileName,
		base:            &provider.DockerComposeBase{},
	}
}

func (s *Polygon) PullImage() error {
	if err := generateRequiredFiles(); err != nil {
		return err
	}
	if err := templateInstance(DeployParams{}); err != nil {
		return err
	}
	return s.base.PullImage(s.composeFileName)
}

func (s *Polygon) Start() error { return s.base.Start(s.composeFileName) }

func (s *Polygon) Stop() error { return s.base.Stop(s.composeFileName) }

func (s *Polygon) GetStatus(containerIDs ...string) (provider.ComposeStatus, error) {
	return s.base.GetStatus(s.composeFileName, containerIDs...)
}
