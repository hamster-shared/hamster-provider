package optimism

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
)

var (
	optimismComposeFileName = "optimism/docker-compose.yml"
)

type DeployParams struct{}

//go:embed templates/docker-compose.yml
var templateFile embed.FS

//go:embed templates
var allFiles embed.FS

func generateRequiredFiles() error {
	// 先创建文件夹
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/docker/grafana/dashboards"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/docker/grafana/provisioning/dashboards"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/docker/grafana/provisioning/datasources"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/docker/prometheus"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/envs/goerli"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/envs/mainnet"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(config.DefaultConfigDir(), "optimism/scripts"), os.ModePerm); err != nil {
		return err
	}

	filenames, err := getAllFilenames(&allFiles, "templates")
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		content, err := allFiles.ReadFile(filename)
		if err != nil {
			return err
		}
		filename = filename[9:]
		_, err = os.Create(filepath.Join(config.DefaultConfigDir(), "optimism", filename))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(config.DefaultConfigDir(), "optimism", filename), content, os.ModePerm); err != nil {
			return err
		}
	}
	err = os.Rename(filepath.Join(config.DefaultConfigDir(), "optimism/env"), filepath.Join(config.DefaultConfigDir(), "optimism/.env"))
	if err != nil {
		return err
	}
	return nil
}

// TemplateInstance Docker compose file instantiation
func templateInstance(deployParam DeployParams) error {
	t, err := template.ParseFS(templateFile, "templates/docker-compose.yml")
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), optimismComposeFileName)
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

func getAllFilenames(fs *embed.FS, dir string) (out []string, err error) {
	if len(dir) == 0 {
		dir = "."
	}

	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fp := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}

			out = append(out, res...)

			continue
		}

		out = append(out, fp)
	}

	return
}
