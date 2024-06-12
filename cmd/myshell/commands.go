package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ExitCommand struct {
	name string
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{name: "exit"}
}

func (e *ExitCommand) Name() string {
	return e.name
}

func (e *ExitCommand) Exec(args []string) error {
	if len(args) == 0 {
		os.Exit(0)
	}
	code, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	os.Exit(code)

	return nil
}

type EchoCommand struct {
	name string
}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{name: "echo"}
}

func (e *EchoCommand) Name() string {
	return "echo"
}

func (e *EchoCommand) Exec(args []string) error {
	msg := args
	newLine := true
	if len(args) > 0 && args[0] == "-n" {
		msg = msg[1:]
		newLine = false
	}

	smsg := strings.Join(msg, " ")
	if newLine {
		_, err := fmt.Fprintln(os.Stderr, smsg)
		return err
	}

	_, err := fmt.Fprint(os.Stderr, smsg)
	return err
}
