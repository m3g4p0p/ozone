package mcp

import (
	"context"

	"m3g4p0p/ozone/internal/resources"
	"m3g4p0p/ozone/ozone"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type Provider struct {
	Client    *mcp.Client
	Transport mcp.Transport

	session resources.Shared[*mcp.ClientSession]
}

func (p *Provider) Prompt(name string) ozone.MessagesProvider {
	return &lazyPrompt{
		provider: p,
		name:     name,
	}
}

func (p *Provider) Toolset(ctx context.Context) (ozone.Toolset, error) {
	lease, err := p.acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &Toolset{
		Session: lease.Resource(),
		closer:  lease,
	}, nil
}

func (p *Provider) acquire(ctx context.Context) (*resources.Lease[*mcp.ClientSession], error) {
	return p.session.Acquire(func() (*mcp.ClientSession, error) {
		return p.Client.Connect(ctx, p.Transport, nil)
	})
}

type lazyPrompt struct {
	provider *Provider
	name     string
}

func (p *lazyPrompt) Messages(ctx context.Context) ([]api.Message, error) {
	lease, err := p.provider.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer lease.Close()

	prompt := &Prompt{
		Name:    p.name,
		Session: lease.Resource(),
	}
	return prompt.Messages(ctx)
}
