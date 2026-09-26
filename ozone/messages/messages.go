package messages

import (
	"context"

	"github.com/ollama/ollama/api"
)

type User string

func (u User) Messages(context.Context) ([]api.Message, error) {
	return static("user", string(u))
}

type System string

func (s System) Messages(context.Context) ([]api.Message, error) {
	return static("system", string(s))
}

type Messages []api.Message

func (m Messages) Messages(context.Context) ([]api.Message, error) {
	return m, nil
}

func static(role, content string) ([]api.Message, error) {
	return []api.Message{{Role: role, Content: content}}, nil
}
