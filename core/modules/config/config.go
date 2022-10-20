package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/hamster-shared/hamster-provider/log"
)

var packageLock sync.Mutex

const (
	CONFIG_DIR_NAME         = ".hamster-provider"
	CONFIG_DEFAULT_FILENAME = "config.yaml"
	SWARM_KEY               = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"
)

// Config  config parameter
type Config struct {
	ApiPort       int           `json:"apiPort" yaml:"apiPort"`           // API port number
	Identity      Identity      `json:"identity" yaml:"identity"`         // p2p id
	Keys          []PublicKey   `json:"keys" yaml:"keys"`                 // public key list
	Bootstraps    []string      `json:"bootstraps" yaml:"bootstraps"`     // local nodes's bootstrap peer addresses
	LinkApi       string        `json:"linkApi" yaml:"linkApi"`           // centralized reporting address
	ChainApi      string        `json:"chainApi" yaml:"chainApi"`         // blockchain address
	SeedOrPhrase  string        `json:"seedOrPhrase" yaml:"seedOrPhrase"` // blockchain account seed or mnemonic
	Vm            VmOption      `json:"vm" yaml:"vm"`                     // theoretical environment config
	ChainRegInfo  ChainRegInfo  `json:"chainRegInfo" yaml:"chainRegInfo"` // chain registration information
	ConfigFlag    ConfigFlag    `json:"configFlag" yaml:"configFlag"`
	PublicIP      string        `json:"publicIP" yaml:"publicIP"`
	Specification Specification `json:"specification" yaml:"specification"`
}

type ConfigFlag string

const (
	DONE ConfigFlag = "done"
	NONE ConfigFlag = "none"
)

// VmOption vm configuration information
type VmOption struct {
	Cpu        uint64 `json:"cpu" yaml:"cpu"`
	Mem        uint64 `json:"mem" yaml:"mem"`
	Disk       uint64 `json:"disk" yaml:"disk"`
	System     string `json:"system" yaml:"system"`
	Image      string `json:"image" yaml:"image"`
	AccessPort int    `json:"accessPort" yaml:"accessPort"`
	// virtualization type,docker/kvm
	Type string `json:"type" yaml:"type"`
}

type Specification = uint32

const (
	General     Specification = 0
	Enhanced    Specification = 1
	HighRanking Specification = 2
)

// Identity p2p identity token structure
type Identity struct {
	PeerID   string `yaml:"peerID"`
	PrivKey  string `json:",omitempty" yaml:"privKey"`
	SwarmKey string `json:"swarm_key" yaml:"swarmKey"`
}

// PublicKey public key information
type PublicKey struct {
	Key string `json:"key" yaml:"key"`
}

type ConfigManager struct {
	configPath string
}

type ChainRegInfo struct {
	ResourceIndex   uint64 `json:"resourceIndex" yaml:"resourceIndex"`
	OrderIndex      uint64 `json:"orderIndex" yaml:"orderIndex"`
	AgreementIndex  uint64 `json:"agreementIndex" yaml:"agreementIndex"`
	RenewOrderIndex uint64 `json:"renewOrderIndex" yaml:"renewOrderIndex"`
	Working         string `json:"working" yaml:"working"`
	Price           uint64 `json:"price" yaml:"price"`
	AccountAddress  string `json:"accountAddress" yaml:"accountAddress"`
	DeployType      uint32 `json:"deployType" yaml:"deployType"`
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configPath: DefaultConfigPath(),
	}
}

func NewConfigManagerWithPath(path string) *ConfigManager {
	return &ConfigManager{
		configPath: path,
	}
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), CONFIG_DEFAULT_FILENAME)
}

func DefaultConfigDir() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.GetLogger().Error(err)
		return CONFIG_DIR_NAME + "."
	}
	dir := filepath.Join(userHomeDir, CONFIG_DIR_NAME)
	return dir
}

func (cm *ConfigManager) GetConfig() (*Config, error) {
	packageLock.Lock()
	defer packageLock.Unlock()

	var cfg Config
	cfgBytes, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, errors.New(
			"hamster-provider not initialized, please run `hamster-provider config init`",
		)
	}
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}
	return &cfg, nil
}

func (cm *ConfigManager) Save(config *Config) error {
	packageLock.Lock()
	defer packageLock.Unlock()
	cfgBytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %s", err)
	}
	err = os.WriteFile(cm.configPath, cfgBytes, 0766)
	if err != nil {
		return errors.New(
			"hamster-provider not initialized, please run `hamster-provider config init`",
		)
	}
	return nil
}

func CreateIdentity() (Identity, error) {
	ident := Identity{
		SwarmKey: SWARM_KEY,
	}

	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return ident, err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}
