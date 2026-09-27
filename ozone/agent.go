package ozone

import (
	"context"
	"slices"

	"m3g4p0p/ozone/internal/cleanup"

	"github.com/ollama/ollama/api"
)

type Agent struct {
	Name    string
	Model   string
	Think   ThinkValueProvider
	System  MessagesProvider
	Toolset ToolsetProvider
}

func (a *Agent) Run(
	ctx context.Context,
	input MessagesProvider,
	options *RunOptions,
) (*RunResult, error) {
	var cleanup cleanup.Cleanup
	defer cleanup.Close()

	if options == nil {
		options = &RunOptions{}
	}

	client, err := options.resolveClient()
	if err != nil {
		return nil, err
	}

	think, err := a.resolveThink(ctx)
	if err != nil {
		return nil, err
	}

	system, err := a.resolveSystem(ctx)
	if err != nil {
		return nil, err
	}

	messages, err := input.Messages(ctx)
	if err != nil {
		return nil, err
	}

	toolset, err := a.resolveToolset(ctx)
	if err != nil {
		return nil, err
	}
	cleanup.AddAny(toolset)

	handler := &TurnHandler{
		Client:  client,
		Toolset: toolset,
	}

	req := &api.ChatRequest{
		Model:    a.Model,
		Think:    think,
		Messages: slices.Concat(system, messages),
	}

	res := &RunResult{
		client:   handler,
		closer:   cleanup.Take(),
		ctx:      ctx,
		req:      req,
		offset:   len(system),
		finished: make(chan struct{}),
	}

	return res, nil
}

func (a *Agent) resolveThink(ctx context.Context) (*api.ThinkValue, error) {
	if a.Think == nil {
		return nil, nil
	}
	return a.Think.ThinkValue(ctx)
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
