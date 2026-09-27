package messages

import (
	"context"

	"m3g4p0p/ozone/ozone"

	"github.com/ollama/ollama/api"
)

type User string

func (u User) Messages(context.Context) ([]api.Message, error) {
	return createMessage("user", string(u))
}

type System string

func (s System) Messages(context.Context) ([]api.Message, error) {
	return createMessage("system", string(s))
}

func createMessage(role, content string) ([]api.Message, error) {
	return []api.Message{{Role: role, Content: content}}, nil
}

func Static(m ...api.Message) ozone.MessagesProvider {
	return staticMessages(m)
}

type staticMessages []api.Message

func (s staticMessages) Messages(context.Context) ([]api.Message, error) {
	return s, nil
}

func Messages(mp ...ozone.MessagesProvider) ozone.MessagesProvider {
	return messageProviders(mp)
}

type messageProviders []ozone.MessagesProvider

func (m messageProviders) Messages(ctx context.Context) ([]api.Message, error) {
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
