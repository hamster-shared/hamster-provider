package vm

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	containerName = "some-test"
	publicKey     = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCWCDS4Io8+PFqGepqy0YNrtw3B7g7lhg7WNcH2VyJmmHlvft69N3S4EzDugEUDgPbgihiL56wyq56GtOG6+RuRuqkEU983MRC6j0yazem/KPs2nAS0NW5A8Nzxm9ixXnF9Bw6qHpO+L8ZbKdsIR+xux5QVriWTmDd/FeaovzRa/Ogr/BdShsp5H1s8aKkj2ygm16rlWAuQcoPQJPDWJVLM9cub8wj/AGrOzRDQCMnbcm69BZT7GPbodVmBIlugICuSVVvKSpZEa0QHCdQW2z2kIan7EwEI7LPYyDpCRAAI2mYEsl9WIIzae1ACK7dKwp9DKfLlKU4YRNfvR5stGgNezelz2pbN0TvK0T6NrqlKDo1eZQbHzRzvUKtDCiwSBdauJVus5Zowqy8lXr9wVosbra8z8cd+vM+e5+82fEjnE3BQm6NUHatOfxe/1MtYeem1Zlru5ISc25ceCXJd/l6qUrlIamHKgMpxvAt8g9pcPpPH2YozLkRlohcdWrhA+kk= gr@gr-Lenovo"
)

func getTemplate() Template {
	return Template{
		Cpu:        1,
		Memory:     1,
		Disk:       50,
		System:     "ubuntu",
		PublicKey:  publicKey,
		Image:      "rastasheep/ubuntu-sshd:18.04",
		AccessPort: 22,
	}
}

func TestSetTemplate(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)
	assert.NotEmpty(t, client)
}

func TestCreate(t *testing.T) {
	client, err := NewDockerManager(getTemplate())

	assert.NoError(t, err)

	id, err := client.Create(containerName)
	defer clean(id)
	assert.NoError(t, err)

	containers, err := client.cli.ContainerList(client.ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", containerName)),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, containers)
	assert.Equal(t, "/"+containerName, containers[0].Names[0])
}

func TestStart(t *testing.T) {

	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.Create(containerName)
	defer clean(id)
	assert.NoError(t, err)

	err = client.Start(containerName)
	assert.NoError(t, err)

	status, err := client.Status(containerName)
	assert.NoError(t, err)

	assert.Equal(t, 1, status.status)
	assert.Equal(t, id, status.id)
}

func TestCreateAndStart(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.CreateAndStart(containerName)
	assert.NoError(t, err)
	defer clean(id)

	status, err := client.Status(containerName)
	assert.NoError(t, err)

	assert.Equal(t, 1, status.status)
	assert.Equal(t, id, status.id)
}

func TestCreateAndStartAndInjectionPublicKey(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.CreateAndStartAndInjectionPublicKey(containerName, publicKey)
	assert.NoError(t, err)
	defer clean(id)

	status, err := client.Status(containerName)
	assert.NoError(t, err)

	assert.Equal(t, 1, status.status)
	assert.Equal(t, id, status.id)
}

func TestStop(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.CreateAndStart(containerName)
	assert.NoError(t, err)
	defer clean(id)

	err = client.Stop(containerName)
	assert.NoError(t, err)

	status, err := client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 0, status.status)
}

func TestReboot(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)
	id, err := client.Create(containerName)
	defer clean(id)

	err = client.Reboot(containerName)
	assert.NoError(t, err)

	status, err := client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 1, status.status)
}

func TestShutdown(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.CreateAndStart(containerName)
	assert.NoError(t, err)
	defer clean(id)

	status, err := client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 1, status.status)

	err = client.Shutdown(containerName)
	assert.NoError(t, err)

	status, err = client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 0, status.status)
}

func TestDestroy(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	_, err = client.Create(containerName)
	assert.NoError(t, err)

	err = client.Destroy(containerName)
	assert.NoError(t, err)

	_, err = client.Status(containerName)
	assert.Error(t, err)
}

func TestStatus(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.Create(containerName)
	defer clean(id)
	assert.NoError(t, err)

	status, err := client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 0, status.status)

	err = client.Start(containerName)
	assert.NoError(t, err)

	status, err = client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 1, status.status)

	err = client.Stop(containerName)
	assert.NoError(t, err)

	status, err = client.Status(containerName)
	assert.NoError(t, err)
	assert.Equal(t, 0, status.status)
}

func TestGetIp(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.Create(containerName)
	defer clean(id)
	assert.NoError(t, err)

	ip, err := client.GetIp(containerName)
	assert.NoError(t, err)

	assert.NotEmpty(t, ip)
}

func TestGetPort(t *testing.T) {
	client, err := NewDockerManager(getTemplate())
	assert.NoError(t, err)

	id, err := client.CreateAndStart(containerName)
	defer clean(id)
	assert.NoError(t, err)

	port := client.GetAccessPort(containerName)

	assert.NotEmpty(t, port)
}

func clean(id string) {
	client, _ := NewDockerManager(getTemplate())
	err := client.cli.ContainerRemove(client.ctx, id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Println("err:", err)
	}
}
