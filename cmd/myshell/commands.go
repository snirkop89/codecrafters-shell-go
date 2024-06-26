package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const builtin = "builtin"

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
	return builtin
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
	return builtin
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
	return builtin
}

func (tc *TypeCommand) Exec(_ []string) error {
	return nil
}

// Worked out of the box since I have the lookup in path,
// but for fun, I'll implement it myself.
type PwdCommand struct{}

func (pc *PwdCommand) Name() string {
	return "pwd"
}

func (pc *PwdCommand) Exec(_ []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, dir)
	return nil
}

func (pc *PwdCommand) Type() string {
	return "builtin"
}

type CdCommand struct{}

func (cd *CdCommand) Name() string {
	return "cd"
}

func (cd *CdCommand) Type() string {
	return builtin
}

func (cd *CdCommand) Exec(args []string) error {
	switch {
	case len(args) == 0:
		homeFolder := os.Getenv("HOME")
		if homeFolder == "" {
			return fmt.Errorf("could not find HOME dir")
		}
		err := os.Chdir(homeFolder)
		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", args[0])
		}

	case len(args) == 1 && args[0] == "~":
		homeFolder := os.Getenv("HOME")
		if homeFolder == "" {
			return fmt.Errorf("could not find HOME dir")
		}
		err := os.Chdir(homeFolder)
		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", args[0])
		}
	default:
		dir, err := handleRelative(args[0])
		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", args[0])
		}
		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", args[0])
		}
	}
	return nil
}

func handleRelative(input string) (string, error) {
	// Absolute path
	if input[0] == '/' {
		return input, nil
	}

	// Relative to current dir
	if input[0] == '.' && len(input) == 1 {
		return os.Getwd()
	} else if len(input) > 1 && input[0] == '.' && input[1] != '.' {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return path.Join(cwd, input[1:]), nil
	}

	if len(input) > 1 {
		idx := strings.Index(input, "/")
		isPrev := input[0] == '.' && input[1] == '.'

		if isPrev && idx == -1 {
			tdir, err := os.Getwd()
			if err != nil {
				return "", err
			}
			prev := path.Dir(tdir)
			return prev, nil
		} else if isPrev && idx > 2 {
			tdir, err := os.Getwd()
			if err != nil {
				return "", err
			}
			prev := path.Dir(tdir)
			return path.Join(prev, input[2:]), nil
		}

	}

	return input, nil
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
	cmd := exec.Command(ex.binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stderr
	return cmd.Run()
}
