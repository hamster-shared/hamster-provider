package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/hamster-shared/hamster-provider/log"
)

var packageLock sync.Mutex

const (
	CONFIG_DIR_NAME         = ".hamster-provider"
	CONFIG_DEFAULT_FILENAME = "config"
	SWARM_KEY               = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"
)

// Config  config parameter
type Config struct {
	ApiPort      int          `json:"apiPort"`      // API port number
	Identity     Identity     `json:"identity"`     // p2p id
	Keys         []PublicKey  `json:"keys"`         // public key list
	Bootstraps   []string     `json:"bootstraps"`   // local nodes's bootstrap peer addresses
	LinkApi      string       `json:"linkApi"`      // centralized reporting address
	ChainApi     string       `json:"chainApi"`     // blockchain address
	SeedOrPhrase string       `json:"seedOrPhrase"` // blockchain account seed or mnemonic
	Vm           VmOption     `json:"vm"`           // theoretical environment config
	ChainRegInfo ChainRegInfo `json:"chainRegInfo"` // chain registration information
	ConfigFlag   ConfigFlag   `json:"configFlag"`
	PublicIP     string       `json:"publicIP"`
}

type ConfigFlag string

const (
	DONE ConfigFlag = "done"
	NONE ConfigFlag = "none"
)

// VmOption vm configuration information
type VmOption struct {
	Cpu        uint64 `json:"cpu"`
	Mem        uint64 `json:"mem"`
	Disk       uint64 `json:"disk"`
	System     string `json:"system"`
	Image      string `json:"image"`
	AccessPort int    `json:"accessPort"`
	// virtualization type,docker/kvm
	Type string `json:"type"`
}

// Identity p2p identity token structure
type Identity struct {
	PeerID   string
	PrivKey  string `json:",omitempty"`
	SwarmKey string `json:"swarm_key"`
}

// PublicKey public key information
type PublicKey struct {
	Key string `json:"key"`
}

type ConfigManager struct {
	configPath string
}

type ChainRegInfo struct {
	ResourceIndex   uint64 `json:"resourceIndex"`
	OrderIndex      uint64 `json:"orderIndex"`
	AgreementIndex  uint64 `json:"agreementIndex"`
	RenewOrderIndex uint64 `json:"renewOrderIndex"`
	Working         string `json:"working"`
	Price           uint64 `json:"price"`
	AccountAddress  string `json:"accountAddress"`
	DeployType      uint32 `json:"deployType"`
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
	return strings.Join(
		[]string{DefaultConfigDir(), CONFIG_DEFAULT_FILENAME},
		string(os.PathSeparator),
	)
}

func DefaultConfigDir() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.GetLogger().Error(err)
		return CONFIG_DIR_NAME + "."
	}
	dir := strings.Join([]string{userHomeDir, CONFIG_DIR_NAME}, string(os.PathSeparator))
	return dir
}

func (cm *ConfigManager) GetConfig() (*Config, error) {
	packageLock.Lock()
	defer packageLock.Unlock()

	var cfg Config
	f, err := os.Open(cm.configPath)
	if err != nil {
		return nil, errors.New(
			"hamster-provider not initialized, please run `hamster-provider config init`",
		)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}

	return &cfg, nil
}

func (cm *ConfigManager) Save(config *Config) error {
	packageLock.Lock()
	defer packageLock.Unlock()
	f, err := os.OpenFile(cm.configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0766)
	if err != nil {
		return errors.New(
			"hamster-provider not initialized, please run `hamster-provider config init`",
		)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		return err
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
