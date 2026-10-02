package ozone

import (
	"context"
	"errors"
	"fmt"

	"github.com/ollama/ollama/api"
)

const defaultMaxSteps = 10

var (
	ErrNotHandled       = errors.New("not handled")
	ErrMaxStepsExceeded = errors.New("max steps exceeded")
)

type TurnHandler struct {
	Client   Chatter
	Toolset  Toolset
	MaxSetps int
}

func (h *TurnHandler) Chat(
	ctx context.Context,
	req *api.ChatRequest,
	fn api.ChatResponseFunc,
) error {
	maxSteps := h.MaxSetps
	if maxSteps == 0 {
		maxSteps = defaultMaxSteps
	}

	for range maxSteps {
		tools, err := h.Tools(ctx)
		if err != nil {
			return err
		}
		req.Tools = tools

		var msg api.Message
		if err := h.Client.Chat(ctx, req, func(cr api.ChatResponse) error {
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
		}); err != nil {
			return err
		}

		req.Messages = append(req.Messages, msg)
		if len(msg.ToolCalls) == 0 {
			return nil
		}

		for _, tc := range msg.ToolCalls {
			res, err := h.Call(ctx, tc)
			if err != nil {
				return err
			}
			req.Messages = append(req.Messages, api.Message{
				Role:       "tool",
				Content:    res,
				ToolCallID: tc.ID,
				ToolName:   tc.Function.Name,
			})
		}
	}

	return ErrMaxStepsExceeded
}

func (h *TurnHandler) Tools(ctx context.Context) (api.Tools, error) {
	if h.Toolset == nil {
		return nil, nil
	}
	return h.Toolset.Tools(ctx)
}

func (h *TurnHandler) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	if h.Toolset == nil {
		return unkownTool(tc)
	}
	res, err := h.Toolset.Call(ctx, tc)
	if errors.Is(err, ErrNotHandled) {
		return unkownTool(tc)
	}
	return res, err
}

func unkownTool(tc api.ToolCall) (string, error) {
	return fmt.Sprintf("unkown tool: %s", tc.Function.Name), nil
}
