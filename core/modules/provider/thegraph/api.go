package thegraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hamster-shared/hamster-provider/core/modules/compose/client"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

var isServe = false

func IsServer() bool {
	return isServe
}

func SetIsServer(status bool) {
	isServe = status
}

func Deploy(data DeployParams) error {

	err := templateInstance(data)
	if err != nil {
		return err
	}

	status, err := GetDockerComposeStatus()
	if RUNNING != status || err != nil {
		err = startDockerCompose()
		return err
	}

	return nil
}

func Uninstall() error {
	return stopDockerCompose()
}

// GetWebSocket 获取容器执行的websocket
func GetWebSocket(conn *websocket.Conn, containerName string) {
	bean := &DockerBean{
		ContainerName: containerName,
		Shell:         "/bin/bash",
	}
	// 执行exec，获取到容器终端的连接
	hr, err := dockerExec(bean)
	if err != nil {
		fmt.Println("hello world", err.Error())
		return
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	// 转发输入/输出至websocket
	go func() {
		wsWriterCopy(hr.Conn, conn)
	}()
	wsReaderCopy(conn, hr.Conn)
}

// GetDockerLog 获取容器日志的websocket
func GetDockerLog(conn *websocket.Conn, containerName string) {
	bean := &DockerBean{
		ContainerName: containerName,
		Shell:         "/bin/sh",
	}

	reader, err, cancel := dockerLogs(bean)

	if err != nil {
		fmt.Println("hello world", err.Error())
		return
	}
	// 关闭I/O流
	defer reader.Close()

	go func() {
		wsWriterCopy(reader, conn)
	}()

	time.Sleep(time.Minute * 10)
	cancel()
}

type ComposeStatus int

const (
	STOP        = 0
	RUNNING     = 1
	SOME_EXITED = 2
)

func GetDockerComposeStatus(containerIDs ...string) (ComposeStatus, error) {
	statusResult, err := getDockerComposeStatus(containerIDs...)
	if err != nil {
		return STOP, err
	}
	if len(statusResult) == 0 {
		return STOP, nil
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
		return RUNNING, nil
	} else if exited == len(statusResult) {
		return STOP, nil
	} else {
		return SOME_EXITED, nil
	}
}

func PullImage() error {
	return pullImages()
}

func ComposeGraphConnect(composeFilePathName string) error {
	cmd := fmt.Sprintf("-f %s exec index-cli graph indexer connect http://index-agent:8500", composeFilePathName)
	out, err := client.RunCompose(context.Background(), strings.Split(cmd, " ")...)
	if err != nil {
		log.Errorf("ComposeGraphConnect error: %s", err.Error())
		return err
	}
	log.Debugf("ComposeGraphConnect out: %s", out)
	return nil
}

func ComposeGraphStart(composeFilePathName string, deploymentID string) error {
	cmd := fmt.Sprintf("-f %s exec index-cli graph indexer rules start %s", composeFilePathName, deploymentID)
	out, err := client.RunCompose(context.Background(), strings.Split(cmd, " ")...)
	if err != nil {
		log.Errorf("ComposeGraphStart error: %s", err.Error())
		return fmt.Errorf("ComposeGraphStart error: %s %s", deploymentID, err.Error())
	}
	log.Debugf("ComposeGraphStart out: %s", out)
	return nil
}

func ComposeGraphStop(composeFilePathName string, deploymentID string) error {
	cmd := fmt.Sprintf("-f %s exec index-cli graph indexer rules stop %s", composeFilePathName, deploymentID)
	out, err := client.RunCompose(context.Background(), strings.Split(cmd, " ")...)
	if err != nil {
		log.Errorf("ComposeGraphStop error: %s", err.Error())
		return fmt.Errorf("ComposeGraphStop error: %s %s", deploymentID, err.Error())
	}
	log.Debugf("ComposeGraphStop out: %s", out)
	return nil
}

func ComposeGraphRules(composeFilePathName string) (result []map[string]interface{}, err error) {
	cmd := fmt.Sprintf("-f %s exec index-cli graph indexer rules get all -o json", composeFilePathName)
	out, err := client.RunCompose(context.Background(), strings.Split(cmd, " ")...)
	if err != nil {
		log.Errorf("ComposeGraphRules error: %s", err.Error())
		return nil, err
	}
	log.Debugf("ComposeGraphRules out: %s", out)
	s := strings.Split(out, "]")[0]
	if s == "" {
		return nil, errors.New("no rules")
	}
	s += "]"
	err = json.Unmarshal([]byte(s), &result)
	if err != nil {
		log.Errorf("ComposeGraphRules error: %s", err.Error())
		return
	}
	return
}

func DefaultComposeGraphConnect() error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	return ComposeGraphConnect(composeFilePathName)
}

func DefaultComposeGraphStart(deploymentID ...string) error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	var errorSet []string
	for _, id := range deploymentID {
		err := ComposeGraphStart(composeFilePathName, id)
		if err != nil {
			errorSet = append(errorSet, err.Error())
		}
	}
	if len(errorSet) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errorSet, ", "))
}

func DefaultComposeGraphStop(deploymentID ...string) error {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	var errorSet []string
	for _, id := range deploymentID {
		err := ComposeGraphStop(composeFilePathName, id)
		if err != nil {
			errorSet = append(errorSet, err.Error())
		}
	}
	if len(errorSet) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errorSet, ", "))
}

func DefaultComposeGraphRules() (result []map[string]interface{}, err error) {
	composeFilePathName := filepath.Join(config.DefaultConfigDir(), "docker-compose.yml")
	return ComposeGraphRules(composeFilePathName)
}
