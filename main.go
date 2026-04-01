package main

// μ (My) Shell - AI-Powered Shell with Local LLM Integration
//
// This is the main entry point for the μ shell. It initializes the
// configuration, sets up the environment, and starts the main REPL loop.
//
// TODO: Implement main function as described in TODO.md Milestone 1
// - Load configuration
// - Initialize Ollama client
// - Handle command-line flags (--init, --version, etc.)
// - Start the shell REPL

import (
	"fmt"
	"os"

	"flag"

	"github.com/nemouu/my/shell"
)

func main() {
	emitPrompt := flag.Bool("p", true, "emit prompt")
	flag.Parse()

	if err := shell.Run(*emitPrompt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
