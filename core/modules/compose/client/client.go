package client

import (
	"context"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
)

var (
	once           sync.Once
	composeCommand *cobra.Command
)

func newDockerCli() (*command.DockerCli, error) {
	initOpt := command.WithInitializeClient(func(dockerCli *command.DockerCli) (client.APIClient, error) {
		return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	})
	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	err = cli.Initialize(flags.NewClientOptions(), initOpt)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func initComposeCommand() error {
	var err error
	once.Do(func() {
		dockerCli, err := newDockerCli()
		if err != nil {
			return
		}
		lazyInit := api.NewServiceProxy()
		lazyInit.WithService(compose.NewComposeService(dockerCli))
		composeCommand = commands.RootCommand(dockerCli, lazyInit)
	})
	return err
}

func Compose(ctx context.Context, args []string) error {
	err := initComposeCommand()
	if err != nil {
		return err
	}
	composeCommand.SetArgs(args)
	err = composeCommand.ExecuteContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RunCompose(ctx context.Context, args ...string) (output string, err error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err = initComposeCommand()
	if err != nil {
		return "", err
	}
	composeCommand.SetArgs(args)
	err = composeCommand.ExecuteContext(ctx)
	os.Stdout = old
	w.Close()
	out, _ := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
