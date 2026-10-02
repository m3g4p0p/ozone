package mcp

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"strings"

	"m3g4p0p/ozone/ozone"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type Toolset struct {
	Session *mcp.ClientSession
	closer  io.Closer
}

func (ts *Toolset) Tools(ctx context.Context) (api.Tools, error) {
	var tools api.Tools
	for t, err := range ts.Session.Tools(ctx, nil) {
		if err != nil {
			return nil, err
		}
		var params api.ToolFunctionParameters
		if err := convert(t.InputSchema, &params); err != nil {
			return nil, err
		}
		tools = append(tools, api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}
	return tools, nil
}

func (ts *Toolset) Call(ctx context.Context, tc api.ToolCall) (string, error) {
	if err := ts.validateTool(ctx, tc.Function.Name); err != nil {
		return "", err
	}

	res, err := ts.Session.CallTool(ctx, &mcp.CallToolParams{
		Name:      tc.Function.Name,
		Arguments: tc.Function.Arguments,
	})
	if err != nil {
		return "", err
	}

	if res.StructuredContent != nil {
		data, err := json.Marshal(res.StructuredContent)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	var s strings.Builder
	if res.IsError {
		fmt.Fprintf(&s, "Error calling %s: ", tc.Function.Name)
	}
	for i, c := range res.Content {
		if i > 0 {
			fmt.Fprintln(&s)
		}
		if tc, ok := c.(*mcp.TextContent); ok {
			s.WriteString(tc.Text)
		}
		if er, ok := c.(*mcp.EmbeddedResource); ok {
			s.WriteString(er.Resource.Text)
		}
	}
	return s.String(), nil
}

func (ts *Toolset) Toolset(context.Context) (ozone.Toolset, error) {
	return ts, nil
}

func (ts *Toolset) Close() error {
	if ts.closer == nil {
		return nil
	}
	return ts.closer.Close()
}

func (ts *Toolset) validateTool(ctx context.Context, name string) error {
	for t, err := range ts.Session.Tools(ctx, nil) {
		if err != nil {
			return err
		}
		if t.Name == name {
			return nil
		}
	}
	return ozone.ErrNotHandled
}

func convert(src, dest any) error {
	tmp, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, dest)
}
