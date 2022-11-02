package bsc

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
)

var (
	bscComposeFileName = "bsc-docker-compose.yml"
)

type DeployParams struct{}

//go:embed templates/bsc-docker-compose.yaml
var templateFile embed.FS

//go:embed templates/config/config.toml
var configToml []byte

//go:embed templates/config/genesis.json
var genesisToml []byte

func generateRequiredFiles() error {
	// 先创建文件夹
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "bsc/config"), os.ModePerm); err != nil {
		return err
	}
	//生成 config.toml
	configTomlFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "bsc/config/config.toml"))
	if err != nil {
		return err
	}
	_, err = configTomlFile.Write(configToml)
	if err != nil {
		return err
	}
	//生成 genesis.json
	genesisTomlFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "bsc/config/genesis.json"))
	if err != nil {
		return err
	}
	_, err = genesisTomlFile.Write(genesisToml)
	if err != nil {
		return err
	}
	return nil
}

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/bsc-docker-compose.yaml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), bscComposeFileName)
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
