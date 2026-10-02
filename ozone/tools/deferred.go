package tools

import (
	"context"
	"fmt"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type Deferred struct {
	Name        string
	Description string
	Provider    ozone.ToolsetProvider
}

func (d *Deferred) Toolset(ctx context.Context) (ozone.Toolset, error) {
	t, err := d.Provider.Toolset(ctx)
	if err != nil {
		return nil, err
	}
	return &deferredToolset{
		name:        d.Name,
		description: d.Description,
		toolset:     t,
	}, nil
}

type deferredToolset struct {
	name        string
	description string
	enabled     bool
	toolset     ozone.Toolset
}

func (dt *deferredToolset) Tools(ctx context.Context) (api.Tools, error) {
	if dt.enabled {
		return dt.toolset.Tools(ctx)
	}
	return api.Tools{{
		Type: "function",
		Function: api.ToolFunction{
			Name:        dt.name,
			Description: dt.description,
		},
	}}, nil
}

func (dt *deferredToolset) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	if dt.enabled {
		return dt.toolset.Call(ctx, tc)
	}
	if tc.Function.Name != dt.name {
		return "", ozone.ErrNotHandled
	}
	if tc.Function.Arguments.Len() > 0 {
		return fmt.Sprintf("Error calling %s: no arguments expected", dt.name), nil
	}
	dt.enabled = true
	return "toolset activated", nil
}
