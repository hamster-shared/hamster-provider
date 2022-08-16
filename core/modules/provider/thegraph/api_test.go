package thegraph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDockerComposeStatus(t *testing.T) {
	_, err := GetDockerComposeStatus()
	assert.NoError(t, err)
}

var composeFilePathName = "/home/vihv/Desktop/Desktop/docker-compose-thegraph-v0.19.yml"

func TestComposeGraphConnect(t *testing.T) {
	err := ComposeGraphConnect(composeFilePathName)
	assert.NoError(t, err)
}

func TestComposeGraphStart(t *testing.T) {
	err := ComposeGraphStart(composeFilePathName, "QmVqMeQUwvQ3XjzCYiMhRvQjRiQLGpVt8C3oHgvDi3agJ2")
	assert.NoError(t, err)
}

func TestComposeGraphRules(t *testing.T) {
	result, err := ComposeGraphRules(composeFilePathName)
	assert.NoError(t, err)
	for _, rule := range result {
		fmt.Println(rule)
	}
}
