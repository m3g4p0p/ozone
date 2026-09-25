package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

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
