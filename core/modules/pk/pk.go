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

// AddPublicKey 新增公钥
func (p *Manager) AddPublicKey(publicKey string) error {

	// 读取config文件
	c, err := p.cm.GetConfig()
	if err != nil {
		logrus.Println(err)
		return err
	}

	// 在Keys中添加公钥
	keys := append(c.Keys, config.PublicKey{Key: publicKey})
	c.Keys = keys

	// 现在的公钥为
	logrus.Printf("现在的公钥为:%s", c.Keys)

	// 本地保存config
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	} else {
		logrus.Printf("本地config更新公钥成功")
	}

	return nil
}

// DeletePublicKey 删除公钥
func (p *Manager) DeletePublicKey(publicKey string) error {
	// 读取config文件
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
			logrus.Println("删除公钥成功")
			ifExist = true
		}
	}

	c.Keys = res
	if ifExist == false {
		logrus.Println("要删除的公钥不存在,请重新输入")
		return errors.New("要删除的公钥不存在,请重新输入")
	}

	// 保存config
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	}

	// 现在的公钥为
	logrus.Printf("现在的公钥为:%s", c.Keys)

	return nil
}

// QueryPublicKey 查询公钥是否存在
func (p *Manager) QueryPublicKey(publicKey string) (bool, error) {
	// 读取config文件
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
		logrus.Printf("查询的公钥存在")
		return true, nil
	} else {
		logrus.Printf("查询的公钥不存在")
		return false, nil
	}

}

// ClearPublicKey 清空公钥列表
func (p *Manager) ClearPublicKey() error {
	var res []config.PublicKey
	c, err := p.cm.GetConfig()
	c.Keys = res

	// 保存config
	err = p.cm.Save(c)
	if err != nil {
		logrus.Println(err)
		return err
	}

	return nil
}
