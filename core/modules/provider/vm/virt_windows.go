package vm

import (
	"errors"
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/hamster-shared/hamster-provider/log"
	"os"
	"path/filepath"
	"strings"
)

type VirtManager struct {
	home       string
	template   *Template
	accessPort int
}

func NewVirtManager(t Template) (*VirtManager, error) {
	homedir, err := os.UserHomeDir()
	manager := &VirtManager{
		home:       filepath.Join(homedir, ".hamster-provider"),
		accessPort: t.AccessPort,
	}
	err = manager.SetTemplate(t)
	return manager, err
}

func (v *VirtManager) SetTemplate(t Template) error {
	v.template = &t
	baeImage := filepath.Join(v.home, filepath.Base(v.template.Image))
	if _, err := os.Stat(baeImage); errors.Is(err, os.ErrNotExist) {
		log.GetLogger().Info("start download template")

		// download file and rename
		err = utils.Download(v.template.Image, baeImage)
		if err != nil {
			log.GetLogger().Error("download template fail")
			return err
		}
	}

	if _, err := os.Stat(v.getBaseImagePath()); errors.Is(err, os.ErrNotExist) {
		if strings.HasSuffix(baeImage, ".tar.gz") {
			file, err := os.Open(baeImage)
			if err != nil {
				log.GetLogger().Error("download template fail")
				return err
			}
			err = utils.UnTar(file, v.getBaseImagePath())
			if err != nil {
				fmt.Println("untar :", err)
				return err
			}
		}
	}
	return nil
}

func (v *VirtManager) getOrderPath(name string) string {
	return filepath.Join(v.home, "orders", name)
}

func (v *VirtManager) getCopyDiskFile(name string) string {
	return filepath.Join(v.getOrderPath(name), fmt.Sprintf("%s_%s", name, v.getBaseImageName()))
}

func (v *VirtManager) getBaseImageName() string {
	imageName := filepath.Base(v.template.Image)
	if strings.HasSuffix(imageName, ".tar.gz") {
		imageName = strings.ReplaceAll(imageName, ".tar.gz", "")
	}
	return imageName
}

func (v *VirtManager) getBaseImagePath() string {
	return filepath.Join(v.home, v.getBaseImageName())
}

func (v *VirtManager) Create(name string) error {

	switchName, err := utils.GetDefaultVirtualSwitch()

	if err != nil {
		return err
	}

	err = utils.CreateVirtualMachine(
		name,
		v.getOrderPath(name),
		v.getBaseImagePath(),
		v.getOrderPath(name),
		int64(v.template.Memory*1<<30),
		int64(v.template.Disk*1<<30),
		switchName,
		1,
		false,
	)
	return err
}

func (v *VirtManager) Start(name string) error {
	return utils.StartVirtualMachine(name)
}

func (v *VirtManager) CreateAndStart(name string) error {
	err := v.Create(name)
	if err != nil {
		return v.Start(name)
	}
	return err
}

func (v *VirtManager) CreateAndStartAndInjectionPublicKey(name, publicKey string) error {
	return v.CreateAndStart(name)
}

func (v *VirtManager) Stop(name string) error {
	return utils.StopVirtualMachine(name)
}

func (v *VirtManager) Reboot(name string) error {
	return utils.RestartVirtualMachine(name)
}

func (v *VirtManager) Shutdown(name string) error {
	return utils.ShutDown(name)
}

func (v *VirtManager) Destroy(name string) error {
	return utils.DeleteVirtualMachine(name)
}

func (v *VirtManager) InjectionPublicKey(name, publicKey string) error {
	return errors.New("not support now")
}

func (v *VirtManager) Status(name string) (*Status, error) {
	isRunning, err := utils.IsRunning(name)

	var status int
	if err != nil {
		status = 2
	} else if isRunning {
		status = 1
	} else {
		status = 0
	}

	return &Status{
		id:     name,
		status: status,
	}, err

}

func (v *VirtManager) GetIp(name string) (string, error) {
	return utils.GetVmIpAddress(name)
}

func (v *VirtManager) GetAccessPort(name string) int {
	return v.accessPort
}
