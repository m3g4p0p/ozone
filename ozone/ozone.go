package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Chatter interface {
	Chat(ctx context.Context, req *api.ChatRequest, fn api.ChatResponseFunc) error
}

type Agent struct {
	Name string
}

type RunOptions struct {
	Client Chatter
}

func (a *Agent) Run(ctx context.Context, input string, options *RunOptions) error {
	return nil
}
