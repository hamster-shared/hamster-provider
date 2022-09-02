package thegraph

import (
	"context"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/hamster-shared/hamster-provider/core/modules/compose/client"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/hamster-shared/hamster-provider/log"
	"os"
	"path/filepath"
	"text/template"
)

// InstallDocker install docker
//func InstallDocker() error {
//	//Determine whether to save and install docker
//	cmd := exec.Command("docker", "-v")
//	_, err := cmd.CombinedOutput()
//	if err != nil {
//		//install docker
//		curlCmd := exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh")
//		curlErr := curlCmd.Run()
//		if curlErr != nil {
//			log.GetLogger().Errorf("cmd.Run() failed with %s\n", curlErr)
//			return curlErr
//		}
//		//installCmd := "sh 'get-docker.sh'"
//		//exec install docker
//		//child, installErr := gexpect.Spawn(installCmd)
//		//if installErr != nil {
//		//	log.GetLogger().Error("cmd.Run() failed with %s\n", installErr)
//		//	return installErr
//		//}
//		//input sudo password
//		//pwd := "some-password"
//		//sendErr := child.SendLine(pwd)
//		//if sendErr != nil {
//		//	log.Printf("password error")
//		//	return sendErr
//		//}
//		//waitErr := child.Wait()
//		//if waitErr != nil {
//		//	log.Printf("cmd.Run() failed with %s\n", waitErr)
//		//	return waitErr
//		//}
//		//start docker
//		startCmd := "sudo service docker start"
//		//startCmd := exec.Command("systemctl","start","docker")
//		//startErr := startCmd.Run()
//		childStart, startErr := gexpect.Spawn(startCmd)
//		if startErr != nil {
//			log.GetLogger().Errorf("cmd.Run() failed with %s\n", startErr)
//			return startErr
//		}
//		//input password
//		//startSendErr := childStart.SendLine(pwd)
//		//if startSendErr != nil {
//		//	log.GetLogger().Error("start pws error")
//		//	return startSendErr
//		//}
//		startWaitErr := childStart.Wait()
//		if startWaitErr != nil {
//			log.GetLogger().Error("start wait error")
//			return startWaitErr
//		}
//		return nil
//	}
//	//start docker
//	//startCmd := "sudo service docker start"
//	//startCmd := exec.Command("systemctl","start","docker")
//	//startErr := startCmd.Run()
//	//childStart, startErr := gexpect.Spawn(startCmd)
//	//if startErr != nil {
//	//	log.GetLogger().Errorf("cmd.Run() installed start failed with %s\n", startErr)
//	//	return startErr
//	//}
//	//startSendErr := childStart.SendLine(pwd)
//	//if startSendErr != nil {
//	//	log.GetLogger().Error("installed start password error")
//	//	return startSendErr
//	//}
//	return nil
//}

// TemplateInstance Docker compose file instantiation
func templateInstance(data DeployParams) error {

	pathExecutable, _ := os.Executable()
	path := filepath.Join(filepath.Dir(pathExecutable), "templates/graph-docker-compose.text")
	t, err := template.ParseFiles(path)
	if err != nil {
		log.GetLogger().Errorf("template failed with %s\n", err)
		return err
	}
	//create file in .hamster-provider
	url := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	file, createErr := os.Create(url)
	if createErr != nil {
		log.GetLogger().Errorf("create file failed %s\n", err)
		return createErr
	}
	writeErr := t.Execute(file, data)
	if writeErr != nil {
		log.GetLogger().Errorf("template write file failed %s\n", err)
		return writeErr
	}
	return nil
}

func pullImages() error {
	pathExecutable, _ := os.Executable()
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	var path string
	if _, err := os.Stat(composeFilePathName); err != nil {
		path = filepath.Join(filepath.Dir(pathExecutable), "templates/graph-docker-compose.text")
	} else {
		path = composeFilePathName
	}

	return client.Compose(context.Background(), []string{"-f", path, "pull"})
}

// StartDockerCompose exec docker-compose
func startDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "up", "-d"})
}

// StopDockerCompose  停止docker compose 服务
func stopDockerCompose() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	return client.Compose(context.Background(), []string{"-f", composeFilePathName, "down", "-v"})
}

func getDockerComposeStatus(containerIDs ...string) ([]api.ContainerSummary, error) {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	args := append([]string{"-f", composeFilePathName, "ps"}, containerIDs...)
	err := client.Compose(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return commands.PsCmdResult, nil
}
