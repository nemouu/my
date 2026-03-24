package shell

import (
	"os"
	"os/signal"
	"syscall"
)

// InitSignalHandlers sets up signal handling for the shell.
// Call this once from repl.go when the shell starts.
func InitSignalHandlers() {
	sigs := make(chan os.Signal, 1)

	// Register the signals we want to handle
	// Add more signals here as the shell grows
	signal.Notify(sigs,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTSTP, // Ctrl+Z
		syscall.SIGCHLD, // Child process state changed
	)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGINT:
				handleSIGINT()
			case syscall.SIGTSTP:
				handleSIGTSTP()
			case syscall.SIGCHLD:
				handleSIGCHLD()
			}
		}
	}()
}

// handleSIGINT handles Ctrl+C (SIGINT).
// Should forward signal to the current foreground job if one exists.
// TODO: Milestone 1 - forward to foreground process group
func handleSIGINT() {
}

// handleSIGTSTP handles Ctrl+Z (SIGTSTP).
// Should stop the current foreground job if one exists.
// TODO: Milestone 1 - forward to foreground process group
func handleSIGTSTP() {
}

// handleSIGCHLD handles child process state changes (SIGCHLD).
// Responsible for reaping zombie processes.
// TODO: Milestone 1 - call waitpid to reap finished children
func handleSIGCHLD() {
}