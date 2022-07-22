package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManager(t *testing.T) {
	// init test config file
	configPath := filepath.Join(DefaultConfigDir(), "test_config.yaml")
	f, err := os.Create(configPath)
	f.Close()
	assert.NoError(t, err)

	// test get config
	cm := NewConfigManagerWithPath(configPath)
	c, err := cm.GetConfig()
	if err != nil {
		t.Error(err)
	}

	// test save config
	c.ApiPort = 1999
	err = cm.Save(c)
	if err != nil {
		t.Error(err)
	}

	// clean
	err = os.Remove(configPath)
	assert.NoError(t, err)
}

func TestConfigManager_Migrate(t *testing.T) {
	cm := NewConfigManager()
	err := cm.Migrate()
	assert.NoError(t, err)
}

func TestConfigManager_GetConfig(t *testing.T) {
	cm := NewConfigManager()
	_, err := cm.GetConfig()
	assert.NoError(t, err)
}
