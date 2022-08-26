package corehttp

import (
	"github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	// some test
}

func TestPortWork(t *testing.T) {
	ctx := context.CoreContext{Cm: config.NewConfigManager()}
	err := StartApi(&ctx)
	assert.NoError(t, err)
}
