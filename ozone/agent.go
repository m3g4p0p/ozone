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

	messages, err := a.resolveSystem(ctx)
	if err != nil {
		return nil, err
	}
	messages = append(messages, options.History...)

	toolset, err := a.resolveToolset(ctx)
	if err != nil {
		return nil, err
	}

	handler := &TurnHandler{
		Client:  client,
		Toolset: toolset,
	}

	req := &api.ChatRequest{
		Model:    a.Model,
		Messages: messages,
	}

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

func (a *Agent) resolveSystem(ctx context.Context) ([]api.Message, error) {
	if a.System == nil {
		return nil, nil
	}
	return a.System.Messages(ctx)
}

func (a *Agent) resolveToolset(ctx context.Context) (Toolset, error) {
	if a.Toolset == nil {
		return nil, nil
	}
	return a.Toolset.Toolset(ctx)
}
