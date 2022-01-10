package vm

// dependency
// yum install libvirt-devel gcc
// brew install libvirt

import (
	"errors"
	"fmt"
	"github.com/docker/docker/pkg/fileutils"
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
	template   *Template
	accessPort int
}

// NewVirtManager create virtManager
func NewVirtManager() (*VirtManager, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	homedir, err := os.UserHomeDir()
	return &VirtManager{
		conn: conn,
		home: homedir + "/.ttchain-compute-provider",
	}, err
}

func (v *VirtManager) SetTemplate(t Template) {
	v.template = &t
	v.accessPort = 22
	baeImage := fmt.Sprintf("%s/%s", v.home, path.Base(v.template.Image))
	if _, err := os.Stat(baeImage); errors.Is(err, os.ErrNotExist) {
		log.Info("start download template")

		// download file and rename
		err = utils.Download(v.template.Image, baeImage)
		if err != nil {
			log.Error("download template fail")
			return
		}

	}

	if _, err := os.Stat(v.getBaseImagePath()); errors.Is(err, os.ErrNotExist) {
		if strings.HasSuffix(baeImage, ".tar.gz") {
			file, err := os.Open(baeImage)
			if err != nil {
				log.Error("download template fail")
				return
			}
			err = utils.UnTar(file, v.home)
		}
	}

	if _, err := os.Stat(v.getCopyDiskFile()); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(path.Dir(v.getCopyDiskFile()), os.ModePerm)

		fmt.Println("cp", v.getBaseImagePath(), v.getCopyDiskFile())
		cmd := exec.Command("cp", v.getBaseImagePath(), v.getCopyDiskFile())
		err := cmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
		}
	}
}

