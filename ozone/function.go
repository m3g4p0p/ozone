package ozone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ollama/ollama/api"
)

type llmError struct {
	message string
}

func (e *llmError) Error() string {
	return e.message
}

func LLMError(message string) error {
	return &llmError{message: message}
}

type FunctionTool struct {
	Definition api.Tool
	Handler    func(ctx context.Context, tc api.ToolCall) (string, error)
}

func (t *FunctionTool) Tools(context.Context) (api.Tools, error) {
	return api.Tools{t.Definition}, nil
}

func (t *FunctionTool) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	return t.Handler(ctx, tc)
}

func NewTool[In, Out any](
	name, description string,
	handler func(ctx context.Context, input In) (Out, error),
) (Toolset, error) {
	schema, err := jsonschema.For[In](nil)
	if err != nil {
		return nil, err
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		return nil, err
	}
	var params api.ToolFunctionParameters
	if err := convert(schema, &params); err != nil {
		return nil, err
	}
	return &FunctionTool{
		Definition: api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        name,
				Description: description,
				Parameters:  params,
			},
		},
		Handler: func(ctx context.Context, tc api.ToolCall) (string, error) {
			if tc.Function.Name != name {
				return "", ErrNotHandled
			}

			if err := resolved.Validate(tc.Function.Arguments.ToMap()); err != nil {
				// Retries simply limited by run level max steps
				return toolCallError(name, err), nil
			}

			var input In
			if err := convert(tc.Function.Arguments, &input); err != nil {
				return "", err
			}

			out, err := handler(ctx, input)
			if err != nil {
				if err, ok := errors.AsType[*llmError](err); ok {
					return toolCallError(name, err), nil
				}
				return "", err
			}

			res, err := json.Marshal(out)
			if err != nil {
				return "", err
			}
			return string(res), nil
		},
	}, nil
}

func MustNewTool[In, Out any](
	name, description string,
	handler func(ctx context.Context, input In) (Out, error),
) Toolset {
	t, err := NewTool(name, description, handler)
	if err != nil {
		panic(err)
	}
	return t
}

func convert(src, dest any) error {
	tmp, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, dest)
}

func toolCallError(name string, err error) string {
	return fmt.Sprintf("error calling %s: %v", name, err)
}
