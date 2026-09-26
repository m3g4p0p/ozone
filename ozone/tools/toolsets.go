package tools

import (
	"context"
	"errors"
	"io"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type Toolsets []ozone.ToolsetProvider

func (t Toolsets) Toolset(ctx context.Context) (ozone.Toolset, error) {
	var toolsets resolvedToolsets
	for _, provider := range t {
		t, err := provider.Toolset(ctx)
		if err != nil {
			_ = toolsets.Close()
			return nil, err
		}
		toolsets = append(toolsets, t)
	}
	return toolsets, nil
}

type resolvedToolsets []ozone.Toolset

func (t resolvedToolsets) Tools(ctx context.Context) (api.Tools, error) {
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

func (t resolvedToolsets) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	for _, t := range t {
		res, err := t.Call(ctx, tc)
		if errors.Is(err, ozone.ErrNotHandled) {
			continue
		}
		return res, err
	}
	return "", ozone.ErrNotHandled
}

func (t resolvedToolsets) Close() error {
	var errs []error
	for _, t := range t {
		if c, ok := t.(io.Closer); ok {
			errs = append(errs, c.Close())
		}
	}
	return errors.Join(errs...)
}
