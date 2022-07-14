package client

import (
	"context"
	"testing"
)

func TestCompose(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{context.Background(), []string{"up"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Compose(tt.args.ctx, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Compose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
