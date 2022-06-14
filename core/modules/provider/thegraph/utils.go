package thegraph

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"io"
)

var docker_cli *client.Client

func GetDockerClient() *client.Client {
	if docker_cli == nil {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return nil
		}
		docker_cli = cli
	}
	return docker_cli
}

type DockerBean struct {
	ID            int    `gorm:"primary_key" json:"id"`
	Title         string `json:"title"`
	ImageName     string `json:"image_name"`
	PortBind      string `json:"port_bind"`
	PortExport    string `json:"port_export"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_id"`
	Mark          string `json:"mark"`
	Status        int    `json:"status"` // 1 start, 2 停止
	Shell         string `json:"shell"`
	UserId        int    `json:"user_id"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
}

func dockerExec(bean *DockerBean) (hr types.HijackedResponse, err error) {
	// 执行/bin/bash命令
	cli := GetDockerClient()
	ctx := context.Background()
	ir, err := cli.ContainerExecCreate(ctx, bean.ContainerName, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{bean.Shell}, //这个是必须选的，以sh或bash进入容器中
		Tty:          true,
	})
	if err != nil {
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	fmt.Println(ir.ID)
	return
}

func dockerLogs(bean *DockerBean) (reader io.ReadCloser, err error, cancel func()) {
	cli := GetDockerClient()
	ctx, cancel := context.WithCancel(context.Background())

	reader, err = cli.ContainerLogs(ctx, bean.ContainerName, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})

	return
}

func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			fmt.Println(string(buf[0:nr]))
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}
