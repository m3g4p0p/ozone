package ozone

import (
	"context"

	"github.com/ollama/ollama/api"
)

type ToolHandler struct {
	Toolset
	Client Chatter
}

func (m *ToolHandler) Chat(
	ctx context.Context,
	req *api.ChatRequest,
	fn api.ChatResponseFunc,
) error {
	for {
		var msg api.Message
		err := m.Client.Chat(ctx, req, func(cr api.ChatResponse) error {
			if err := fn(cr); err != nil {
				return err
			}
			msg.Thinking += cr.Message.Thinking
			msg.Content += cr.Message.Content
			msg.ToolCalls = append(
				msg.ToolCalls,
				cr.Message.ToolCalls...,
			)
			return nil
		})
		if err != nil {
			return err
		}

		req.Messages = append(req.Messages, msg)
		if len(msg.ToolCalls) == 0 {
			return nil
		}

		for _, tc := range msg.ToolCalls {
			res, err := m.Call(ctx, tc)
			if err != nil {
				return err
			}
			req.Messages = append(req.Messages, api.Message{
				Role:       "tool",
				Content:    res,
				ToolCallID: tc.ID,
			})
		}
	}
}
