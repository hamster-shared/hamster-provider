package vm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getVmManager() (*VirtManager, error) {
	template := Template{
		Cpu:        1,
		Memory:     2,
		Disk:       50,
		System:     "ubuntu",
		Image:      "https://s3.ttchain.tntlinking.com/compute/ubuntu.vhdx.tar.gz",
		AccessPort: 22,
	}
	return NewVirtManager(template)
}

func TestCreate(t *testing.T) {

	template := Template{
		Cpu:        2,
		Memory:     4,
		Disk:       50,
		System:     "ubuntu",
		Image:      "https://s3.ttchain.tntlinking.com/compute/windows10.vhdx.tar.gz",
		AccessPort: 22,
	}
	vmManager, err := NewVirtManager(template)
	assert.NoError(t, err)

	vmName := "test_win2"

	err = vmManager.Create(vmName)

	assert.NoError(t, err)

	status, err := vmManager.Status(vmName)
	assert.NoError(t, err)

	assert.Equal(t, 0, status.status)

	err = vmManager.Start(vmName)
	assert.NoError(t, err)

	ip, err := vmManager.GetIp(vmName)

	assert.NoError(t, err)
	fmt.Println(ip)
}

func TestGetStatus(t *testing.T) {
	vmName := "test_create"
	manager, _ := getVmManager()
	result, err := manager.Status(vmName)
	assert.NoError(t, err)
	fmt.Println(result.status)

}

func TestStart(t *testing.T) {

	vmName := "test_create"
	manager, _ := getVmManager()
	err := manager.Start(vmName)
	assert.NoError(t, err)
	result, err := manager.Status(vmName)
	assert.NoError(t, err)

	assert.Equal(t, 1, result.status)
}

func TestStop(t *testing.T) {
	vmName := "test_create"
	manager, _ := getVmManager()
	err := manager.Stop(vmName)
	assert.NoError(t, err)
	result, err := manager.Status(vmName)
	assert.NoError(t, err)

	assert.Equal(t, 0, result.status)
}

func TestDelete(t *testing.T) {
	vmName := "test_win"
	manager, _ := getVmManager()
	err := manager.Destroy(vmName)
	assert.NoError(t, err)
}

func TestGetIp(t *testing.T) {
	vmName := "test_win2"
	manager, _ := getVmManager()
	ip, err := manager.GetIp(vmName)
	assert.NoError(t, err)
	fmt.Println(ip)
}
