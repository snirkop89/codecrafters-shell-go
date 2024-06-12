package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

var (
	ErrCommandNotFound = errors.New("command not found")
	ErrInvalidInput    = errors.New("invalid input") // Can be used to show available list of commands
)

type cmd interface {
	Name() string
	Exec(args []string) error
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	input := readInput()
	commands := []cmd{NewExitCommand()}

outer:
	for {
		fmt.Fprint(os.Stdout, "$ ")
		select {
		case <-ctx.Done():
			break outer
		case userInput := <-input:
			userCmd, err := extractCommand(userInput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				continue
			}
			for _, cmd := range commands {
				if cmd.Name() != userCmd {
					continue
				}
				if err := cmd.Exec(strings.Fields(userInput)[1:]); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				}
				continue outer
			}
			fmt.Fprintf(os.Stdout, "%s: command not found\n", userInput)
		}
	}
	fmt.Fprintln(os.Stdout, "quiting...")
}

func readInput() <-chan string {
	scanner := bufio.NewScanner(os.Stdin)

	r := make(chan string)
	go func() {
		for scanner.Scan() {
			r <- strings.TrimSpace(scanner.Text())
		}
	}()
	return r
}

func extractCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ErrInvalidInput
	}
	return parts[0], nil
}
