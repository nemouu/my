package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CompletionContext holds shell state that completers may need.
type CompletionContext struct {
	// TODO: add Jobs []*Job when job-aware completers are implemented
}

// Completer is the interface for per-command tab completion.
// args is everything typed so far, current is the word being completed.
type Completer interface {
	Complete(args []string, current string, ctx CompletionContext) []string
}

// CompletionRegistry maps command names to their Completer.
type CompletionRegistry struct {
	completers map[string]Completer
}

func NewCompletionRegistry() *CompletionRegistry {
	r := &CompletionRegistry{completers: make(map[string]Completer)}

	r.Register("cd", &DirCompleter{})
	r.Register("git", &GitCompleter{})

	// TODO: register per-command completers here as they are built, e.g.:
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

// ShellCompleter implements readline.AutoCompleter.
// It is the single entry point readline calls on every Tab press.
type ShellCompleter struct {
	registry     *CompletionRegistry
	builtins     []string
	pathBinaries []string
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
// line is the full input buffer, pos is the cursor position.
func (sc *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {

	if pos == 0 || strings.TrimSpace(string(line[:pos])) == "" {
		return nil, 0
	}

	// Only consider the input up to the cursor
	input := string(line[:pos])
	words := strings.Fields(input)

	// Are we completing the first word (the command)?
	completingCommand := len(words) == 0 || (len(words) == 1 && !strings.HasSuffix(input, " "))

	if completingCommand {
		current := ""
		if len(words) == 1 {
			current = words[0]
		}
		return toReadlineCandidates(sc.completeCommands(current), len(current))
	}

	// We are completing an argument
	command := words[0]
	current := ""
	if !strings.HasSuffix(input, " ") {
		current = words[len(words)-1]
	}

	ctx := CompletionContext{}

	// Look up a registered completer for this command
	if completer, ok := sc.registry.Get(command); ok {
		candidates := completer.Complete(words[1:], current, ctx)
		return toReadlineCandidates(candidates, len(current))
	}

	// Default: file completion
	candidates := completeFiles(current)
	return toReadlineCandidates(candidates, len(current))
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

// this helper function now scans paths and stores where it was
// for later use
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

// completeFiles returns file/directory names matching prefix.
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
	fmt.Println(candidates)
	return candidates
}

// isExecutable reports whether the file at path is executable.
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// toReadlineCandidates converts string candidates to the format readline expects.
// length is how many characters of the current word readline should replace.
func toReadlineCandidates(candidates []string, length int) ([][]rune, int) {
	result := make([][]rune, len(candidates))
	for i, c := range candidates {
		// readline appends these to the remaining suffix, so strip the prefix
		result[i] = []rune(c[length:])
	}
	return result, length
}

// DirCompleter completes only directories. Register for "cd".
type DirCompleter struct{}

func (d *DirCompleter) Complete(args []string, current string, ctx CompletionContext) []string {
	currCandidates := completeFiles(current)

	var actualCandidates []string

	// filter the current directory entries for folders
	for _, candidate := range currCandidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			actualCandidates = append(actualCandidates, candidate)
		}
	}

	return actualCandidates
}

// GitCompleter completes git subcommands. Register for "git".
type GitCompleter struct{}

func (g *GitCompleter) Complete(args []string, current string, ctx CompletionContext) []string {
	// TODO: return subcommands for word 0 works but smarter completion for word 1+
	// has to be done e.g. `git checkout <Tab>` → branch names via `git branch`
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
