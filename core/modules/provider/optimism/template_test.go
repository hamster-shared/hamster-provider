package optimism

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_generateRequiredFiles(t *testing.T) {
	err := generateRequiredFiles()
	assert.NilError(t, err)
}

func Test_getAllFilenames(t *testing.T) {
	filenames, err := getAllFilenames(&allFiles, "templates")
	assert.NilError(t, err)
	fmt.Println(filenames)
}

func Test_templateInstance(t *testing.T) {
	err := templateInstance(DeployParams{})
	assert.NilError(t, err)
}
