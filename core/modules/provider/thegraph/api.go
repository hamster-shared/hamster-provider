package thegraph

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

func GetWebSocket(conn *websocket.Conn) {
	bean := &DockerBean{
		ContainerName: "nginx",
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

func GetDockerLog(conn *websocket.Conn) {
	bean := &DockerBean{
		ContainerName: "nginx",
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
