package shell

import (
	"os"
	"path/filepath"
	"strings"
)

// =============================================================================
// Interface & context
// =============================================================================

// CompletionContext holds shell state that completers may need.
// Extended as new job-aware or state-aware completers are added.
type CompletionContext struct {
	// TODO: add Jobs []*Job when job-aware completers are implemented
}

// Completer is the interface every per-command completer must satisfy.
// args is everything typed after the command, current is the word being completed.
type Completer interface {
	Complete(args []string, current string, ctx CompletionContext) []string
}

// =============================================================================
// Registry — maps command names to their Completer
// =============================================================================

type CompletionRegistry struct {
	completers map[string]Completer
}

func NewCompletionRegistry() *CompletionRegistry {
	r := &CompletionRegistry{completers: make(map[string]Completer)}

	r.Register("cd", &DirCompleter{})
	r.Register("git", &GitCompleter{})

	// TODO: register per-command completers here as they are built:
	// r.Register("fg", &JobCompleter{})
	// r.Register("bg", &JobCompleter{})

	return r
}

func (r *CompletionRegistry) Register(command string, c Completer) {
	r.completers[command] = c
}

func (r *CompletionRegistry) Get(command string) (Completer, bool) {
	c, ok := r.completers[command]
	return c, ok
}

// =============================================================================
// ShellCompleter — the single entry point readline calls on every Tab press
// =============================================================================

type ShellCompleter struct {
	registry     *CompletionRegistry
	builtins     []string
	pathBinaries []string // cached on startup, see scanPath
}

func NewShellCompleter(registry *CompletionRegistry) *ShellCompleter {
	sc := &ShellCompleter{
		registry: registry,
		builtins: []string{"cd", "pwd", "exit", "history", "jobs", "fg", "bg"},
	}
	sc.pathBinaries = sc.scanPath()
	return sc
}

// Do is called by readline on every Tab press.
// It routes to command completion, a registered completer, or file completion.
func (sc *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	if pos == 0 || strings.TrimSpace(string(line[:pos])) == "" {
		return nil, 0
	}

	input := string(line[:pos])
	words := strings.Fields(input)

	completingCommand := len(words) == 0 || (len(words) == 1 && !strings.HasSuffix(input, " "))
	if completingCommand {
		current := ""
		if len(words) == 1 {
			current = words[0]
		}
		return toReadlineCandidates(sc.completeCommands(current), len(current))
	}

	command := words[0]
	current := ""
	if !strings.HasSuffix(input, " ") {
		current = words[len(words)-1]
	}

	ctx := CompletionContext{}

	if completer, ok := sc.registry.Get(command); ok {
		candidates := completer.Complete(words[1:], current, ctx)
		return toReadlineCandidates(candidates, len(current))
	}

	return toReadlineCandidates(completeFiles(current), len(current))
}

// completeCommands returns builtins and $PATH binaries matching prefix.
func (sc *ShellCompleter) completeCommands(prefix string) []string {
	var candidates []string
	for _, b := range sc.builtins {
		if strings.HasPrefix(b, prefix) {
			candidates = append(candidates, b)
		}
	}
	for _, name := range sc.pathBinaries {
		if strings.HasPrefix(name, prefix) {
			candidates = append(candidates, name)
		}
	}
	return candidates
}

// scanPath scans all directories in $PATH and returns their executable binaries.
// Called once on startup to avoid rescanning on every Tab press.
func (sc *ShellCompleter) scanPath() []string {
	seen := make(map[string]bool)
	for _, b := range sc.builtins {
		seen[b] = true
	}
	var binaries []string
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if seen[name] {
				continue
			}
			if isExecutable(filepath.Join(dir, name)) {
				binaries = append(binaries, name)
				seen[name] = true
			}
		}
	}
	return binaries
}

// =============================================================================
// Helpers
// =============================================================================

// completeFiles returns file and directory names matching prefix.
// Handles tilde expansion and converts results back to ~ form.
func completeFiles(prefix string) []string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(prefix, "~/") {
		prefix = home + prefix[1:]
	}

	dir, filePrefix := filepath.Split(prefix)
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var candidates []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, filePrefix) {
			continue
		}
		candidate := filepath.Join(dir, name)
		if dir == "." {
			candidate = name
		}
		if e.IsDir() {
			candidate += "/"
		}
		if home != "" && strings.HasPrefix(candidate, home) {
			candidate = "~" + candidate[len(home):]
		}
		candidates = append(candidates, candidate)
	}
	return candidates
}

// isExecutable reports whether the file at path has execute permission bits set.
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// toReadlineCandidates converts candidates to the format readline expects.
// readline already has the typed prefix, so we strip it and return only the suffix.
func toReadlineCandidates(candidates []string, length int) ([][]rune, int) {
	result := make([][]rune, len(candidates))
	for i, c := range candidates {
		result[i] = []rune(c[length:])
	}
	return result, length
}

// =============================================================================
// Concrete completers — add new ones here and register them in NewCompletionRegistry
// =============================================================================

// DirCompleter completes only directories. Registered for "cd".
type DirCompleter struct{}

func (d *DirCompleter) Complete(args []string, current string, ctx CompletionContext) []string {
	home, _ := os.UserHomeDir()
	var candidates []string
	for _, candidate := range completeFiles(current) {
		realPath := candidate
		if home != "" && strings.HasPrefix(candidate, "~/") {
			realPath = home + candidate[1:]
		}
		info, err := os.Stat(realPath)
		if err == nil && info.IsDir() {
			candidates = append(candidates, candidate)
		}
	}
	return candidates
}

// GitCompleter completes git subcommands. Registered for "git".
type GitCompleter struct{}

func (g *GitCompleter) Complete(args []string, current string, ctx CompletionContext) []string {
	// TODO: word 1+ completion, e.g. `git checkout <Tab>` show branch names via `git branch`
	subcommands := []string{
		"add", "commit", "push", "pull", "checkout", "branch",
		"status", "log", "diff", "merge", "rebase", "stash",
	}
	var candidates []string
	for _, s := range subcommands {
		if strings.HasPrefix(s, current) {
			candidates = append(candidates, s)
		}
	}
	return candidates
}
