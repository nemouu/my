package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// inttest sleeps for a specified number of seconds, then sends SIGINT to itself.
// This tests the shell's handling of interrupted jobs (e.g., with Ctrl+C).
//
// Usage:
//
//	inttest <seconds>
//
// Example:
//
//	inttest 4   # Sleeps for 4 seconds, then interrupts itself.
//	inttest 2 & # Runs in the background and interrupts itself.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: inttest <seconds>")
		os.Exit(1)
	}
	secs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}

	time.Sleep(time.Duration(secs) * time.Second)
	// Ignore SIGINT so we can send it to ourselves
	signal.Ignore(syscall.SIGINT)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
}
