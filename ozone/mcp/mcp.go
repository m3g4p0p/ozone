package mcp

import (
	"context"

	"m3g4p0p/ozone/ozone"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Provider struct {
	Client    *mcp.Client
	Transport mcp.Transport
}

func (m *Provider) Toolset(ctx context.Context) (ozone.Toolset, error) {
	session, err := m.Client.Connect(ctx, m.Transport, nil)
	if err != nil {
		return nil, err
	}
	return &Toolset{
		Session: session,
		cleanup: true,
	}, nil
}
