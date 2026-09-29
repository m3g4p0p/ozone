package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type Prompt struct {
	Session *mcp.ClientSession
	Name    string
	cleanup bool
}

func (p *Prompt) Messages(ctx context.Context) ([]api.Message, error) {
	res, err := p.Session.GetPrompt(ctx, &mcp.GetPromptParams{
		Name: p.Name,
	})
	if err != nil {
		return nil, err
	}
	var messages []api.Message
	for _, m := range res.Messages {
		switch c := m.Content.(type) {
		case *mcp.TextContent:
			messages = append(messages, api.Message{
				Role:    string(m.Role),
				Content: c.Text,
			})
		}
	}
	return messages, nil
}

func (p *Prompt) Close() error {
	if !p.cleanup {
		return nil
	}
	return p.Session.Close()
}
