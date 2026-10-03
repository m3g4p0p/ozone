package tools

import (
	"context"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type Lazy struct {
	Provider ozone.ToolsetProvider
}

func (l *Lazy) Toolset(context.Context) (ozone.Toolset, error) {
	return &lazyToolset{provider: l.Provider}, nil
}

type lazyToolset struct {
	provider ozone.ToolsetProvider
	ozone.Toolset
}

func (lt *lazyToolset) Tools(ctx context.Context) (api.Tools, error) {
	if lt.Toolset == nil {
		t, err := lt.provider.Toolset(ctx)
		if err != nil {
			return nil, err
		}
		lt.Toolset = t
	}
	return lt.Toolset.Tools(ctx)
}
