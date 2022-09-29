package polygon

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_generateRequiredFiles(t *testing.T) {
	err := generateRequiredFiles()
	assert.NilError(t, err)
}

func Test_templateInstance(t *testing.T) {
	err := templateInstance(DeployParams{})
	assert.NilError(t, err)
}
