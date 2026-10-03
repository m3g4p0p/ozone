package tools

import (
	"context"

	"m3g4p0p/ozone/ozone"
)

type ToolsetProviderFunc func(ctx context.Context) (ozone.Toolset, error)

func (f ToolsetProviderFunc) Toolset(ctx context.Context) (ozone.Toolset, error) {
	return f(ctx)
}
