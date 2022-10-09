package sui

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
)

var (
	suiComposeFileName = "sui-docker-compose.yml"
)

type DeployParams struct{}

//go:embed templates/sui-docker-compose.yaml
var templateFile embed.FS

//go:embed templates/genesis.blob
var genesisBlob []byte

//go:embed templates/fullnode-template.yaml
var fullNodeYaml []byte

func generateRequiredFiles() error {
	// 先创建文件夹
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "sui"), os.ModePerm); err != nil {
		return err
	}
	//生成 genesis.blob
	genesisBlobFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "sui/genesis.blob"))
	if err != nil {
		return err
	}
	_, err = genesisBlobFile.Write(genesisBlob)
	if err != nil {
		return err
	}
	//生成 full_node.yaml
	fullNodeFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "sui/fullnode-template.yaml"))
	if err != nil {
		return err
	}
	_, err = fullNodeFile.Write(fullNodeYaml)
	if err != nil {
		return err
	}
	return nil
}

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/sui-docker-compose.yaml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), suiComposeFileName)
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

