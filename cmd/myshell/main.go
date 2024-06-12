package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	input := readInput()

outer:
	for {
		fmt.Fprint(os.Stdout, "$ ")
		select {
		case <-ctx.Done():
			break outer
		case command := <-input:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		}
	}
	fmt.Fprintln(os.Stdout, "quiting...")
}

func readInput() <-chan string {
	scanner := bufio.NewScanner(os.Stdin)

	r := make(chan string)
	go func() {
		for scanner.Scan() {
			r <- scanner.Text()
		}
	}()
	return r
}
