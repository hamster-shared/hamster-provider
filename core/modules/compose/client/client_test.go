package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCompose(t *testing.T) {
	cmd := "-f /home/vihv/.hamster-provider/docker-compose.yml exec index-cli graph indexer rules get all -o json"
	cmd2 := "-f /home/vihv/.hamster-provider/docker-compose.yml exec index-cli echo hello"
	err := Compose(context.Background(), strings.Split(cmd, " "))
	if err != nil {
		t.Errorf("RunCompose error: %s", err.Error())
	}

	err = Compose(context.Background(), strings.Split(cmd2, " "))
	if err != nil {
		t.Errorf("RunCompose error: %s", err.Error())
	}
}

func TestRunCompose(t *testing.T) {
	cmd := "-f /home/vihv/.hamster-provider/docker-compose.yml exec index-cli graph indexer rules get all -o json"
	cmd2 := "-f /home/vihv/.hamster-provider/docker-compose.yml exec index-cli echo hello"
	output, err := RunCompose(context.Background(), strings.Split(cmd, " ")...)
	if err != nil {
		t.Errorf("RunCompose error: %s", err.Error())
	}
	assert.NotEmpty(t, output)

	output, err = RunCompose(context.Background(), strings.Split(cmd2, " ")...)
	if err != nil {
		t.Errorf("RunCompose error: %s", err.Error())
	}
	assert.NotEmpty(t, output)
}
