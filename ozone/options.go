package ozone

import "github.com/ollama/ollama/api"

type RunOptions struct {
	Client  Chatter
	History []api.Message
}

func (r *RunOptions) resolveClient() (Chatter, error) {
	if r.Client != nil {
		return r.Client, nil
	}
	return api.ClientFromEnvironment()
}
