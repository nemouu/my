package shell

import (
	"os"
	"os/exec"
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
func execute(args []string) error {
	if len(args) == 0 {
		return nil
	}

	if isBuiltin(args[0]) {
		return executeBuiltin(args)
	}

	return executeExternal(args)
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
