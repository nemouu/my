package shell

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
	// TODO: Implement command execution logic
	return nil
}
