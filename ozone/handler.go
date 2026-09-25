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
		var msg api.Message
		err := h.Client.Chat(ctx, req, func(cr api.ChatResponse) error {
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
			res, err := h.Call(ctx, tc)
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

	return ErrMaxStepsExceeded
}

func (h *TurnHandler) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	res, err := h.Call(ctx, tc)
	if errors.Is(err, ErrNotHandled) {
		return fmt.Sprintf("unkown tool: %s", tc.Function.Name), nil
	}
	return res, err
}
