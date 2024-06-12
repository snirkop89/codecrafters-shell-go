package main

import (
	"os"
	"strconv"
)

type ExitCommand struct {
	name string
	args []string
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
