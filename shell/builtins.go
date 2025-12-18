package shell

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
func IsBuiltin(cmd string) bool {
	// TODO: Check against list of built-in commands
	return false
}

// ExecuteBuiltin executes a built-in command
func ExecuteBuiltin(args []string) error {
	// TODO: Implement built-in command routing and execution
	return nil
}
