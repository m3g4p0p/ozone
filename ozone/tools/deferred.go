package tools

import (
	"context"
	"fmt"
	"sync"

	"m3g4p0p/ozone/internal/ctxval"
	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type EnabledFunc func(ctx context.Context, name string) bool

type Deferred struct {
	Name        string
	Description string
	Provider    ozone.ToolsetProvider
	Enabled     EnabledFunc
}

func (d *Deferred) Toolset(ctx context.Context) (ozone.Toolset, error) {
	t, err := d.Provider.Toolset(ctx)
	if err != nil {
		return nil, err
	}
	return &deferredToolset{
		name:        d.Name,
		description: d.Description,
		enabledFunc: d.Enabled,
		toolset:     t,
	}, nil
}

type deferredToolset struct {
	name        string
	description string
	toolset     ozone.Toolset
	enabledFunc EnabledFunc
	once        sync.Once
	enabled     bool
}

func (dt *deferredToolset) Tools(ctx context.Context) (api.Tools, error) {
	dt.once.Do(func() {
		dt.initialize(ctx)
	})
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

func (dt *deferredToolset) initialize(ctx context.Context) {
	if dt.enabledFunc != nil {
		dt.enabled = dt.enabledFunc(ctx, dt.name)
	}
}

func HasActivateMessage(ctx context.Context, name string) bool {
	req, ok := ctxval.From[*api.ChatRequest](ctx)
	if !ok {
		return false
	}

	for _, m := range req.Messages {
		if m.Role == "tool" && m.ToolName == name {
			return true
		}
	}

	return false
}