//func (v *VirtManager) newXml() (string, error) {
//
//	if v.template == nil {
//		return "", errors.New("template not set,please set template first!")
//	}
//
//	domcfg := &libvirtxml.Domain{
//		Type: "kvm",
//		Name: v.template.Name,
//		Memory: &libvirtxml.DomainMemory{
//			Unit:  "GiB",
//			Value: uint(v.template.Memory),
//		},
//		CurrentMemory: &libvirtxml.DomainCurrentMemory{
//			Unit:  "GiB",
//			Value: uint(v.template.Memory),
//		},
//		VCPU: &libvirtxml.DomainVCPU{
//			Placement: "static",
//			Value:     uint(v.template.Cpu),
//		},
//		OS: &libvirtxml.DomainOS{
//			Type: &libvirtxml.DomainOSType{
//				Arch:    "x86_64",
//				Machine: "pc-i440fx-rhel7.0.0",
//				Type:    "hvm",
//			},
//			BootDevices: []libvirtxml.DomainBootDevice{
//				{
//					Dev: "hd",
//				},
//			},
//		},
//		Features: &libvirtxml.DomainFeatureList{
//			APIC: &libvirtxml.DomainFeatureAPIC{},
//			ACPI: &libvirtxml.DomainFeature{},
//		},
//		CPU: &libvirtxml.DomainCPU{
//			Mode:  "custom",
//			Match: "exact",
//			Check: "partial",
//			Model: &libvirtxml.DomainCPUModel{
//				Fallback: "allow",
//				Value:    "Haswell-noTSX-IBRS",
//			},
//			Features: []libvirtxml.DomainCPUFeature{
//				{
//					Policy: "require",
//					Name:   "md-clear",
//				},
//				{
//					Policy: "require",
//					Name:   "spec-ctrl",
//				},
//				{
//					Policy: "require",
//					Name:   "ssbd",
//				},
//			},
//		},
//		OnCrash:    "destroy",
//		OnReboot:   "restart",
//		OnPoweroff: "destroy",
//		Devices: &libvirtxml.DomainDeviceList{
//			Emulator: "/usr/libexec/qemu-kvm",
//			Disks: []libvirtxml.DomainDisk{
//				{
//					Device: "disk",
//					Driver: &libvirtxml.DomainDiskDriver{
//						Name: "qemu",
//						Type: "qcow2",
//					},
//					Source: &libvirtxml.DomainDiskSource{
//						File: &libvirtxml.DomainDiskSourceFile{
//							File: v.image,
//						},
//					},
//					Target: &libvirtxml.DomainDiskTarget{
//						Dev: "hda",
//						Bus: "ide",
//					},
//					Address: &libvirtxml.DomainAddress{
//						Drive: &libvirtxml.DomainAddressDrive{
//							Controller: helpUint(0),
//							Bus:        helpUint(0),
//							Target:     helpUint(0),
//							Unit:       helpUint(0),
//						},
//					},
//				},
//			},
//			Controllers: []libvirtxml.DomainController{
//				{
//					Type:  "usb",
//					Index: helpUint(0),
//					Model: "ich9-ehci1",
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(4),
//							Function: helpUint(7),
//						},
//					},
//				},
//				{
//					Type:  "usb",
//					Index: helpUint(0),
//					Model: "ich9-uhci1",
//					USB: &libvirtxml.DomainControllerUSB{
//						Master: &libvirtxml.DomainControllerUSBMaster{
//							StartPort: 0,
//						},
//					},
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:        helpUint(0),
//							Bus:           helpUint(0),
//							Slot:          helpUint(4),
//							Function:      helpUint(0),
//							MultiFunction: "on",
//						},
//					},
//				},
//				{
//					Type:  "usb",
//					Index: helpUint(0),
//					Model: "ich9-uhci2",
//					USB: &libvirtxml.DomainControllerUSB{
//						Master: &libvirtxml.DomainControllerUSBMaster{
//							StartPort: 2,
//						},
//					},
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(4),
//							Function: helpUint(1),
//						},
//					},
//				},
//				{
//					Type:  "usb",
//					Index: helpUint(0),
//					Model: "ich9-uhci3",
//					USB: &libvirtxml.DomainControllerUSB{
//						Master: &libvirtxml.DomainControllerUSBMaster{
//							StartPort: 4,
//						},
//					},
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(4),
//							Function: helpUint(2),
//						},
//					},
//				},
//				{
//					Type:  "pci",
//					Index: helpUint(0),
//					Model: "pci-root",
//				},
//				{
//					Type:  "ide",
//					Index: helpUint(0),
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(1),
//							Function: helpUint(1),
//						},
//					},
//				},
//			},
//			Interfaces: []libvirtxml.DomainInterface{
//				{
//					Source: &libvirtxml.DomainInterfaceSource{
//						Network: &libvirtxml.DomainInterfaceSourceNetwork{
//							Network: "default",
//						},
//					},
//					Model: &libvirtxml.DomainInterfaceModel{
//						Type: "e1000",
//					},
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(3),
//							Function: helpUint(0),
//						},
//					},
//				},
//			},
//			Serials: []libvirtxml.DomainSerial{
//				{
//					Target: &libvirtxml.DomainSerialTarget{
//						Type: "isa-serial",
//						Port: helpUint(0),
//						Model: &libvirtxml.DomainSerialTargetModel{
//							Name: "isa-serial",
//						},
//					},
//				},
//			},
//			Consoles: []libvirtxml.DomainConsole{
//				{
//					Target: &libvirtxml.DomainConsoleTarget{
//						Type: "serial",
//						Port: helpUint(0),
//					},
//				},
//			},
//			Inputs: []libvirtxml.DomainInput{
//				{
//					Type: "tablet",
//					Bus:  "usb",
//					Address: &libvirtxml.DomainAddress{
//						USB: &libvirtxml.DomainAddressUSB{
//							Bus:  helpUint(0),
//							Port: "1",
//						},
//					},
//				},
//				{
//					Type: "mouse",
//					Bus:  "ps2",
//				},
//				{
//					Type: "keyboard",
//					Bus:  "ps2",
//				},
//			},
//			Videos: []libvirtxml.DomainVideo{
//				{
//					Model: libvirtxml.DomainVideoModel{
//						Type:    "vga",
//						VRam:    16384,
//						Heads:   1,
//						Primary: "yes",
//					},
//					Address: &libvirtxml.DomainAddress{
//						PCI: &libvirtxml.DomainAddressPCI{
//							Domain:   helpUint(0),
//							Bus:      helpUint(0),
//							Slot:     helpUint(2),
//							Function: helpUint(0),
//						},
//					},
//				},
//			},
//			MemBalloon: &libvirtxml.DomainMemBalloon{
//				Model: "virtio",
//				Address: &libvirtxml.DomainAddress{
//					PCI: &libvirtxml.DomainAddressPCI{
//						Domain:   helpUint(0),
//						Bus:      helpUint(0),
//						Slot:     helpUint(5),
//						Function: helpUint(0),
//					},
//				},
//			},
//		},
//	}
//
//	xml, err := domcfg.Marshal()
//	if err != nil {
//		return "", err
//	}
//
//	return xml, nil
//}

