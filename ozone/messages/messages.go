package messages

import (
	"context"

	"github.com/ollama/ollama/api"
)

type System string

func (s System) Messages(context.Context) ([]api.Message, error) {
	return []api.Message{{
		Role:    "system",
		Content: string(s),
	}}, nil
}

type Messages []api.Message

func (m Messages) Messages(context.Context) ([]api.Message, error) {
	return m, nil
}
