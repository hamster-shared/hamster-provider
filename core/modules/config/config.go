package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

var packageLock sync.Mutex

const (
	CONFIG_DIR_NAME         = ".ttchain-compute-provider"
	CONFIG_DEFAULT_FILENAME = "config"
	SWARM_KEY               = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"
)

// Config  config参数
type Config struct {
	ApiPort      int          // API port number
	Identity     Identity     // p2p id
	Keys         []PublicKey  // public key list
	Bootstraps   []string     // local nodes's bootstrap peer addresses
	LinkApi      string       // centralized reporting address
	ChainApi     string       // blockchain address
	SeedOrPhrase string       // blockchain account seed or mnemonic
	Vm           VmOption     // theoretical environment config
	ChainRegInfo ChainRegInfo // chain registration information
}

// Identity p2p identity token structure
type Identity struct {
	PeerID   string
	PrivKey  string `json:",omitempty"`
	SwarmKey string `json:"swarm_key"`
}

// PublicKey public key information
type PublicKey struct {
	Key string
}

type ConfigManager struct {
	configPath string
}

type ChainRegInfo struct {
	ResourceIndex   uint64
	OrderIndex      uint64
	AgreementIndex  uint64
	RenewOrderIndex uint64
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
	return strings.Join([]string{DefaultConfigDir(), CONFIG_DEFAULT_FILENAME}, string(os.PathSeparator))
}

func DefaultConfigDir() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return CONFIG_DIR_NAME + "."
	}
	if err != nil {
		logrus.Error(err)
	}
	dir := strings.Join([]string{userHomeDir, CONFIG_DIR_NAME}, string(os.PathSeparator))
	if err != nil {
		logrus.Error(err)
	}
	return dir
}

func (cm *ConfigManager) GetConfig() (*Config, error) {
	packageLock.Lock()
	defer packageLock.Unlock()

	var cfg Config
	f, err := os.Open(cm.configPath)
	defer f.Close()
	if err != nil {
		return nil, errors.New("ttchain-computer-provider not initialized, please run `ttchain-compute-provider config init`")
	}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}

	return &cfg, nil
}

func (cm *ConfigManager) Save(config *Config) error {
	packageLock.Lock()
	defer packageLock.Unlock()
	f, err := os.OpenFile(cm.configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0766)
	defer f.Close()
	if err != nil {
		return errors.New("ttchain-computer-provider not initialized, please run `ttchain-compute-provider config init`")
	}
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
