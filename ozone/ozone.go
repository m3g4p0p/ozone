package ozone

import (
	"context"
	"errors"
	"iter"

	"github.com/ollama/ollama/api"
)

var errStop = errors.New("stop")

type Chatter interface {
	Chat(ctx context.Context, req *api.ChatRequest, fn api.ChatResponseFunc) error
}

type RunOptions struct {
	Client  Chatter
	History []api.Message
}

func (r *RunOptions) resolveClient() (Chatter, error) {
	if r.Client != nil {
		return r.Client, nil
	}
	return api.ClientFromEnvironment()
}

type RunResult struct {
	ctx    context.Context
	client Chatter
	req    *api.ChatRequest
}

func (r *RunResult) Stream() iter.Seq2[api.ChatResponse, error] {
	return func(yield func(api.ChatResponse, error) bool) {
		err := r.client.Chat(r.ctx, r.req, func(cr api.ChatResponse) error {
			if !yield(cr, nil) {
				return errStop
			}
			return nil
		})
		if err != nil && !errors.Is(err, errStop) {
			yield(api.ChatResponse{}, err)
		}
	}
}

type Agent struct {
	Name   string
	Model  string
	System string
}

func (a *Agent) Run(ctx context.Context, input string, options *RunOptions) (*RunResult, error) {
	if options == nil {
		options = &RunOptions{}
	}

	client, err := options.resolveClient()
	if err != nil {
		return nil, err
	}

	req := &api.ChatRequest{
		Model: a.Model,
	}

	if a.System != "" {
		req.Messages = append(req.Messages, api.Message{
			Role:    "system",
			Content: a.System,
		})
	}

	req.Messages = append(
		req.Messages,
		options.History...,
	)

	req.Messages = append(req.Messages, api.Message{
		Role:    "user",
		Content: input,
	})

	res := &RunResult{
		client: client,
		ctx:    ctx,
		req:    req,
	}

	return res, nil
}
