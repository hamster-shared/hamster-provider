package vm

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type DockerManager struct {
	template *Template
	//access port
	accessPort int
	// docker client
	cli *client.Client
	// context
	ctx context.Context

	// image
	image string
}

func NewDockerManager(t Template) (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	manager := &DockerManager{
		cli: cli,
		ctx: context.Background(),
	}
	err = manager.SetTemplate(t)
	return manager, err
}

func (d *DockerManager) SetTemplate(t Template) error {
	d.template = &t
	d.image = t.Image
	d.accessPort = 22
	return nil
}

func (d *DockerManager) Status(name string) (*Status, error) {
	containers, err := d.cli.ContainerList(d.ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", name)),
	})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return &Status{}, errors.New("container not exists")
	}

	// status Status 0: Off, 1: On, 2: Paused, 3, other
	var status int

	switch containers[0].State {
	case "created":
		status = 0
		break
	case "restarting":
		status = 3
		break
	case "running":
		status = 1
		break
	case "removing":
		status = 3
		break
	case "paused":
		status = 2
		break
	case "exited":
	case "dead":
		status = 0
		break
	}

	return &Status{
		status: status,
		id:     containers[0].ID,
	}, nil
}

// query container ip address
func (d *DockerManager) GetIp(name string) (string, error) {
	//status, err := d.Status()
	//
	//if err != nil {
	//	return "", err
	//}
	//
	//if status.id == "" {
	//	return "", errors.New("container id cannot be empty")
	//}
	//
	//containerJson, err := d.cli.ContainerInspect(d.ctx, status.id)
	//if err != nil {
	//	return "", err
	//}
	//return containerJson.NetworkSettings.IPAddress, nil

	// macos cannot directly access container ip
	return "127.0.0.1", nil
}

func (d *DockerManager) GetAccessPort(name string) int {
	inspect, err := d.cli.ContainerInspect(d.ctx, name)
	if err != nil {
		return 0
	}
	portMap := inspect.NetworkSettings.Ports
	port, _ := nat.NewPort("tcp", strconv.Itoa(d.accessPort))
	arrays := portMap[port]
	if len(arrays) > 0 {
		hostPort, _ := strconv.Atoi(arrays[0].HostPort)
		return hostPort
	}
	return 0
}

func (d *DockerManager) Create(name string) (string, error) {

	// view all images
	imageLists, err := d.cli.ImageList(d.ctx, types.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("reference", d.image)),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(imageLists) == 0 {
		// pull image
		out, err := d.cli.ImagePull(d.ctx, d.image, types.ImagePullOptions{})
		if err != nil {
			log.Println(err)
			return "", err
		}
		_, err = io.Copy(os.Stdout, out)
		if err != nil {
			return "", err
		}
	}

	// determine whether there is a repeated start
	status, err := d.Status(name)
	if err != nil {
		log.Info(err.Error())
	}

	if status.id != "" {
		err = d.cli.ContainerRemove(d.ctx, status.id, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return status.id, err
		}
	}

	port, err := nat.NewPort("tcp", strconv.Itoa(d.accessPort))
	// create a container
	resp, err := d.cli.ContainerCreate(d.ctx, &container.Config{
		Image: d.image, //image name
		//Tty:        true,
		//OpenStdin:  true,
		//Cmd:        []string{cmd},
		//WorkingDir: workDir,
		ExposedPorts: nat.PortSet{
			port: struct{}{}, //docker container open port
		},
	},
		&container.HostConfig{

			Resources: container.Resources{
				CPUCount: int64(d.template.Cpu),
				Memory:   int64(d.template.Memory << 30),
			},
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{
					{
						HostPort: strconv.Itoa(utils.RandomPort()),
					},
				},
			},
		}, nil, nil, name)

	if err != nil {
		log.Println(err)
		return resp.ID, err
	}

	log.WithField("containerId", resp.ID).Info("container created")

	return resp.ID, err
}

// StartContainer running containers in the background
func (d *DockerManager) Start(name string) error {
	status, err := d.Status(name)
	if err != nil {
		return err
	}

	id := status.id

	if status.status == 0 && id != "" {
		err = d.cli.ContainerStart(d.ctx, id, types.ContainerStartOptions{})
		return err
	} else {
		return errors.New("container status is invalid")
	}

}

func (d *DockerManager) CreateAndStart(name string) (string, error) {
	id, err := d.Create(name)
	if err != nil {
		return id, err
	}
	err = d.Start(name)
	return id, err
}

func (d *DockerManager) CreateAndStartAndInjectionPublicKey(name, publicKey string) (string, error) {
	// create a virtual machine
	id, err := d.CreateAndStart(name)
	if err != nil {
		return id, err
	}
	// wait for the virtual machine to start successfully
	for {
		status, err := d.Status(name)
		if err != nil {
			return id, err
		}
		if status.IsRunning() {
			break
		}
		time.Sleep(time.Second * 3)
	}
	err = d.InjectionPublicKey(name, publicKey)
	return id, err
}

func (d *DockerManager) Stop(name string) error {
	status, err := d.Status(name)
	if status.status != 1 {
		return errors.New("invalid container status")
	}
	if err != nil {
		return err
	}
	id := status.id

	timeout := time.Second * 3
	if id != "" {
		return d.cli.ContainerStop(d.ctx, status.id, &timeout)
	} else {
		return errors.New("container id is invalid")
	}
}

func (d *DockerManager) Reboot(name string) error {
	status, err := d.Status(name)
	if err != nil {
		return err
	}
	id := status.id

	timeout := time.Second * 3
	if id != "" {
		return d.cli.ContainerRestart(d.ctx, status.id, &timeout)
	} else {
		return errors.New("container id is invalid")
	}
}

func (d *DockerManager) Shutdown(name string) error {
	status, err := d.Status(name)
	if status.status != 1 {
		return errors.New("invalid container status")
	}
	if err != nil {
		return err
	}
	id := status.id
	timeout := time.Second * 3
	if id != "" {
		return d.cli.ContainerStop(d.ctx, status.id, &timeout)
	} else {
		return errors.New("container id is invalid")
	}
}

// Destroy delete container
func (d *DockerManager) Destroy(name string) error {
	status, err := d.Status(name)
	if err != nil {
		return err
	}
	id := status.id

	if id != "" {
		return d.cli.ContainerRemove(d.ctx, status.id, types.ContainerRemoveOptions{Force: true})
	} else {
		return errors.New("container id is invalid")
	}
}

// InjectionPublicKey add the public key to the container
func (d *DockerManager) InjectionPublicKey(name, publicKey string) error {

	status, err := d.Status(name)
	if !status.IsRunning() {
		return errors.New("invalid container status")
	}
	if err != nil {
		return err
	}
	id := status.id
	if id != "" {
		cmd := fmt.Sprintf("echo %s  > /root/.ssh/authorized_keys", publicKey)
		command := exec.Command("docker", "exec", id, "bash", "-c", cmd)
		return command.Run()
	} else {
		return errors.New("container status is invalid")
	}
}
