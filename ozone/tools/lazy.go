package tools

import (
	"context"

	"m3g4p0p/ozone/internal/resources"
	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

func Lazy(provider ozone.ToolsetProvider) ozone.ToolsetProvider {
	return ToolsetProviderFunc(func(ctx context.Context) (ozone.Toolset, error) {
		return &lazyToolset{provider: provider}, nil
	})
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

func (lt *lazyToolset) Close() error {
	return resources.CloseAny(lt.Toolset)
}
