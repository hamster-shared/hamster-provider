package thegraph

import (
	"context"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompose(t *testing.T) {

	lazyInit := api.NewServiceProxy()

	cli, err := command.NewDockerCli()
	assert.NoError(t, err)
	lazyInit.WithService(compose.NewComposeService(cli))
	ctx := context.Background()

	p := types.Project{
		Services: []types.ServiceConfig{
			{
				Name:  "foo",
				Image: "nginx",
			},
		},
	}

	optoins := api.UpOptions{
		Create: api.CreateOptions{},
		Start: api.StartOptions{
			Project: &p,
			Wait:    false,
		},
	}
	err = lazyInit.Up(ctx, &p, optoins)
	assert.NoError(t, err)
}
