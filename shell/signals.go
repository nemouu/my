package shell

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
)

// InitSignalHandlers sets up signal handling for the shell.
// Call this once from repl.go when the shell starts.
func initSignalHandlers() {
	sigs := make(chan os.Signal, 10)

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
func handleSIGINT() {
	job, err := getForegroundJob()
	if err != nil {
		return // no foreground job, nothing to do
	}
	syscall.Kill(-job.pid, syscall.SIGINT)
}

// handleSIGTSTP handles Ctrl+Z (SIGTSTP).
// Should stop the current foreground job if one exists.
func handleSIGTSTP() {
	job, err := getForegroundJob()
	if err != nil {
		return // no foreground job, nothing to do
	}
	syscall.Kill(-job.pid, syscall.SIGTSTP)
	job.state = ST
}

// handleSIGCHLD handles child process state changes (SIGCHLD).
// Responsible for reaping zombie processes.
func handleSIGCHLD() {
	for {
		var status syscall.WaitStatus
		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
		if err != nil || pid <= 0 {
			break
		}
		if status.Stopped() {
			job, err := getJobByPid(pid)
			if err == nil {
				job.state = ST
			}
		} else {
			deleteJob(pid)
		}
	}
}

// Helper function
func getForegroundJob() (*Job, error) {
	for i := range jobs {
		if jobs[i].state == FG {
			return &jobs[i], nil
		}
	}
	return nil, errors.New("no foreground job found")
}
