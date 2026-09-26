package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

type RunOptions struct {
	Client  Chatter
	History MessagesProvider
}

func (r *RunOptions) resolveClient() (Chatter, error) {
	if r.Client != nil {
		return r.Client, nil
	}
	return api.ClientFromEnvironment()
}

func (r *RunOptions) resolveHistory(ctx context.Context) ([]api.Message, error) {
	if r.History == nil {
		return nil, nil
	}
	return r.History.Messages(ctx)
}
