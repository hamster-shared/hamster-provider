package aptos

import (
	"context"
	"embed"
	commands "github.com/docker/compose/v2/cmd/compose"
	"os"
	"path/filepath"
	"text/template"

	"github.com/docker/compose/v2/pkg/api"

	"github.com/hamster-shared/hamster-provider/core/modules/compose/client"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
)

var (
	starkwareComposeFileName = "starkware-docker-compose.yml"
	dirName                  = "starkware"
)

type DeployParams struct{}

//go:embed templates/starkware-docker-compose.yaml
var templateFile embed.FS

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/starkware-docker-compose.yaml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), starkwareComposeFileName)
	file, createErr := os.Create(url)
	if createErr != nil {
		log.GetLogger().Errorf("create file failed %s\n", err)
		return createErr
	}
	writeErr := t.Execute(file, deployParam)
	if writeErr != nil {
		log.GetLogger().Errorf("template write file failed %s\n", err)
		return writeErr
	}
	return nil
}

func pullImages() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), starkwareComposeFileName)
	if _, err := os.Stat(composeFilePathName); err != nil {
		_ = templateInstance(DeployParams{})
	}
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "pull"})
}

// StartDockerCompose exec docker-compose
func startDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), starkwareComposeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "up", "-d"})
}

// StopDockerCompose  停止docker compose 服务
func stopDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), starkwareComposeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "down", "-v"})
}

func getDockerComposeStatus(containerIDs ...string) ([]api.ContainerSummary, error) {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), starkwareComposeFileName)
	args := append([]string{"-f", composeFilePathName, "ps"}, containerIDs...)
	err := client.Compose(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return commands.PsCmdResult, nil
}
