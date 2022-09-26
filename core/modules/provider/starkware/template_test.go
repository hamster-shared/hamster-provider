package aptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_templateInstance(t *testing.T) {
	err := templateInstance(DeployParams{})
	assert.NoError(t, err)
}

func Test_pullImages(t *testing.T) {
	err := pullImages()
	assert.NoError(t, err)
}

func Test_startDockerCompose(t *testing.T) {
	err := startDockerCompose()
	assert.NoError(t, err)
}

func Test_stopDockerCompose(t *testing.T) {
	err := stopDockerCompose()
	assert.NoError(t, err)
}

func Test_getDockerComposeStatus(t *testing.T) {
	cs, err := getDockerComposeStatus()
	assert.NoError(t, err)
	fmt.Println(cs)
}
