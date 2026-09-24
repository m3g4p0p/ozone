package tools

import (
	"context"
	"errors"
	"fmt"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

var ErrNotHandled = errors.New("not handled")

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
		if errors.Is(err, ErrNotHandled) {
			continue
		}
		return res, err
	}
	return "", fmt.Errorf("%w: %s", ErrNotHandled, tc.Function.Name)
}
