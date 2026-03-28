package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// spintest simulates a long-running process by sleeping for a specified number of seconds,
// printing its progress each second. It is used to test foreground/background job control
// and signal handling (e.g., SIGINT, SIGTSTP) in the shell.
//
// Usage:
//
//	spintest <seconds>
//
// Example:
//
//	spintest 5   # Sleeps for 5 seconds, printing progress.
//	spintest 3 & # Runs in the background.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: spintest <seconds>")
		os.Exit(1)
	}
	secs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}
	for i := 0; i < secs; i++ {
		fmt.Printf("Sleeping %d/%d\n", i+1, secs)
		time.Sleep(1 * time.Second)
	}
}
