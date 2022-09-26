package aptos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateRequiredFiles(t *testing.T) {
	err := generateRequiredFiles()
	assert.NoError(t, err)
}

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
