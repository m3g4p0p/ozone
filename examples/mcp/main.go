package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"m3g4p0p/ozone/ozone"
	ozonemcp "m3g4p0p/ozone/ozone/mcp"
	"m3g4p0p/ozone/ozone/messages"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	t1, t2 := mcp.NewInMemoryTransports()
	srv := mcp.NewServer(&mcp.Implementation{}, nil)

	srv.AddPrompt(&mcp.Prompt{
		Name: "weather_expert",
	}, func(
		ctx context.Context,
		gpr *mcp.GetPromptRequest,
	) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{
				Role: mcp.Role("system"),
				Content: &mcp.TextContent{
					Text: "You are a friendly weather expert.",
				},
			}},
		}, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_weather",
		Description: "get the weather for the provided location",
	}, func(
		ctx context.Context,
		rec *mcp.CallToolRequest,
		params struct {
			Location string `json:"location"`
		},
	) (*mcp.CallToolResult, string, error) {
		return nil, "mostly sunny", nil
	})

	if session, err := srv.Connect(ctx, t1, nil); err != nil {
		return err
	} else {
		defer session.Close()
	}

	mcpProvider := &ozonemcp.Provider{
		Client:    mcp.NewClient(&mcp.Implementation{}, nil),
		Transport: t2,
	}

	agent := &ozone.Agent{
		Name:    "main",
		Model:   "qwen3.5:2b",
		Toolset: mcpProvider,
		System:  mcpProvider.Prompt("weather_expert"),
	}

	res, err := agent.Run(ctx, messages.User(flag.Arg(0)), nil)
	if err != nil {
		return err
	}
	defer res.Close()

	for cr, err := range res.Stream() {
		if err != nil {
			return err
		}
		fmt.Print(cr.Message.Thinking + cr.Message.Content)
		if cr.Done {
			fmt.Println()
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
