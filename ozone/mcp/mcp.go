package mcp

import (
	"context"
	"encoding/json/v2"
	"strings"

	"m3g4p0p/ozone/ozone"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type Provider struct {
	Client    *mcp.Client
	Transport mcp.Transport
}

func (m *Provider) Toolset(ctx context.Context) (ozone.Toolset, error) {
	session, err := m.Client.Connect(ctx, m.Transport, nil)
	if err != nil {
		return nil, err
	}
	return &Toolset{
		Session: session,
		cleanup: true,
	}, nil
}

type Toolset struct {
	Session *mcp.ClientSession
	cleanup bool
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
		s.WriteString("Error: ")
	}
	for _, c := range res.Content {
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
	if !ts.cleanup {
		return nil
	}
	return ts.Session.Close()
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
