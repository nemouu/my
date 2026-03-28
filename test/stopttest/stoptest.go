package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// stoptest sleeps for a specified number of seconds, then sends SIGTSTP to itself.
// This tests the shell's ability to handle stopped jobs (e.g., with `fg` or `bg`).
//
// Usage:
//
//	stoptest <seconds>
//
// Example:
//
//	stoptest 3   # Sleeps for 3 seconds, then stops itself.
//	stoptest 2 & # Runs in the background and stops itself.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: stoptest <seconds>")
		os.Exit(1)
	}
	secs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}

	time.Sleep(time.Duration(secs) * time.Second)
	// Ignore SIGTSTP so we can send it to ourselves
	signal.Ignore(syscall.SIGTSTP)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTSTP)
}
