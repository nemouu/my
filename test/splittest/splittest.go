package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// splittest forks a child process that sleeps for a specified number of seconds.
// The parent process waits for the child to finish. This is used to test
// process groups, job control, and signal handling in the shell.
//
// Usage:
//
//	splittest <seconds>
//
// Example:
//
//	splittest 4   # Forks a child that sleeps for 4 seconds.
//	splittest 2 & # Runs in the background.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: splittest <seconds>")
		os.Exit(1)
	}
	secs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}

	// Launch a subprocess (e.g., "sleep" or your myspin)
	cmd := exec.Command("sleep", strconv.Itoa(secs))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	fmt.Printf("Child process (PID: %d) started\n", cmd.Process.Pid)
	cmd.Wait()
	fmt.Println("Child process finished")
}
