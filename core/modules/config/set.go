package config

import (
	"fmt"

	"github.com/hamster-shared/hamster-provider/core/modules/utils"
)

// ConfigVM Configure
func (cm *ConfigManager) ConfigVM(vmOption VmOption) error {
	config, err := cm.GetConfig()
	if err != nil {
		return err
	}
	config.Vm = vmOption

	err = cm.Save(config)
	return err
}

func (cm *ConfigManager) AddBootstrap(bootstrap string) error {
	config, err := cm.GetConfig()
	if err != nil {
		return err
	}

	bootstraps := config.Bootstraps

	if utils.Contains(bootstraps, bootstrap) {
		return nil
	} else {
		config.Bootstraps = append(bootstraps, bootstrap)
	}

	return cm.Save(config)
}

func (cm *ConfigManager) RemoveBootstrap(bootstrap string) error {
	config, err := cm.GetConfig()
	if err != nil {
		return err
	}

	config.Bootstraps = utils.Remove(config.Bootstraps, bootstrap)

	return cm.Save(config)
}

func (cm *ConfigManager) SetPublicIP(publicIP string) error {
	config, err := cm.GetConfig()
	if err != nil {
		return err
	}

	config.PublicIP = publicIP

	return cm.Save(config)
}

func GetSpecifiction(cpu uint64, mem uint64) (Specification, error) {
	if cpu >= 8 && mem >= 64 {
		return HighRanking, nil
	} else if cpu >= 4 && mem >= 32 {
		return Enhanced, nil
	} else if cpu >= 2 && mem >= 8 {
		return General, nil
	} else {
		return 0, fmt.Errorf("cpu and mem is not enough, minimum cpu: 2, minimum mem: 8")
	}
}
