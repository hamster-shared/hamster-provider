package thegraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDockerComposeStatus(t *testing.T) {
	_, err := GetDockerComposeStatus()
	assert.NoError(t, err)
}
