package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Agent struct {
	Name    string
	Model   string
	System  MessagesProvider
	Toolset ToolsetProvider
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

	handler := &TurnHandler{
		Client: client,
	}

	if a.Toolset != nil {
		toolset, err := a.Toolset.Toolset(ctx)
		if err != nil {
			return nil, err
		}
		handler.Toolset = toolset
	}

	req := &api.ChatRequest{
		Model: a.Model,
	}

	if a.System != nil {
		messages, err := a.System.Messages(ctx)
		if err != nil {
			return nil, err
		}
		req.Messages = append(req.Messages, messages...)
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
		client: handler,
		ctx:    ctx,
		req:    req,
	}

	return res, nil
}
