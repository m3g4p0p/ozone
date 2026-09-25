package ozone

import (
	"context"
	"errors"
	"iter"
	"sync/atomic"

	"github.com/ollama/ollama/api"
)

var errStop = errors.New("stop")

var ErrStarted = errors.New("stream already started")

type RunResult struct {
	client  Chatter
	ctx     context.Context
	req     *api.ChatRequest
	started atomic.Bool
}

func (r *RunResult) Stream() iter.Seq2[api.ChatResponse, error] {
	return func(yield func(api.ChatResponse, error) bool) {
		if !r.started.CompareAndSwap(false, true) {
			yield(api.ChatResponse{}, ErrStarted)
			return
		}

		err := r.client.Chat(r.ctx, r.req, func(cr api.ChatResponse) error {
			if !yield(cr, nil) {
				return errStop
			}
			return nil
		})
		if err != nil && !errors.Is(err, errStop) {
			yield(api.ChatResponse{}, err)
		}
	}
}
