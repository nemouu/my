package shell

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Command represents a parsed command.
type Command struct {
	Args []string // Command and arguments
	Bg   bool     // Run in background?
}

// Parse splits a line of input into a Command.
func parse(line string) (*Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, errors.New("empty input")
	}

	// Remove trailing '&' for background jobs
	bg := false
	if strings.HasSuffix(line, "&") {
		bg = true
		line = strings.TrimSpace(line[:len(line)-1])
	}

	// Get the arguments from the input (also handles quoted strings)
	var args []string
	var currentArg strings.Builder
	inside := false

	for i := 0; i < len(line); i++ {

		// switch between true and false depending on whether we are inside or not
		if (line[i] == '"' || line[i] == '\'') && !inside {
			inside = true
			continue
		} else if (line[i] == '"' || line[i] == '\'') && inside {
			inside = false
			continue
		}

		// On space append the current argument unless inside quotations
		if line[i] == ' ' && !inside {
			args = append(args, currentArg.String())
			currentArg.Reset()
			continue
		}

		currentArg.WriteByte(line[i])
	}

	// If still inside quote return error, otherwise append the last argument
	if inside {
		return nil, errors.New("unclosed quote")
	}
	args = append(args, currentArg.String())

	// expand tilde in each argument
	for i, arg := range args {
		args[i] = expandTilde(arg)
	}

	// expand globs in each argument
	args = expandGlobs(args)

	if len(args) == 0 {
		return nil, errors.New("no command")
	}

	return &Command{Args: args, Bg: bg}, nil
}

// Helper function for tilde character
func expandTilde(line string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return line
	}
	// replace ~ at start of line
	if strings.HasPrefix(line, "~/") {
		return home + line[1:]
	}
	if line == "~" {
		return home
	}
	return line
}

// Helper function for globs characters
func expandGlobs(args []string) []string {
	var expanded []string
	for _, arg := range args {
		if strings.Contains(arg, "*") || strings.Contains(arg, "?") || strings.Contains(arg, "[") {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return args
			}
			if len(matches) == 0 {
				// no matches, keep the pattern as-is (bash behavior)
				expanded = append(expanded, arg)
			} else {
				expanded = append(expanded, matches...)
			}
		} else {
			// no glob characters, keep the argument unchanged
			expanded = append(expanded, arg)
		}
	}
	return expanded
}
