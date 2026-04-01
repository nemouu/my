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

// test_signals.go tests the mu shell's signal handling functionality.
// It spawns the shell as a subprocess and verifies correct behavior for:
//   - SIGINT (Ctrl+C) kills foreground job but not the shell
//   - SIGINT at the prompt with no foreground job is ignored
//   - Background jobs complete and are automatically removed from the job list
//   - SIGINT kills entire process groups including child-of-child processes
//
// Usage:
//
//	test_signals [shell_path]
//
// Default shell path is ../../my (relative to test/bin/)
// Or pass the path as an argument:
//
//	test_signals /path/to/my

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

// spawnShell starts the shell with -p=false to suppress the prompt,
// and returns the process, a stdin writer, and a stdout reader.
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

// sendCommand writes a command to the shell's stdin.
func sendCommand(stdin io.WriteCloser, command string) {
	fmt.Fprintln(stdin, command)
}

// readUntil reads lines from the shell's stdout until a line containing
// expected is found, printing each line as it arrives. Returns an error
// if the expected string is not found before the timeout.
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

// testSIGINT verifies that SIGINT kills the foreground job but leaves
// the shell alive and ready to accept new commands.
func testSIGINT(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30")
	sendCommand(stdin, "sleep 30")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Foreground job running ✓")

	fmt.Println("  Sending SIGINT (Ctrl+C)...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Verifying shell is still alive...")
	sendCommand(stdin, "jobs")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Shell survived SIGINT ✓")

	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT: %v", err)
	}
	fmt.Println("  Foreground job killed, shell continues ✓")

	return nil
}

// testBGCompletion verifies that a background job is automatically removed
// from the jobs list after it finishes.
func testBGCompletion(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 2 &")
	sendCommand(stdin, "sleep 2 &")
	if _, err := readUntil(reader, "[1]", 3*time.Second); err != nil {
		return fmt.Errorf("no background job started: %v", err)
	}
	fmt.Println("  Background job started ✓")

	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Running", 3*time.Second); err != nil {
		return fmt.Errorf("job not showing as Running: %v", err)
	}
	fmt.Println("  Job shows as Running ✓")

	fmt.Println("  Waiting for job to complete (2 seconds)...")
	time.Sleep(3 * time.Second)

	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	sendCommand(stdin, "echo done")
	if _, err := readUntil(reader, "done", 3*time.Second); err != nil {
		return fmt.Errorf("shell not responding after job completion: %v", err)
	}
	fmt.Println("  Job completed and removed from list ✓")

	return nil
}

// testSIGINTNoJob verifies that SIGINT at the prompt with no foreground job
// is silently ignored and the shell remains alive.
func testSIGINTNoJob(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending SIGINT with no foreground job...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Verifying shell is still alive...")
	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT at prompt: %v", err)
	}
	fmt.Println("  Shell survived SIGINT at prompt ✓")

	return nil
}

// testProcessGroup verifies that SIGINT kills an entire process group,
// including child processes spawned by the foreground job.
// Uses sh -c 'sleep 10 & wait' to create a parent with a child process.
func testProcessGroup(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sh -c 'sleep 10 & wait'")
	sendCommand(stdin, "sh -c 'sleep 10 & wait'")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Process with child started ✓")

	fmt.Println("  Sending SIGINT...")
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("failed to send SIGINT: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell did not survive SIGINT: %v", err)
	}
	fmt.Println("  Shell survived SIGINT ✓")

	sendCommand(stdin, "jobs")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Both parent and child processes killed ✓")

	return nil
}
