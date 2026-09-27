package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"m3g4p0p/ozone/ozone"
	"m3g4p0p/ozone/ozone/messages"
	"m3g4p0p/ozone/ozone/think"
)

func runAgent(agent *ozone.Agent, input ozone.MessagesProvider) (ozone.MessagesProvider, error) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	res, err := agent.Run(ctx, input, nil)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for cr, err := range res.Stream() {
		if err != nil {
			return nil, err
		}

		fmt.Print(cr.Message.Content)
		if cr.Done {
			fmt.Println()
		}
	}

	return res, err
}

func readline(s *bufio.Scanner) bool {
	fmt.Print("> ")
	return s.Scan()
}

func run() error {
	agent := &ozone.Agent{
		Name:  "main",
		Model: "qwen3.5:2b",
		Think: think.False,
	}

	conversation := messages.Messages()
	scanner := bufio.NewScanner(os.Stdin)

	for readline(scanner) {
		input := messages.Messages(
			conversation,
			messages.User(scanner.Text()),
		)

		res, err := runAgent(agent, input)
		if err != nil {
			return err
		}
		conversation = res
	}

	return scanner.Err()
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
