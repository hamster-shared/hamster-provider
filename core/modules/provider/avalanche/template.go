package avalanche

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
)

var (
	avalancheComposeFileName = "avalanche-docker-compose.yml"
)

type DeployParams struct{}

//go:embed templates/docker-compose.yml
var templateFile embed.FS

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/docker-compose.yml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), avalancheComposeFileName)
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

