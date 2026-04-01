package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
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
func Run(emitPrompt bool) error {
	initSignalHandlers()

	os.MkdirAll(os.Getenv("HOME")+"/.config/mu", 0755)

	p := prompt()
	if !emitPrompt {
		p = ""
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       p,
		HistoryFile:  os.Getenv("HOME") + "/.config/mu/history",
		HistoryLimit: 500,
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		if emitPrompt {
			rl.SetPrompt(prompt())
		}

		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		cmd, err := parse(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := execute(cmd.Args, cmd.Bg); err != nil {
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
