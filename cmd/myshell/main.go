package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	Type() string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	input := readInput()
	commands := []cmd{
		NewExitCommand(),
		NewEchoCommand(),
		&PwdCommand{},
		&CdCommand{},
		&TypeCommand{},
	}

outer:
	for {
		fmt.Fprint(os.Stdout, "$ ")
		select {
		case <-ctx.Done():
			break outer
		case userInput := <-input:
			userCmd, args, err := extractCommand(userInput)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			runCommand(commands, userCmd, args)
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

// extractCommand returns the command and its args
func extractCommand(input string) (string, []string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil, ErrInvalidInput
	}
	if len(parts) == 1 {
		return parts[0], nil, nil
	}
	return parts[0], parts[1:], nil
}

func runCommand(commands []cmd, command string, args []string) {
	var execCmd cmd
	var reflect bool

	if command == "type" {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "expected 1 arg, got none")
			return
		}
		reflect = true
		command = args[0]
	}

	// Check for builtin commands
	for _, cmd := range commands {
		if cmd.Name() == command {
			execCmd = cmd
			break
		}
	}
	// Fallback for external command
	if execCmd == nil {
		cmd, err := exec.LookPath(command)
		if err == nil {
			execCmd = &ExternalCommand{name: command, binPath: cmd}
		}
	}

	if execCmd == nil {
		fmt.Fprintf(os.Stdout, "%s: not found\n", command)
		return
	}

	// type is a unique command of reflect a command
	if reflect {
		if execCmd.Type() == "builtin" {
			fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", args[0])
		} else {
			fmt.Fprintf(os.Stdout, "%s is %s\n", args[0], execCmd.Type())
		}
		return
	}

	if err := execCmd.Exec(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
