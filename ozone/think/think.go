package think

import (
	"context"

	"github.com/ollama/ollama/api"
)

const (
	False thinkValue = iota
	True
	Low
	Medium
	High
	Max
)

type thinkValue int

func (v thinkValue) ThinkValue(context.Context) (*api.ThinkValue, error) {
	return &api.ThinkValue{Value: v.resolve()}, nil
}

func (v thinkValue) resolve() any {
	switch v {
	case False:
		return false
	case True:
		return True
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	case Max:
		return "max"
	default:
		return nil
	}
}
