package thegraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplateInstance(t *testing.T) {
	var data DeployParams
	data.EthereumNetwork = "rinkeby"
	data.EthereumUrl = "https://rinkeby.infura.io/v3/af7a79eb36f64e609b5dda130cd62946"
	data.IndexerAddress = "0x9438BbE4E7AF1ec6b13f75ECd1f53391506A12DF"
	data.Mnemonic = "please output text solve glare exit divert boil nerve eagle attack turkey"
	data.NodeEthereumUrl = "rinkeby:https://rinkeby.infura.io/v3/af7a79eb36f64e609b5dda130cd62946"
	err := templateInstance(data)

	assert.NoError(t, err)
}

func TestStartDockerCompose(t *testing.T) {
	err := startDockerCompose()
	assert.NoError(t, err)
}

func Test_getDockerComposeStatus(t *testing.T) {
	_, err := getDockerComposeStatus()
	assert.NoError(t, err)
}
