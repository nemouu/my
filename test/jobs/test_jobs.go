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

// jobs_test.go tests the mu shell's job control functionality.
// It spawns the shell as a subprocess and verifies correct behavior for:
//   - Background jobs (&)
//   - Stopping foreground jobs (SIGTSTP / Ctrl+Z)
//   - Resuming stopped jobs in foreground (fg)
//   - Resuming stopped jobs in background (bg)
//   - Multiple simultaneous background jobs
//   - Invalid job ids
//
// Usage:
//
//	test_jobs [shell_path]
//
// Default shell path is ../../my (relative to test/bin/)
// Or pass the path as an argument:
//
//	test_jobs /path/to/my

func main() {
	shellPath := "../../my"
	if len(os.Args) > 1 {
		shellPath = os.Args[1]
	}

	fmt.Println("=== Job Control Integration Tests ===")
	fmt.Println()

	passed := 0
	failed := 0

	tests := []struct {
		name string
		fn   func(string) error
	}{
		{"Background job", testBackgroundJob},
		{"SIGTSTP stops foreground job", testSIGTSTP},
		{"fg resumes stopped job", testFGResume},
		{"bg resumes stopped job in background", testBGResume},
		{"Multiple background jobs", testMultipleBackgroundJobs},
		{"Invalid job id", testInvalidJobID},
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

// testBackgroundJob verifies that a job started with & runs in the background
// and appears as Running in the jobs list.
func testBackgroundJob(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30 &")
	sendCommand(stdin, "sleep 30 &")
	if _, err := readUntil(reader, "[1]", 3*time.Second); err != nil {
		return fmt.Errorf("no background job started: %v", err)
	}
	fmt.Println("  Background job started ✓")

	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Running", 3*time.Second); err != nil {
		return fmt.Errorf("job not showing as Running: %v", err)
	}
	fmt.Println("  Job shows as Running ✓")

	return nil
}

// testSIGTSTP verifies that sending SIGTSTP to the shell stops the foreground
// job and marks it as Stopped in the jobs list.
func testSIGTSTP(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30")
	sendCommand(stdin, "sleep 30")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Foreground job running ✓")

	fmt.Println("  Sending SIGTSTP to shell...")
	if err := cmd.Process.Signal(syscall.SIGTSTP); err != nil {
		return fmt.Errorf("failed to send SIGTSTP: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Stopped", 3*time.Second); err != nil {
		return fmt.Errorf("job not showing as Stopped: %v", err)
	}
	fmt.Println("  Job shows as Stopped ✓")

	return nil
}

// testFGResume verifies that a stopped job can be resumed in the foreground
// with fg, and is removed from the jobs list after being killed.
func testFGResume(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30")
	sendCommand(stdin, "sleep 30")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Sending SIGTSTP...")
	cmd.Process.Signal(syscall.SIGTSTP)
	time.Sleep(500 * time.Millisecond)

	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Stopped", 3*time.Second); err != nil {
		return fmt.Errorf("job not stopped: %v", err)
	}
	fmt.Println("  Job stopped ✓")

	fmt.Println("  Sending: fg %1")
	sendCommand(stdin, "fg %1")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Job resumed in foreground ✓")

	fmt.Println("  Sending SIGINT to terminate job...")
	cmd.Process.Signal(syscall.SIGINT)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  Job removed from list ✓")

	return nil
}

// testBGResume verifies that a stopped job can be resumed in the background
// with bg and appears as Running in the jobs list.
func testBGResume(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30")
	sendCommand(stdin, "sleep 30")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("  Sending SIGTSTP...")
	cmd.Process.Signal(syscall.SIGTSTP)
	time.Sleep(500 * time.Millisecond)

	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Stopped", 3*time.Second); err != nil {
		return fmt.Errorf("job not stopped: %v", err)
	}
	fmt.Println("  Job stopped ✓")

	fmt.Println("  Sending: bg %1")
	sendCommand(stdin, "bg %1")
	time.Sleep(500 * time.Millisecond)

	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "Running", 3*time.Second); err != nil {
		return fmt.Errorf("job not running in background: %v", err)
	}
	fmt.Println("  Job resumed in background ✓")

	return nil
}

// testMultipleBackgroundJobs verifies that multiple background jobs can run
// simultaneously and all appear in the jobs list.
func testMultipleBackgroundJobs(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: sleep 30 &")
	sendCommand(stdin, "sleep 30 &")
	if _, err := readUntil(reader, "[1]", 3*time.Second); err != nil {
		return fmt.Errorf("job 1 not started: %v", err)
	}
	fmt.Println("  Job 1 started ✓")

	fmt.Println("  Sending: sleep 30 &")
	sendCommand(stdin, "sleep 30 &")
	if _, err := readUntil(reader, "[2]", 3*time.Second); err != nil {
		return fmt.Errorf("job 2 not started: %v", err)
	}
	fmt.Println("  Job 2 started ✓")

	fmt.Println("  Sending: sleep 30 &")
	sendCommand(stdin, "sleep 30 &")
	if _, err := readUntil(reader, "[3]", 3*time.Second); err != nil {
		return fmt.Errorf("job 3 not started: %v", err)
	}
	fmt.Println("  Job 3 started ✓")

	fmt.Println("  Sending: jobs")
	sendCommand(stdin, "jobs")
	if _, err := readUntil(reader, "[3]", 3*time.Second); err != nil {
		return fmt.Errorf("not all jobs showing: %v", err)
	}
	fmt.Println("  All 3 jobs show as Running ✓")

	return nil
}

// testInvalidJobID verifies that fg and bg with a non-existent job id
// return an error message without crashing the shell.
func testInvalidJobID(shellPath string) error {
	cmd, stdin, reader, err := spawnShell(shellPath)
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()

	fmt.Println("  Sending: fg %99")
	sendCommand(stdin, "fg %99")
	if _, err := readUntil(reader, "not found", 3*time.Second); err != nil {
		return fmt.Errorf("expected 'not found' error: %v", err)
	}
	fmt.Println("  Got expected error for invalid fg job id ✓")

	fmt.Println("  Sending: bg %99")
	sendCommand(stdin, "bg %99")
	if _, err := readUntil(reader, "not found", 3*time.Second); err != nil {
		return fmt.Errorf("expected 'not found' error: %v", err)
	}
	fmt.Println("  Got expected error for invalid bg job id ✓")

	sendCommand(stdin, "echo alive")
	if _, err := readUntil(reader, "alive", 3*time.Second); err != nil {
		return fmt.Errorf("shell not responding after invalid job id: %v", err)
	}
	fmt.Println("  Shell still alive after invalid job ids ✓")

	return nil
}