func (v *VirtManager) getCopyDiskFile() string {
	return fmt.Sprintf("%s/orders/%s_%s", v.home, v.template.Name, v.getBaseImageName())
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

// Create 创建
func (v *VirtManager) Create() error {
	log.Info("start the virtual machine")
	//xml, err := v.newXml()
	//if err != nil {
	//	return err
	//}
	//log.Info(xml)
	//domain, err := v.conn.DomainDefineXML(xml)
	//defer domain.Free()

	fmt.Println("virt-install", "--virt-type", "kvm", "--name", v.template.Name, "--vcpus", fmt.Sprintf("%d", v.template.Cpu), "--ram", fmt.Sprintf("%d", v.template.Memory<<10), "--disk", fmt.Sprintf("path=%s", v.getCopyDiskFile()), "--network", "network=default", "--graphics", "vnc,listen=0.0.0.0", "--noautoconsole", "--boot", "hd")
	cmd := exec.Command("virt-install", "--virt-type", "kvm", "--name", v.template.Name, "--vcpus", fmt.Sprintf("%d", v.template.Cpu), "--ram", fmt.Sprintf("%d", v.template.Memory<<10), "--disk", fmt.Sprintf("path=%s", v.getCopyDiskFile()), "--network", "network=default", "--graphics", "vnc,listen=0.0.0.0", "--noautoconsole", "--boot", "hd")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
	}
	return err
}

// Start start the virtual machine
func (v *VirtManager) Start() error {

	d, err := v.conn.LookupDomainByName(v.template.Name)
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
func (v *VirtManager) CreateAndStart() error {

	err := v.Create()
	if err != nil {
		return err
	}

	return v.Start()
}

func (v *VirtManager) CreateAndStartAndInjectionPublicKey(publicKey string) error {

	// injection public key
	err := v.InjectionPublicKey(publicKey)
	if err != nil {
		return err
	}

	// create a virtual machine
	err = v.CreateAndStart()
	if err != nil {
		return err
	}

	// wait for the virtual machine to start successfully
	for {
		status, err := v.Status()
		if err != nil {
			return err
		}
		if status.IsRunning() {
			break
		}
		time.Sleep(time.Second * 3)
	}

	log.Info("processing order complete")
	return nil
}

// Stop shut down the virtual machine
func (v *VirtManager) Stop() error {
	d, err := v.conn.LookupDomainByName(v.template.Name)
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
func (v *VirtManager) Reboot() error {
	d, err := v.conn.LookupDomainByName(v.template.Name)
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

// Shutdown shut down the virtual machine
func (v *VirtManager) Shutdown() error {
	d, err := v.conn.LookupDomainByName(v.template.Name)
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
func (v *VirtManager) Destroy() error {
	d, err := v.conn.LookupDomainByName(v.template.Name)
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
func (v *VirtManager) InjectionPublicKey(publicKey string) error {

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
	fmt.Println("virt-customize", "-a", v.getCopyDiskFile(), "--ssh-inject", fmt.Sprintf("root:file:%s", publicKeyFileName))
	cmd := exec.Command("virt-customize", "-a", v.getCopyDiskFile(), "--ssh-inject", fmt.Sprintf("root:file:%s", publicKeyFileName))
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
func (v *VirtManager) Status() (*Status, error) {
	dom, err := v.conn.LookupDomainByName(v.template.Name)
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
	return &Status{
		status: status,
		id:     strconv.Itoa(int(id)),
	}, err
}

// GetIp get runtime ip
func (v *VirtManager) GetIp() (string, error) {
	d, err := v.conn.LookupDomainByName(v.template.Name)
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
func (v *VirtManager) GetAccessPort() int {
	return v.accessPort
}

func helpUint(x uint) *uint { return &x }
