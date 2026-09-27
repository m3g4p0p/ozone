package messages

import (
	"context"

	"m3g4p0p/ozone/ozone"

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

func static(role, content string) ([]api.Message, error) {
	return []api.Message{{Role: role, Content: content}}, nil
}

type Messages []ozone.MessagesProvider

func (m Messages) Messages(ctx context.Context) ([]api.Message, error) {
	var messages []api.Message
	for _, provider := range m {
		m, err := provider.Messages(ctx)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m...)
	}
	return messages, nil
}
