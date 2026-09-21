package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"m3g4p0p/ozone/ozone"
)

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	agent := &ozone.Agent{
		Name:  "main",
		Model: "qwen3.5:2b",
	}

	res, err := agent.Run(ctx, "say hello", nil)
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
