package vm

// dependency
// yum install libvirt-devel gcc
// brew install libvirt
// apt install libvirt-dev

import (
	"errors"
	"fmt"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	libvirt "github.com/libvirt/libvirt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// VirtManager virtual machine management client
type VirtManager struct {
	conn       *libvirt.Connect
	home       string
	template   *provider.Template
	accessPort int
}

// NewVirtManager create virtManager
func NewVirtManager(t provider.Template) (*VirtManager, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	homedir, err := os.UserHomeDir()
	manager := &VirtManager{
		conn: conn,
		home: homedir + "/.hamster-provider",
	}
	err = manager.SetTemplate(t)
	return manager, err
}

func (v *VirtManager) SetTemplate(t provider.Template) error {
	v.template = &t
	v.accessPort = 22
	baeImage := fmt.Sprintf("%s/%s", v.home, path.Base(v.template.Image))
	if _, err := os.Stat(baeImage); errors.Is(err, os.ErrNotExist) {
		log.Info("start download provider.Template")

		// download file and rename
		err = utils.Download(v.template.Image, baeImage)
		if err != nil {
			log.Error("download provider.Template fail")
			return err
		}

	}

	if _, err := os.Stat(v.getBaseImagePath()); errors.Is(err, os.ErrNotExist) {
		if strings.HasSuffix(baeImage, ".tar.gz") {
			file, err := os.Open(baeImage)
			if err != nil {
				log.Error("download provider.Template fail")
				return err
			}
			return utils.UnTar(file, v.home)
		}
	}

	return nil
}

func (v *VirtManager) getCopyDiskFile(name string) string {
	return fmt.Sprintf("%s/orders/%s_%s", v.home, name, v.getBaseImageName())
}

func (v *VirtManager) getBaseImageName() string {
	imageName := path.Base(v.template.Image)
	if strings.HasSuffix(imageName, ".tar.gz") {
		imageName = strings.ReplaceAll(imageName, ".tar.gz", "")
	}
	return imageName
}

func (v *VirtManager) getBaseImagePath() string {
	return v.home + "/" + v.getBaseImageName()
}

// Create create
func (v *VirtManager) Create(name string) (string, error) {
	log.Info("start the virtual machine")
	//xml, err := v.newXml()
	//if err != nil {
	//	return err
	//}
	//log.Info(xml)
	//domain, err := v.conn.DomainDefineXML(xml)
	//defer domain.Free()

	if _, err := os.Stat(v.getCopyDiskFile(name)); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(path.Dir(v.getCopyDiskFile(name)), os.ModePerm)

		fmt.Println("cp", v.getBaseImagePath(), v.getCopyDiskFile(name))
		cmd := exec.Command("cp", v.getBaseImagePath(), v.getCopyDiskFile(name))
		err := cmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
		}
	}

	fmt.Println("virt-install", "--virt-type", "kvm", "--name", name, "--vcpus", fmt.Sprintf("%d", v.template.Cpu), "--ram", fmt.Sprintf("%d", v.template.Memory<<10), "--disk", fmt.Sprintf("path=%s", v.getCopyDiskFile(name)), "--network", "network=default", "--graphics", "vnc,listen=0.0.0.0", "--noautoconsole", "--boot", "hd")
	cmd := exec.Command("virt-install", "--virt-type", "kvm", "--name", name, "--vcpus", fmt.Sprintf("%d", v.template.Cpu), "--ram", fmt.Sprintf("%d", v.template.Memory<<10), "--disk", fmt.Sprintf("path=%s", v.getCopyDiskFile(name)), "--network", "network=default", "--graphics", "vnc,listen=0.0.0.0", "--noautoconsole", "--boot", "hd")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
	}
	return name, err
}

// Start start the virtual machine
func (v *VirtManager) Start(name string) error {

	d, err := v.conn.LookupDomainByName(name)
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(d)

	if err != nil {
		return err
	}

	if active, err := d.IsActive(); active {
		return err
	}

	return d.Create()
}

// CreateAndStart create and start
func (v *VirtManager) CreateAndStart(name string) (string, error) {

	id, err := v.Create(name)
	if err != nil {
		return "", err
	}

	err = v.Start(name)
	return id, err
}

func (v *VirtManager) CreateAndStartAndInjectionPublicKey(name, publicKey string) (string, error) {

	// injection public key
	err := v.InjectionPublicKey(name, publicKey)
	if err != nil {
		return "", err
	}

	// create a virtual machine
	id, err := v.CreateAndStart(name)
	if err != nil {
		return id, err
	}

	// wait for the virtual machine to start successfully
	for {
		status, err := v.Status(name)
		if err != nil {
			return "", err
		}
		if status.IsRunning() {
			break
		}
		time.Sleep(time.Second * 3)
	}

	log.Info("processing order complete")
	return id, err
}

