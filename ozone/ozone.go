package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Chatter interface {
	Chat(ctx context.Context, req *api.ChatRequest, fn api.ChatResponseFunc) error
}

type MessagesProvider interface {
	Messages(ctx context.Context) ([]api.Message, error)
}

type ToolsProvider interface {
	Tools(ctx context.Context) (api.Tools, error)
}

type ToolCallHandler interface {
	Call(ctx context.Context, tc api.ToolCall) (string, error)
}

type Toolset interface {
	ToolsProvider
	ToolCallHandler
}
