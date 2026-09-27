package ozone

import (
	"context"
	"errors"
	"io"
	"iter"
	"sync/atomic"

	"github.com/ollama/ollama/api"
)

var errStop = errors.New("stop")

var (
	ErrStarted     = errors.New("stream already started")
	ErrNotFinished = errors.New("run not finished")
)

type RunResult struct {
	client Chatter
	closer io.Closer
	ctx    context.Context
	req    *api.ChatRequest
	err    error

	started  atomic.Bool
	finished chan struct{}

	chatStart int
	turnStart int
}

func (r *RunResult) Stream() iter.Seq2[api.ChatResponse, error] {
	return func(yield func(api.ChatResponse, error) bool) {
		if !r.started.CompareAndSwap(false, true) {
			yield(api.ChatResponse{}, ErrStarted)
			return
		}
		defer close(r.finished)

		err := r.client.Chat(r.ctx, r.req, func(cr api.ChatResponse) error {
			if !yield(cr, nil) {
				return errStop
			}
			return nil
		})
		if err != nil && !errors.Is(err, errStop) {
			r.err = err
			yield(api.ChatResponse{}, err)
		}
	}
}

func (r *RunResult) NewMessages() ([]api.Message, error) {
	select {
	default:
		return nil, ErrNotFinished
	case <-r.finished:
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.req.Messages[r.turnStart:], nil
}

func (r *RunResult) Messages(ctx context.Context) ([]api.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.finished:
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.req.Messages[r.chatStart:], nil
}

func (r *RunResult) Close() error {
	if r.closer == nil {
		return nil
	}
	return r.closer.Close()
}
