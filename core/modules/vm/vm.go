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

// Manager virtualized interface
type Manager interface {
	// SetTemplate  config template
	SetTemplate(t Template)
	// Create create the virtual machine
	Create() error
	// Start start the virtual machine
	Start() error
	// CreateAndStart create and start
	CreateAndStart() error
	// CreateAndStartAndInjectionPublicKey start and inject public key
	CreateAndStartAndInjectionPublicKey(publicKey string) error
	// Stop shut down the virtual machine
	Stop() error
	// Reboot restart the virtual machine
	Reboot() error
	// Shutdown shut down the virtual machine
	Shutdown() error
	// Destroy destroy the virtual machine
	Destroy() error
	// InjectionPublicKey Inject publicKey (vm implementation and docker implementation timing are different)
	InjectionPublicKey(publicKey string) error
	// Status view status
	Status() (*Status, error)

	// GetIp get runtime ip
	GetIp() (string, error)
	// GetAccessPort get runtime port
	GetAccessPort() int
}

type Status struct {
	id string
	// status status 0 off 1 running 2 other
	status int
}

// IsRunning is it running
func (s *Status) IsRunning() bool {
	return s.status == 1
}

type Template struct {
	Cpu, Memory, Dist uint64
	Name              string
	System            string
	PublicKey         string
	Image             string
}
