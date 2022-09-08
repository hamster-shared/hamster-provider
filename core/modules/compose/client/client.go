package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
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

func newComposeCommand() (*cobra.Command, error) {
	dockerCli, err := newDockerCli()
	if err != nil {
		return nil, err
	}
	lazyInit := api.NewServiceProxy()
	lazyInit.WithService(compose.NewComposeService(dockerCli))
	composeCommand := commands.RootCommand(dockerCli, lazyInit)
	return composeCommand, nil
}

func Compose(ctx context.Context, args []string) error {
	composeCommand, err := newComposeCommand()
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
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	composeCmd, err := newComposeCommand()
	if err != nil {
		return "", fmt.Errorf("failed to create compose command: %v", err)
	}
	composeCmd.SetArgs(args)
	err = composeCmd.ExecuteContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute compose command: %v", err)
	}
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			fmt.Printf("failed to read from pipe: %v", err)
			return
		}
		outC <- buf.String()
	}()
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close pipe: %v", err)
	}
	os.Stdout = old
	out := <-outC
	return out, nil
}
