package ozone

import (
	"context"
	"errors"
	"iter"
	"sync/atomic"

	"github.com/ollama/ollama/api"
)

var errStop = errors.New("stop")

var ErrStarted = errors.New("stream already started")

type Chatter interface {
	Chat(ctx context.Context, req *api.ChatRequest, fn api.ChatResponseFunc) error
}

type MessagesProvider interface {
	Messages(ctx context.Context) ([]api.Message, error)
}

type ToolsProvider interface {
	Tools(ctx context.Context) (api.Tools, error)
}

type ToolCallHandler interface {
	Call(ctx context.Context, tc api.ToolCall) (string, error)
}

type Toolset interface {
	ToolsProvider
	ToolCallHandler
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
	client  Chatter
	ctx     context.Context
	req     *api.ChatRequest
	started atomic.Bool
}

func (r *RunResult) Stream() iter.Seq2[api.ChatResponse, error] {
	return func(yield func(api.ChatResponse, error) bool) {
		if !r.started.CompareAndSwap(false, true) {
			yield(api.ChatResponse{}, ErrStarted)
			return
		}

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
	Name    string
	Model   string
	System  string
	Toolset Toolset
}

func (a *Agent) Run(
	ctx context.Context,
	input string,
	options *RunOptions,
) (*RunResult, error) {
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

	if a.Toolset != nil {
		tools, err := a.Toolset.Tools(ctx)
		if err != nil {
			return nil, err
		}
		req.Tools = tools
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
		client: &TurnHandler{
			Client:  client,
			Handler: a.Toolset,
		},
		ctx: ctx,
		req: req,
	}

	return res, nil
}
