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
	"m3g4p0p/ozone/ozone/tools"
)

type GetWeatherParams struct {
	Location string `json:"location" jsonschema:"the location to get the weather for"`
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	tool := tools.MustNewTool(
		"get_weather",
		"get the weather for the provided location",
		func(_ context.Context, params GetWeatherParams) (string, error) {
			return "mostly sunny", nil
		},
	)

	agent := &ozone.Agent{
		Name:    "main",
		Model:   "qwen3.5:2b",
		Toolset: tools.Toolsets{tool},
	}

	res, err := agent.Run(ctx, flag.Arg(0), nil)
	if err != nil {
		return err
	}

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