// Stop shutdown the virtual machine
func (v *VirtManager) Stop(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.Shutdown()
}

// Reboot restart the virtual machine
func (v *VirtManager) Reboot(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
}

// Shutdown shutdown the virtual machine
func (v *VirtManager) Shutdown(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_ACPI_POWER_BTN)
}

// Destroy destroy the virtual machine
func (v *VirtManager) Destroy(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(d)
	return d.Destroy()
}

// InjectionPublicKey injection publickey (The timing of vm implementation and docker implementation are different)
func (v *VirtManager) InjectionPublicKey(name, publicKey string) error {

	publicKeyFileName := fmt.Sprintf("/tmp/%s.pub", utils.GetRandomString(10))

	err := fileutils.CreateIfNotExists(publicKeyFileName, false)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(publicKeyFileName, os.O_WRONLY, os.ModePerm)
	_, err = file.WriteString(publicKey)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	fmt.Println("virt-customize", "-a", v.getCopyDiskFile(name), "--ssh-inject", fmt.Sprintf("root:file:%s", publicKeyFileName))
	cmd := exec.Command("virt-customize", "-a", v.getCopyDiskFile(name), "--ssh-inject", fmt.Sprintf("root:file:%s", publicKeyFileName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		fmt.Println("===output===")
		fmt.Println(string(output))
		fmt.Println("============")

	}
	_ = os.Remove(publicKeyFileName)
	return err
}

// Status View status
func (v *VirtManager) Status(name string) (*provider.Status, error) {
	dom, err := v.conn.LookupDomainByName(name)
	if err != nil {
		fmt.Println("err", err)
		return nil, nil
	}

	info, err := dom.GetInfo()
	state := info.State

	//DOMAIN_NOSTATE     = DomainState(C.VIR_DOMAIN_NOSTATE)
	//	DOMAIN_RUNNING     = DomainState(C.VIR_DOMAIN_RUNNING)
	//	DOMAIN_BLOCKED     = DomainState(C.VIR_DOMAIN_BLOCKED)
	//	DOMAIN_PAUSED      = DomainState(C.VIR_DOMAIN_PAUSED)
	//	DOMAIN_SHUTDOWN    = DomainState(C.VIR_DOMAIN_SHUTDOWN)
	//	DOMAIN_CRASHED     = DomainState(C.VIR_DOMAIN_CRASHED)
	//	DOMAIN_PMSUSPENDED = DomainState(C.VIR_DOMAIN_PMSUSPENDED)
	//	DOMAIN_SHUTOFF     = DomainState(C.VIR_DOMAIN_SHUTOFF)

	var status int
	id, err := dom.GetID()

	switch state {
	case libvirt.DOMAIN_NOSTATE:
		status = 0
		break
	case libvirt.DOMAIN_BLOCKED:
		status = 3
		break
	case libvirt.DOMAIN_RUNNING:
		status = 1
		break
	case libvirt.DOMAIN_PMSUSPENDED:
		status = 3
		break
	case libvirt.DOMAIN_PAUSED:
		status = 2
		break
	case libvirt.DOMAIN_SHUTOFF:
	case libvirt.DOMAIN_SHUTDOWN:
	case libvirt.DOMAIN_CRASHED:
		status = 0
		break
	}

	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			log.Error("free libvirt.Domain fail")
		}
	}(dom)
	return &provider.Status{
		Status: status,
		Id:     strconv.Itoa(int(id)),
	}, err
}

// GetIp get runtime ip
func (v *VirtManager) GetIp(name string) (string, error) {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return "", err
	}
	var dis []libvirt.DomainInterface
	failTimes := 0
	for {
		if failTimes > 180 {
			return "", err
		}
		dis, err = d.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
		if err != nil || len(dis) == 0 {
			failTimes++
			time.Sleep(time.Second)
			fmt.Println("fail time is :", failTimes)
			continue
		} else {
			fmt.Println("success time is :", failTimes)
			break
		}
	}

	for _, di := range dis {
		if len(di.Addrs) == 0 {
			continue
		}

		for _, ipAddress := range di.Addrs {
			return ipAddress.Addr, nil
		}
	}
	return "", errors.New("cannot get vm ip address")
}

// GetAccessPort get runtime port
func (v *VirtManager) GetAccessPort(name string) int {
	return v.accessPort
}

func helpUint(x uint) *uint { return &x }
