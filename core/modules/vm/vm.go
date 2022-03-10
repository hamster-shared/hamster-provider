package vm

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

// Manager 虚拟化接口
type Manager interface {
	// SetTemplate 配置模板
	SetTemplate(t Template) error
	// Create 创建
	Create(name string) (string, error)
	// Start 启动虚拟机
	Start(name string) error
	// CreateAndStart 创建并启动
	CreateAndStart(name string) (string, error)
	// 启动并注入公钥
	CreateAndStartAndInjectionPublicKey(name string, publicKey string) (string, error)
	// Stop 关闭虚拟机
	Stop(name string) error
	// Reboot 重启虚拟机
	Reboot(name string) error
	// Shutdown 关闭虚拟机
	Shutdown(name string) error
	// Destroy 销毁虚拟机
	Destroy(name string) error
	// InjectionPublicKey 注入公玥 (vm实现和docker实现时机不同)
	InjectionPublicKey(name string, publicKey string) error
	// Status 查看状态
	Status(name string) (*Status, error)

	// GetIp 获取运行时ip
	GetIp(name string) (string, error)
	// GetAccessPort 获取运行时端口
	GetAccessPort(name string) int
}

type Status struct {
	id string
	// status 状态 0: 关闭,1: running , 2：其他
	status int
}

// IsRunning 是否正在运行
func (s *Status) IsRunning() bool {
	return s.status == 1
}

type Template struct {
	Cpu, Memory, Disk uint64
	System            string
	PublicKey         string
	Image             string
	AccessPort        int
}
