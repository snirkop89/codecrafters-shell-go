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

func (e *ExitCommand) Type() string {
	return "builtin"
}

type EchoCommand struct {
	name string
}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{name: "echo"}
}

func (e *EchoCommand) Name() string {
	return e.name
}

func (e *EchoCommand) Type() string {
	return "builtin"
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

type TypeCommand struct{}

func (tc *TypeCommand) Name() string {
	return "type"
}

func (e *TypeCommand) Type() string {
	return "builtin"
}

func (tc *TypeCommand) Exec(_ []string) error {
	return nil
}

type ExternalCommand struct {
	name    string
	binPath string
}

func (ex *ExternalCommand) Name() string {
	return ex.name
}

func (e *ExternalCommand) Type() string {
	return e.binPath
}

func (ex *ExternalCommand) Exec(args []string) error {
	// TODO: execute the command
	return nil
}
