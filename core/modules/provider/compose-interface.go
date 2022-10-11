package provider

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"

	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"

	"github.com/hamster-shared/hamster-provider/core/modules/compose/client"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
)

type Chain interface {
	DockerCompose
}

type DockerCompose interface {
	InitParam(c *gin.Context) error
	PullImage() error
	Start() error
	Stop() error
	GetStatus(containerIDs ...string) (ComposeStatus, error)
}

type ComposeStatus int

const (
	ComposeStop ComposeStatus = iota
	ComposeRunning
	ComposeSomeExited
)

type DockerComposeBase struct{}

func (*DockerComposeBase) PullImage(composeFileName string) error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), composeFileName)
	if _, err := os.Stat(composeFilePathName); err != nil {
		return err
	}
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "pull"})
}

func (*DockerComposeBase) Start(composeFileName string) error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), composeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "up", "-d"})
}

func (*DockerComposeBase) Stop(composeFileName string) error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), composeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "down", "-v"})
}

func (*DockerComposeBase) GetStatus(composeFileName string, containerIDs ...string) (ComposeStatus, error) {
	fmt.Println("get compose status: ", containerIDs)
	statusResult, err := getDockerComposeStatus(composeFileName, containerIDs...)
	if err != nil {
		return ComposeStop, err
	}
	if len(statusResult) == 0 {
		return ComposeStop, nil
	}
	running := 0
	exited := 0
	for _, status := range statusResult {
		switch status.State {
		case "running":
			running++
		case "exited":
			exited++
		}
	}
	if running == len(statusResult) {
		return ComposeRunning, nil
	} else if exited == len(statusResult) {
		return ComposeStop, nil
	} else {
		return ComposeSomeExited, nil
	}
}

func getDockerComposeStatus(composeFileName string, containerIDs ...string) ([]api.ContainerSummary, error) {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), composeFileName)
	args := append([]string{"-f", composeFilePathName, "ps"}, containerIDs...)
	err := client.Compose(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return commands.PsCmdResult, nil
}
