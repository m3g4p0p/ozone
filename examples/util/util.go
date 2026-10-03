package util

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"m3g4p0p/ozone/ozone"
	"m3g4p0p/ozone/ozone/messages"
)

func RunAgent(ctx context.Context, agent *ozone.Agent, input ozone.MessagesProvider) (ozone.MessagesProvider, error) {
	ctx, stop := signal.NotifyContext(
		ctx,
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

		fmt.Print(cr.Message.Thinking + cr.Message.Content)
		if cr.Done {
			fmt.Println()
		}
	}

	return res, err
}

func RunSingle(ctx context.Context, agent *ozone.Agent, input string) error {
	_, err := RunAgent(ctx, agent, messages.User(input))
	return err
}

func RunLoop(ctx context.Context, agent *ozone.Agent) error {
	conversation := messages.Messages()
	scanner := bufio.NewScanner(os.Stdin)

	for Readline(scanner) {
		input := messages.Messages(
			conversation,
			messages.User(scanner.Text()),
		)

		res, err := RunAgent(ctx, agent, input)
		if err != nil {
			return err
		}
		conversation = res
	}

	return scanner.Err()
}

func Readline(s *bufio.Scanner) bool {
	fmt.Print("> ")
	return s.Scan()
}
