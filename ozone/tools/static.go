package tools

import (
	"context"

	"m3g4p0p/ozone/ozone"
)

func Static(toolset ozone.Toolset) ozone.ToolsetProvider {
	return ToolsetProviderFunc(func(ctx context.Context) (ozone.Toolset, error) {
		return toolset, nil
	})
}
