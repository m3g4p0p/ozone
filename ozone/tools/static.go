package tools

import (
	"context"

	"m3g4p0p/ozone/ozone"
)

type staticToolset struct {
	toolset ozone.Toolset
}

func (t *staticToolset) Toolset(context.Context) (ozone.Toolset, error) {
	return t.toolset, nil
}

func Static(toolset ozone.Toolset) ozone.ToolsetProvider {
	return &staticToolset{toolset}
}
