package chain

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRegisterResource(t *testing.T) {

	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)

	assert.NoError(t, err)

	r := ResourceInfo{
		PeerId:     "peerId",
		Cpu:        8,
		Memory:     16,
		System:     "ubuntu",
		VmType:     "docker",
		CpuModel:   "E5-2670v3",
		Creator:    "",
		ExpireTime: time.Now().Add(time.Hour * 1000),
		User:       "root",
		Status:     0,
	}
	err = cc.RegisterResource(r)
	assert.NoError(t, err)
}

func TestExecOrder(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)
	err = cc.OrderExec(22)
	assert.NoError(t, err)
}

func TestRemoveResource(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)

	assert.NoError(t, err)

	err = cc.OrderExec(0)
	//assert.Error(t, err)
	assert.NoError(t, err)
}

func TestKeyPair(t *testing.T) {
	seedOrPhrase := "gesture region bundle fix hazard assume ozone yellow chronic camp walnut cactus"
	kp, err := signature.KeyringPairFromSecret(seedOrPhrase, 42)
	assert.NoError(t, err)
	fmt.Println(kp.URI)
	fmt.Println(kp.Address)
	fmt.Printf(string(kp.PublicKey))
}

func TestEvent(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)
	evt, err := cc.GetEvent(189)
	fmt.Println(len(evt.ResourceOrder_ReNewOrderSuccess))
}

func TestDuration(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)

	duration := cc.CalculateInstanceOverdue(6)
	fmt.Println(duration.String())
}

func TestResource(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)

	resource, err := cc.GetResource(41)
	assert.NoError(t, err)
	fmt.Println(resource)
}

func TestChainClient_ProcessApplyFreeResource(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)

	err = cc.ProcessApplyFreeResource(3, "abcd")
	assert.NoError(t, err)
}

func TestGetAgreementIndex(t *testing.T) {

}
