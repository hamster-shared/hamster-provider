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
	aptosComposeFileName = "aptos-docker-compose.yml"
)

type DeployParams struct{}

//go:embed templates/aptos-docker-compose.yaml
var templateFile embed.FS

//go:embed templates/genesis.blob
var genesisBlob []byte

//go:embed templates/public_full_node.yaml
var publicFullNodeYaml []byte

//go:embed templates/waypoint.txt
var waypointTxtFile []byte

func generateRequiredFiles() error {
	// 先创建文件夹
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "aptos"), os.ModePerm); err != nil {
		return err
	}
	//生成genesis.blob
	genesisBlobFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "aptos/genesis.blob"))
	if err != nil {
		return err
	}
	_, err = genesisBlobFile.Write(genesisBlob)
	if err != nil {
		return err
	}
	//生成waypoint.txt
	waypointFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "aptos/waypoint.txt"))
	if err != nil {
		return err
	}
	_, err = waypointFile.Write(waypointTxtFile)
	if err != nil {
		return err
	}
	//生成public_full_node.yaml
	publicFullNodeFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "aptos/public_full_node.yaml"))
	if err != nil {
		return err
	}
	_, err = publicFullNodeFile.Write(publicFullNodeYaml)
	if err != nil {
		return err
	}
	return nil
}

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/aptos-docker-compose.yaml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), aptosComposeFileName)
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
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), aptosComposeFileName)
	if _, err := os.Stat(composeFilePathName); err != nil {
		_ = templateInstance(DeployParams{})
		_ = generateRequiredFiles()
	}
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "pull"})
}

// StartDockerCompose exec docker-compose
func startDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), aptosComposeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "up", "-d"})
}

// StopDockerCompose  停止docker compose 服务
func stopDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), aptosComposeFileName)
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "down", "-v"})
}

func getDockerComposeStatus(containerIDs ...string) ([]api.ContainerSummary, error) {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), aptosComposeFileName)
	args := append([]string{"-f", composeFilePathName, "ps"}, containerIDs...)
	err := client.Compose(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return commands.PsCmdResult, nil
}
