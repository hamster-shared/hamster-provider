package aptos

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"

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
	//生成 genesis.blob
	genesisBlobFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "aptos/genesis.blob"))
	if err != nil {
		return err
	}
	_, err = genesisBlobFile.Write(genesisBlob)
	if err != nil {
		return err
	}
	//生成 waypoint.txt
	waypointFile, err := os.Create(filepath.Join(config.DefaultConfigDir(), "aptos/waypoint.txt"))
	if err != nil {
		return err
	}
	_, err = waypointFile.Write(waypointTxtFile)
	if err != nil {
		return err
	}
	//生成 public_full_node.yaml
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
