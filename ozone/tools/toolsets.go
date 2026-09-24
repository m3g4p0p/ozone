package tools

import (
	"context"
	"errors"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type Toolsets []ozone.Toolset

func (t Toolsets) Tools(ctx context.Context) (api.Tools, error) {
	var tools api.Tools
	for _, t := range t {
		ts, err := t.Tools(ctx)
		if err != nil {
			return nil, err
		}
		tools = append(tools, ts...)
	}
	return tools, nil
}

func (t Toolsets) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	for _, t := range t {
		res, err := t.Call(ctx, tc)
		if errors.Is(err, ozone.ErrNotHandled) {
			continue
		}
		return res, err
	}
	return "", ozone.ErrNotHandled
}
