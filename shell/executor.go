package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// This file will contain command execution logic for the shell.
// It determines whether a command is built-in or external and
// executes it appropriately using fork+exec for external commands.
//
// TODO: Implement execution as described in TODO.md Milestone 1 "Command Execution"
// - Check if command is built-in
// - Execute built-ins directly
// - Execute external commands using exec.Command()
// - Handle command not found errors

// Execute runs a parsed command with its arguments
func Execute(args []string) error {
	if len(args) == 0 {
		return nil
	}

	if isBuiltin(args[0]) {
		return executeBuiltin(args)
	}

	return executeExternal(args)
}

func isBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "exit", "pwd":
		return true
	}
	return false
}

func executeExternal(args []string) error {
	// Create Command
	cmd := exec.Command(args[0], args[1:]...)

	// Wire up stdin/stdout/stderr so the child process can talk to the terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run command and return errors if they happen
	return cmd.Run()
}

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
