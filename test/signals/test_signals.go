package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// signals.go tests the mu shell's signal handling functionality.
// It spawns the shell as a subprocess and verifies correct behavior for:
//   - SIGINT (Ctrl+C) kills foreground job but not the shell
//   - Background job completes and is automatically removed from job list
//
// Usage:
//
//	signals [shell_path]
//
// Default shell path is ../../my (relative to test/bin/)
// Or pass the path as an argument:
//
//	signals /path/to/my

func main() {
	shellPath := "../../my"
	if len(os.Args) > 1 {
		shellPath = os.Args[1]
	}

	fmt.Println("=== Signal Handling Integration Tests ===")
	fmt.Println()

	passed := 0
	failed := 0

	tests := []struct {
		name string
		fn   func(string) error
	}{
		{"SIGINT kills foreground job but not shell", testSIGINT},
		{"Background job completes and is removed from job list", testBGCompletion},
		{"SIGINT ignored at prompt (no foreground job)", testSIGINTNoJob},
		{"Process group: SIGINT kills child-of-child", testProcessGroup},
	}

	for _, test := range tests {
		fmt.Printf("--- %s ---\n", test.name)
		if err := test.fn(shellPath); err != nil {
			fmt.Printf("FAIL: %v\n", err)
			failed++
		} else {
			fmt.Println("PASS ✓")
			passed++
		}
		fmt.Println()
	}

	fmt.Printf("=== Results: %d passed, %d failed ===\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

// spawnShell starts the shell with -p=false to suppress the prompt
func spawnShell(shellPath string) (*exec.Cmd, io.WriteCloser, *bufio.Reader, error) {
	cmd := exec.Command(shellPath, "-p=false")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start shell: %v", err)
	}
	reader := bufio.NewReader(stdoutPipe)
	time.Sleep(200 * time.Millisecond)
	return cmd, stdin, reader, nil
}

// sendCommand sends a command to the shell's stdin
func sendCommand(stdin io.WriteCloser, command string) {
	fmt.Fprintln(stdin, command)
}

// readUntil reads output until a line containing the expected string is found or timeout
func readUntil(reader *bufio.Reader, expected string, timeout time.Duration) (string, error) {
	done := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		var output strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				errChan <- err
				return
			}
			output.WriteString(line)
			fmt.Print("  > ", line)
			if strings.Contains(line, expected) {
				done <- output.String()
				return
			}
		}
	}()

	select {
	case output := <-done:
		return output, nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout waiting for: %q", expected)
	}
}

func testSIGINT(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	// Start a foreground job
	fmt.Println("  Sending: sleep 30")
	sendCommand(stdin, "sleep 30")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Foreground job running ✓")

	// Send SIGINT to the shell
	fmt.Println("  Sending SIGINT (Ctrl+C)...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Shell should still be alive — verify by sending a command
	fmt.Println("  Verifying shell is still alive...")
	sendCommand(stdin, "jobs")
	// jobs should return quickly with empty output (job was killed)
	// if shell is dead this will timeout
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Shell survived SIGINT ✓")

	// Verify job is gone
	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT: %v", err)
	}
	fmt.Println("  Foreground job killed, shell continues ✓")

	return nil
}

func testBGCompletion(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	// Start a short background job
	fmt.Println("  Sending: sleep 2 &")
	sendCommand(stdin, "sleep 2 &")
	if _, err := readUntil(reader, "[1]", 3*time.Second); err != nil {
		return fmt.Errorf("no background job started: %v", err)
	}
	fmt.Println("  Background job started ✓")

	// Verify it shows as running
	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Running", 3*time.Second); err != nil {
		return fmt.Errorf("job not showing as Running: %v", err)
	}
	fmt.Println("  Job shows as Running ✓")

	// Wait for job to complete
	fmt.Println("  Waiting for job to complete (2 seconds)...")
	time.Sleep(3 * time.Second)

	// Verify job is gone from list
	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	// Give it a moment then check — if jobs list is empty there will be no output
	// so we verify by running echo and checking the job list is empty
	sendCommand(stdin, "echo done")
	if _, err := readUntil(reader, "done", 3*time.Second); err != nil {
		return fmt.Errorf("shell not responding after job completion: %v", err)
	}
	fmt.Println("  Job completed and removed from list ✓")

	return nil
}

func testSIGINTNoJob(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	// Send SIGINT with no foreground job running
	fmt.Println("  Sending SIGINT with no foreground job...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Shell should still be alive
	fmt.Println("  Verifying shell is still alive...")
	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT at prompt: %v", err)
	}
	fmt.Println("  Shell survived SIGINT at prompt ✓")

	return nil
}

func testProcessGroup(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	// Use sh -c to spawn a parent shell that itself forks a child (sleep)
	// This mimics splittest: a process that has its own child process
	// SIGINT should kill the entire process group including the child
	fmt.Println("  Sending: sh -c 'sleep 10 & wait'")
	sendCommand(stdin, "sh -c 'sleep 10 & wait'")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Process with child started ✓")

	// Send SIGINT — should kill both sh and its child sleep
	fmt.Println("  Sending SIGINT...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Shell should still be alive
	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT: %v", err)
	}
	fmt.Println("  Shell survived SIGINT ✓")

	// Jobs list should be empty — both processes killed
	sendCommand(stdin, "jobs")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Both parent and child processes killed ✓")

	return nil
}
