package shell

import (
	"fmt"
	"os"
	"strconv"
)

// This file will contain implementations of built-in shell commands
// that cannot be external programs (cd, exit, pwd, etc.) as well as
// AI-powered commands (ask, commit, code).
//
// TODO: Implement built-ins as described in TODO.md Milestone 1 "Built-in Commands"
// - cd: Change directory with os.Chdir()
// - exit: Clean exit with os.Exit(0)
// - pwd: Print working directory with os.Getwd()
//
// TODO: Later add AI commands from Milestones 3-5
// - ask: Query local LLM
// - commit: AI-powered git commit messages
// - code: Interactive code mode

// IsBuiltin checks if a command name is a built-in command
func isBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "exit", "pwd":
		return true
	}
	return false
}

// ExecuteBuiltin executes a built-in command
func executeBuiltin(args []string) error {
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			// go to home directory
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return os.Chdir(home)
		}
		return os.Chdir(args[1])
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(dir)
		return nil
	case "exit":
		if len(args) > 1 {
			status, err := strconv.Atoi((args[1]))
			if err != nil {
				return err
			}
			os.Exit(status)
		} else {
			os.Exit(0)
		}
	default:
		return nil
	}
	return nil
}
