package pk

import (
	"errors"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	cm *config.ConfigManager
}

func NewManager(cm *config.ConfigManager) *Manager {
	return &Manager{
		cm: cm,
	}
}

// AddPublicKey add public key
func (p *Manager) AddPublicKey(publicKey string) error {

	// read config file
	c, err := p.cm.GetConfig()
	if err != nil {
		logrus.Println(err)
		return err
	}

	// add public key in keys
	keys := append(c.Keys, config.PublicKey{Key: publicKey})
	c.Keys = keys

	// the current public key is
	logrus.Printf("the current public key is:%s", c.Keys)

	// save config locally
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	} else {
		logrus.Printf("local config update public key successfully")
	}

	return nil
}

// DeletePublicKey delete public key
func (p *Manager) DeletePublicKey(publicKey string) error {
	// read config file
	c, err := p.cm.GetConfig()
	if err != nil {
		logrus.Println(err)
		return err
	}

	ifExist := false
	var res []config.PublicKey

	for _, k := range c.Keys {
		if k.Key != publicKey {
			res = append(res, k)
		} else {
			logrus.Println("delete public key successfully")
			ifExist = true
		}
	}

	c.Keys = res
	if ifExist == false {
		logrus.Println("The public key to delete does not exist, please re-enter")
		return errors.New("The public key to delete does not exist, please re-enter")
	}

	// save config
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	}

	// the current public key is
	logrus.Printf("the current public key is:%s", c.Keys)

	return nil
}

// QueryPublicKey check if the public key exists
func (p *Manager) QueryPublicKey(publicKey string) (bool, error) {
	// read config file
	c, err := p.cm.GetConfig()
	if err != nil {
		logrus.Println(err)
		return false, err
	}

	ifExist := false
	for _, k := range c.Keys {
		if k.Key == publicKey {
			ifExist = true
		}
	}

	if ifExist {
		logrus.Printf("the queried public key exists")
		return true, nil
	} else {
		logrus.Printf("the queried public key does not exist")
		return false, nil
	}

}

// ClearPublicKey clear public key list
func (p *Manager) ClearPublicKey() error {
	var res []config.PublicKey
	c, err := p.cm.GetConfig()
	c.Keys = res

	// save config
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	}

	return nil
}
