package thegraph

import (
	"fmt"
	"github.com/gorilla/websocket"
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

	//if !isServe {
	//	return errors.New("status is error")
	//}

	err := templateInstance(data)
	if err != nil {
		return err
	}

	if STOP == composeStatus() {
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
		Shell:         "/bin/sh",
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
