package shell

import (
	"bufio"
	"fmt"
	"os"
)

// This file will contain the main REPL (Read-Eval-Print Loop) implementation
// for the μ shell. It handles reading user input, executing commands, and
// printing the results in a loop.
//
// TODO: Implement the REPL as described in TODO.md Milestone 1 "Core REPL Implementation"
// - printPrompt() to display the shell prompt
// - readLine() using bufio.Scanner
// - Main loop that continuously reads and executes commands

// Run starts the main shell REPL loop
func Run() error {
	initSignalHandlers()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt())
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		cmd, err := parse(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := execute(cmd.Args); err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

// Print the prompt --- change this later!
func prompt() string {
	dir, err := os.Getwd()
	if err != nil {
		return "μ> "
	}
	return fmt.Sprintf("μ %s> ", dir)
}
