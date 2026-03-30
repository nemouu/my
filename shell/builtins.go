package shell

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

// This file will contain implementations of built-in shell commands
// that cannot be external programs (cd, exit, pwd, etc.) as well as
// AI-powered commands (ask, commit, code).
//
// TODO: Later add AI commands
// - ask: Query local LLM
// - commit: AI-powered git commit messages
// - code: Interactive code mode

var previousDir string

// IsBuiltin checks if a command name is a built-in command
func isBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "exit", "pwd", "history", "jobs", "fg", "bg":
		return true
	}
	return false
}

// ExecuteBuiltin executes a built-in command
func executeBuiltin(args []string) error {
	switch args[0] {
	case "cd":
		// Save the current directory in global variable
		current, err := os.Getwd()
		if err != nil {
			return err
		}

		// No second argument, go to home directory
		if len(args) < 2 {
			// go to home directory
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			previousDir = current
			return os.Chdir(home)
		}

		// check second argument
		if args[1] == "-" {
			if previousDir == "" {
				return errors.New("cd: no previous directory")
			}
			target := previousDir
			previousDir = current
			return os.Chdir(target)
		}
		previousDir = current
		return os.Chdir(args[1])
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(dir)
		return nil
	case "exit":
		if len(args) > 1 {
			status, err := strconv.Atoi((args[1]))
			if err != nil {
				return err
			}
			os.Exit(status)
		} else {
			os.Exit(0)
		}
	case "history":
		data, err := os.ReadFile(os.Getenv("HOME") + "/.config/mu/history")
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	case "jobs":
		listJobs()
		return nil
	case "fg":
		// Check if input has correct format
		if len(args) < 2 || args[1][0] != '%' {
			return errors.New("fg: usage: fg %<jobid>")
		}

		// Parse the job id from args
		jid, err := strconv.Atoi(args[1][1:])
		if err != nil {
			return err
		}

		// Find the job
		job, err := getJobByJid(jid)
		if err != nil {
			return err
		}

		// Send SIGCONT
		err = syscall.Kill(-job.pid, syscall.SIGCONT)
		if err != nil {
			return err
		}

		// Update state of job
		job.state = FG

		// Wait for it since its now foreground, delete it when done
		var status syscall.WaitStatus
		pid, err := syscall.Wait4(job.pid, &status, 0, nil)
		if err != nil {
			return err
		}
		deleteJob(pid)

		return nil
	case "bg":
		// Check if input has correct format
		if len(args) < 2 || args[1][0] != '%' {
			return errors.New("bg: usage: bg %<jobid>")
		}

		// Parse the job id from args
		jid, err := strconv.Atoi(args[1][1:])
		if err != nil {
			return err
		}

		// Find the job
		job, err := getJobByJid(jid)
		if err != nil {
			return err
		}

		// Send SIGCONT
		err = syscall.Kill(-job.pid, syscall.SIGCONT)
		if err != nil {
			return err
		}

		// Update state of job
		job.state = BG

		// Print cmdline for info to user
		fmt.Printf("[%d] %s", jid, job.cmdline)

		return nil
	default:
		return nil
	}
	return nil
}
