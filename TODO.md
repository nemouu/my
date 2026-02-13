# μ (My) Shell - Development Roadmap & TODO

## Project Overview

Building a real Unix shell from scratch in Go with integrated local LLM capabilities. This is both a practical tool and a long-term learning project for systems programming.

**Timeline:**
- MVP (basic shell + AI): 4 weeks
- Feature complete v1.0: 3 months
- Advanced features: 6-12 months ongoing

---

## Design Reference

### Philosophy

Build a real shell from the ground up as a long-term systems programming learning project, while delivering immediate practical value through AI integration. Start simple, grow organically.

### What Makes μ Different from Existing AI Shell Tools

| Existing AI shell tools | μ approach |
|---|---|
| Shell plugins or CLI wrappers (not actual shells) | Actual shell you run as your primary environment |
| Require hotkeys and explicit invocation | AI commands are built-in (like `cd` or `exit`) |
| Context switching between shell and AI | Automatic context from current directory and session |
| Usually cloud-dependent | Local-first by design (Ollama) |

**Inspiration:** Like having Claude Code's capabilities built into your shell, but using local models and maintaining privacy. Technically modeled on standard Unix shell architecture (similar to dash/sh) with AI as a native feature rather than a plugin.

### Command Execution Model

```
User input
    ↓
Built-in? (cd, exit, ask, commit, code)
    ↓ YES → Handle in Go
    ↓ NO  → External command (gcc, git, ls, etc.) → fork() + exec()
```

**Built-ins:** `cd`, `exit`, `pwd`, `ask`, `commit`, `code`, `history`, `session`

### Session & Config Storage

```
~/.config/my/
├── config.yaml
├── sessions/
│   ├── {sha256hash}.json   # per-directory session
│   └── ...
└── history
```

### Default Configuration Reference

```yaml
ollama:
  host: "http://localhost:11434"

models:
  ask: "qwen2.5:7b"
  code: "qwen2.5-coder:7b"
  commit: "deepseek-coder:6.7b"

context:
  max_tokens: 4096
  warning_threshold: 0.8
  auto_load_files: true

session:
  auto_save: true
  cleanup_days: 7

shell:
  prompt: "μ"
  show_git_branch: true
  history_size: 1000
```

### Technology Choices

**Language:** Go — easy process management (exec package), single binary distribution, fast startup, compiled performance.

**Dependencies:**
- Ollama Go API client
- gopkg.in/yaml.v3 (YAML parsing)
- Cobra (CLI framework — TBD if needed)
- github.com/chzyer/readline (TBD)

### Non-Goals

- Not trying to replace bash/zsh for power users (initially)
- Not a full IDE or development environment
- Not trying to compete with cloud AI quality
- Not implementing bash/zsh compatibility (fresh start)

### Use Cases

**Good fit:** learning codebases, quick git commits, exploring new projects, solo dev with AI pair programming, systems programming practice.

**Not ideal for:** complex shell scripts, production deployment automation, heavy piping workflows (will improve over time).

---

## Phase 0: Setup & Prerequisites

### Environment Setup
- [ ] Install Go (1.21+): `https://go.dev/dl/`
- [ ] Install Ollama: `https://ollama.ai`
- [ ] Pull required models:
  ```bash
  ollama pull qwen2.5:7b
  ollama pull qwen2.5-coder:7b
  ollama pull deepseek-coder:6.7b
  ```
- [ ] Verify Ollama is running: `ollama list`
- [ ] Test a model: `ollama run qwen2.5:7b "hello"`

### Project Initialization
- [ ] Create project directory: `mkdir my && cd my`
- [ ] Initialize Go module: `go mod init github.com/yourusername/my`
- [ ] Install dependencies:
  ```bash
  go get github.com/ollama/ollama/api
  go get gopkg.in/yaml.v3
  # Maybe: go get github.com/spf13/cobra (TBD if needed)
  ```
- [ ] Create directory structure:
  ```
  my/
  ├── main.go
  ├── shell/
  │   ├── repl.go
  │   ├── parser.go
  │   ├── executor.go
  │   └── builtins.go
  ├── ai/
  │   ├── client.go
  │   ├── modes.go
  │   └── session.go
  ├── git/
  │   └── integration.go
  ├── config/
  │   └── config.go
  ├── go.mod
  ├── go.sum
  ├── README.md
  └── TODO.md
  ```

### Learning Resources
- [ ] Read about Unix fork/exec model
- [ ] Review Go's `os/exec` package documentation
- [ ] Study simple shell implementations (dash source, toy shells)

---

## Milestone 1: Basic Shell (Week 1)

**Goal:** Build a minimal working shell that can execute any command

### Core REPL Implementation
- [ ] Create main.go with entry point
- [ ] Implement basic REPL in `shell/repl.go`:
  ```go
  for {
      printPrompt()
      input := readLine()
      execute(input)
  }
  ```
- [ ] Implement `readLine()` using `bufio.Scanner`
- [ ] Create simple prompt: `μ ~/current/dir> `

### Input Parsing
- [ ] Create `shell/parser.go`
- [ ] Implement `Parse(input string) []string`:
  - [ ] Split by whitespace
  - [ ] Handle quoted strings: `"hello world"` → single arg
  - [ ] Handle empty input
  - [ ] Trim whitespace
- [ ] Add error handling for malformed input

### Command Execution
- [ ] Create `shell/executor.go`
- [ ] Implement command execution logic:
  ```go
  func Execute(args []string) error {
      if isBuiltin(args[0]) {
          return executeBuiltin(args)
      }
      return executeExternal(args)
  }
  ```
- [ ] Implement `executeExternal()`:
  ```go
  cmd := exec.Command(args[0], args[1:]...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
  ```
- [ ] Handle command not found errors
- [ ] Handle Ctrl+C gracefully (signal handling)

### Built-in Commands
- [ ] Create `shell/builtins.go`
- [ ] Implement `cd <directory>`:
  - [ ] Change current directory with `os.Chdir()`
  - [ ] Handle `cd` (go to home)
  - [ ] Handle `cd -` (go to previous directory)
  - [ ] Update prompt to show new directory
- [ ] Implement `exit`:
  - [ ] Clean exit with `os.Exit(0)`
- [ ] Implement `pwd`:
  - [ ] Print current directory with `os.Getwd()`

### Testing
- [ ] Test basic commands:
  ```bash
  μ> ls
  μ> ls -la
  μ> gcc main.c -o main
  μ> git status
  μ> cd /tmp
  μ> pwd
  μ> exit
  ```
- [ ] Test error handling (invalid commands, missing files)
- [ ] Test Ctrl+C doesn't crash shell
- [ ] Test with various command-line tools

### Deliverable
Working shell that can:
- Execute any system command (gcc, git, flutter, ls, etc.)
- Navigate directories with `cd`
- Exit cleanly
- Handle errors gracefully

---

## Milestone 2: Configuration System (Week 2 Part 1)

**Goal:** Load configuration from YAML file

### Config Structure
- [ ] Create `config/config.go`
- [ ] Define config structs:
  ```go
  type Config struct {
      Ollama  OllamaConfig
      Models  ModelsConfig
      Context ContextConfig
      Session SessionConfig
      Shell   ShellConfig
      Prompts PromptsConfig
  }
  ```
- [ ] Create default config values
- [ ] Implement config validation

### Config Loading
- [ ] Implement `LoadConfig() (*Config, error)`:
  - [ ] Check `~/.config/my/config.yaml`
  - [ ] If missing, use defaults
  - [ ] Parse YAML with `gopkg.in/yaml.v3`
  - [ ] Validate required fields
- [ ] Add `--init` flag to generate default config:
  ```bash
  μ --init
  # Creates ~/.config/my/config.yaml with defaults
  ```
- [ ] Implement config reload (for development)

### Default Config Template
- [ ] Create default config with sensible values
- [ ] Add comments explaining each option
- [ ] Include all system prompt templates

### Testing
- [ ] Test config loading with valid YAML
- [ ] Test with missing config (uses defaults)
- [ ] Test with invalid YAML (error handling)
- [ ] Test `--init` flag creates config correctly

---

## Milestone 3: Ollama Integration & Ask Command (Week 2 Part 2)

**Goal:** Connect to Ollama and implement `ask` command

### Ollama Client
- [ ] Create `ai/client.go`
- [ ] Implement Ollama client wrapper:
  ```go
  type Client struct {
      baseURL string
      client  *http.Client
  }
  
  func (c *Client) Generate(model, prompt string) (string, error)
  ```
- [ ] Add connection testing
- [ ] Handle timeouts
- [ ] Handle Ollama not running error
- [ ] Add streaming support (for future)

### Session Management
- [ ] Create `ai/session.go`
- [ ] Define session structure:
  ```go
  type Session struct {
      Directory   string
      Messages    []Message
      LoadedFiles map[string]string
      TotalTokens int
      LastActive  time.Time
  }
  ```
- [ ] Implement session file path generation:
  ```go
  // Use SHA256 hash of directory path
  hash := sha256(directory)
  path := ~/.config/my/sessions/{hash}.json
  ```
- [ ] Implement `LoadSession(dir string) (*Session, error)`
- [ ] Implement `SaveSession(session *Session) error`
- [ ] Implement token counting (approximate from char count)

### Ask Command (Built-in)
- [ ] Add `ask` to built-in commands in `shell/builtins.go`
- [ ] Implement basic ask flow:
  ```go
  func executeAsk(query string) error {
      // 1. Get current directory
      // 2. Load or create session
      // 3. Check for file references in query
      // 4. Load referenced files
      // 5. Build prompt with context
      // 6. Call Ollama
      // 7. Display response
      // 8. Save session
  }
  ```
- [ ] Add auto-file detection (regex for file paths)
- [ ] Implement context building:
  - [ ] System prompt from config
  - [ ] Current directory info
  - [ ] Previous messages from session
  - [ ] Loaded file contents
- [ ] Add token usage warning (80% threshold)

### File Loading
- [ ] Implement file content reading
- [ ] Add file to session's loaded files
- [ ] Handle file not found errors
- [ ] Support reading multiple files
- [ ] Cache file contents in session

### Testing
- [ ] Test Ollama connection
- [ ] Test ask command:
  ```bash
  μ ~/project> ask what is 2+2?
  μ ~/project> ask explain main.c
  μ ~/project> ask compare utils.c and helper.c
  ```
- [ ] Test session persistence (ask, exit, restart, ask again)
- [ ] Test context warning at 80%
- [ ] Test with Ollama not running (error message)

### Deliverable
- `ask` command works with automatic context
- Sessions saved per directory
- File loading automatic from queries

---

## Milestone 4: Git Commit Integration (Week 3)

**Goal:** Implement AI-powered git commit message generation

### Git Integration
- [ ] Create `git/integration.go`
- [ ] Implement `GetStagedDiff() (string, error)`:
  ```bash
  git diff --staged
  ```
- [ ] Implement `HasStagedChanges() bool`
- [ ] Implement `GetCurrentBranch() string`
- [ ] Handle not in git repo error

### Commit Command
- [ ] Add `commit` to built-in commands
- [ ] Implement commit flow:
  ```go
  func executeCommit() error {
      // 1. Check for staged changes
      // 2. Get diff
      // 3. Generate commit message via Ollama
      // 4. Display message
      // 5. Prompt user: [C]ommit [E]dit [R]egenerate [A]bort
      // 6. Handle user choice
  }
  ```
- [ ] Implement user prompt with options:
  - [ ] C: Execute `git commit -m "<message>"`
  - [ ] E: Open in `$EDITOR`, then commit
  - [ ] R: Regenerate message
  - [ ] A: Abort without committing
- [ ] Add commit message prompt in config (conventional commits)

### System Prompt for Commits
- [ ] Create commit-specific prompt:
  ```
  Generate a conventional commit message.
  Format: <type>(<scope>): <description>
  
  Types: feat, fix, docs, style, refactor, test, chore
  Keep subject under 50 chars.
  Add bullet points for details if needed.
  ```
- [ ] Add few-shot examples for better quality

### Editor Integration
- [ ] Detect `$EDITOR` environment variable
- [ ] Create temp file with message
- [ ] Open editor
- [ ] Read edited message
- [ ] Commit with edited message

### Testing
- [ ] Test with single file change
- [ ] Test with multiple files
- [ ] Test with add/delete/modify operations
- [ ] Test all user options (C/E/R/A)
- [ ] Test with no staged changes (error)
- [ ] Test editor integration
- [ ] Verify conventional commit format

### Deliverable
```bash
μ ~/project> git add .
μ ~/project> commit
Generated commit message:
━━━━━━━━━━━━━━━━━━━━━━━
feat: Add user authentication

- Implement JWT middleware
- Add login/logout routes
━━━━━━━━━━━━━━━━━━━━━━━
[C]ommit [E]dit [R]egenerate [A]bort >
```

---

## Milestone 5: Code Mode (Week 4)

**Goal:** Interactive REPL for code-related queries

### Code Mode Implementation
- [ ] Create `ai/modes.go`
- [ ] Implement code mode REPL:
  ```go
  func EnterCodeMode() {
      for {
          input := readPrompt("code> ")
          if input == "exit" { break }
          handleCodeQuery(input)
      }
  }
  ```
- [ ] Add `code` built-in command to enter mode
- [ ] Use code-specific model from config
- [ ] Use code-specific system prompt

### Code Mode Features
- [ ] File explanation:
  ```
  code> explain divide() in utils.c
  ```
- [ ] Multi-file analysis:
  ```
  code> how do main.c and utils.c connect?
  ```
- [ ] Code generation:
  ```
  code> create helper.c with file I/O functions
  ```
- [ ] Code review suggestions:
  ```
  code> review auth.go for improvements
  ```

### Code Generation with Confirmation
- [ ] Detect code generation requests
- [ ] Extract code from markdown blocks in response
- [ ] Show code preview (first/last 20 lines if >50 lines)
- [ ] Warn if code >100 lines
- [ ] Implement per-session permission:
  ```
  Allow file creation this session? [y/N]
  ```
- [ ] Cache permission for session
- [ ] Prompt for filename if not specified
- [ ] Write file after confirmation

### Code Mode System Prompt
- [ ] Create teaching-focused prompt:
  ```
  You are a coding teacher and assistant.
  - Explain code clearly with examples
  - Suggest concrete improvements
  - Keep generated code concise (<100 lines)
  - Add helpful comments
  - Ask for clarification if ambiguous
  ```

### Context in Code Mode
- [ ] Maintain conversation history in code mode
- [ ] Track loaded files separately
- [ ] Token counting and warnings
- [ ] Clear context on mode exit (optional)

### Commands in Code Mode
- [ ] `exit` - Leave code mode
- [ ] `/clear` - Clear conversation context
- [ ] `/files` - List loaded files
- [ ] `/help` - Show available commands

### Testing
- [ ] Test code explanations
- [ ] Test multi-file context
- [ ] Test code generation
- [ ] Test file creation permission flow
- [ ] Test context limits in long sessions
- [ ] Test exiting and re-entering

### Deliverable
```bash
μ ~/project> code
code> explain the divide function in calc.c
[AI explains function]
code> can you improve error handling?
[AI suggests improvements with code]
code> create a test file for this
[Shows preview, asks permission, creates file]
code> exit
μ ~/project>
```

---

## Milestone 6: Shell History & Improvements (Month 2 Week 1)

**Goal:** Add command history and quality-of-life improvements

### Command History
- [ ] Implement history storage:
  ```go
  type History struct {
      commands []string
      maxSize  int
      file     string
  }
  ```
- [ ] Save history to `~/.config/my/history`
- [ ] Implement `history` built-in command
- [ ] Add up/down arrow navigation (readline library)
- [ ] Consider using: `github.com/chzyer/readline`

### Improved Prompt
- [ ] Show current directory (shortened if long)
- [ ] Optionally show git branch when in git repo
- [ ] Color support (green prompt, red if last command failed)
- [ ] Customize prompt in config

### Better Error Messages
- [ ] User-friendly error for Ollama not running:
  ```
  ⚠️  Cannot connect to Ollama at localhost:11434
  Is Ollama running? Start it with: ollama serve
  ```
- [ ] Better git error messages
- [ ] Suggest `my --init` if config missing

### Signal Handling
- [ ] Proper Ctrl+C handling (interrupt current command)
- [ ] Ctrl+D to exit shell
- [ ] Don't exit on Ctrl+C in REPL

### Testing
- [ ] Test history persists across sessions
- [ ] Test up/down arrow navigation
- [ ] Test prompt updates with directory changes
- [ ] Test error messages are helpful

---

## Milestone 7: Context Optimization (Month 2 Week 2-3)

**Goal:** Better context management and token optimization

### Token Counting
- [ ] Implement accurate token counting (or close approximation)
- [ ] Show token usage in verbose mode
- [ ] Add `/tokens` command in interactive modes

### Context Summarization
- [ ] When approaching token limit:
  - [ ] Keep last N messages (e.g., 10)
  - [ ] Keep all loaded files
  - [ ] Option: Summarize old messages into single message
- [ ] Implement summarization via LLM
- [ ] Make strategy configurable

### Session Management Commands
- [ ] `session info` - Show current session stats
- [ ] `session clear` - Clear context, keep session
- [ ] `session list` - List all sessions
- [ ] `session clean` - Delete old sessions

### File Management
- [ ] Track which files are loaded
- [ ] Unload files from context: `/unload utils.c`
- [ ] Reload file if modified: check mtime
- [ ] Show loaded files in prompt or status

### Testing
- [ ] Test with very long conversations (>4096 tokens)
- [ ] Test summarization quality
- [ ] Test file reloading on changes
- [ ] Test session cleanup

---

## Milestone 8: Polish & Distribution (Month 2 Week 4)

**Goal:** Make μ ready for real use

### Documentation
- [ ] Write comprehensive README
- [ ] Create USAGE.md with examples
- [ ] Document all commands
- [ ] Add troubleshooting guide
- [ ] Create demo GIFs/videos

### Build System
- [ ] Create Makefile:
  ```makefile
  build:
      go build -o my
  
  install:
      cp my /usr/local/bin/
  
  clean:
      rm my
  ```
- [ ] Add version info to binary
- [ ] Cross-compile for Linux/macOS

### Installation Script
- [ ] Create `install.sh`:
  ```bash
  #!/bin/bash
  # Download binary
  # Move to /usr/local/bin
  # Run my --init
  # Show getting started message
  ```

### Quality Checks
- [ ] Add comments to all public functions
- [ ] Run `go fmt` on all files
- [ ] Run `go vet` and fix warnings
- [ ] Basic tests for core functions

### Performance
- [ ] Profile startup time (should be <50ms)
- [ ] Check memory usage during long sessions
- [ ] Optimize session loading/saving

### Release
- [ ] Tag v0.1.0
- [ ] Create GitHub releases
- [ ] Write release notes
- [ ] Share on relevant forums/subreddits

---

## Milestone 9: Advanced Shell Features (Month 3-4)

**Goal:** Move toward feature parity with basic shells

### Pipes
- [ ] Implement pipe operator: `ls | grep txt`
- [ ] Chain multiple processes
- [ ] Connect stdout to stdin between processes

### Redirection
- [ ] Output redirection: `ls > files.txt`
- [ ] Append: `echo "text" >> file.txt`
- [ ] Input redirection: `sort < input.txt`
- [ ] Error redirection: `2>`, `2>&1`

### Background Jobs
- [ ] Background operator: `long-command &`
- [ ] Job control built-ins:
  - [ ] `jobs` - List background jobs
  - [ ] `fg` - Bring job to foreground
  - [ ] `bg` - Resume job in background
- [ ] Ctrl+Z to suspend current job

### Environment Variables
- [ ] `export VAR=value`
- [ ] Variable expansion: `$VAR`, `${VAR}`
- [ ] Special variables: `$?`, `$!`, `$$`

### Scripting Support
- [ ] Run script files: `my script.sh`
- [ ] Shebang support: `#!/usr/local/bin/my`
- [ ] Conditionals: `if`, `then`, `else`, `fi`
- [ ] Loops: `for`, `while`

---

## Milestone 10: Advanced AI Features (Month 4-6)

**Goal:** Enhanced AI capabilities

### RAG over Project
- [ ] Generate embeddings for project files
- [ ] Semantic search over codebase
- [ ] Answer questions about entire project

### Memory System
- [ ] Extract persistent facts from sessions
- [ ] Store: `memory add "project uses PostgreSQL"`
- [ ] Recall: Load relevant memories into context
- [ ] Manage: `memory list`, `memory remove`

### Multiple Model Support
- [ ] Quick compare: Ask same question to multiple models
- [ ] Model switching: `use qwen` / `use deepseek`
- [ ] Model benchmarking for specific tasks

### Git Hooks
- [ ] `install-hooks` command
- [ ] Pre-commit: Auto-generate message
- [ ] Pre-push: Check for issues
- [ ] Commit-msg: Validate format

---

## Future Ideas (Backlog)

### Nice to Have
- [ ] Tab completion for commands and files
- [ ] Syntax highlighting in prompt
- [ ] Plugin system for custom commands
- [ ] Web interface option (browser-based REPL)
- [ ] Shell themes and customization
- [ ] Alias support
- [ ] Shell functions
- [ ] Multi-line input support
- [ ] Better Unicode support (for μ symbol)

### Advanced Features
- [ ] Distributed mode (connect to remote Ollama)
- [ ] Team features (shared sessions)
- [ ] Integration with IDE (LSP server)
- [ ] Model fine-tuning on project code
- [ ] Automated testing generation
- [ ] Security scanning with AI

### Community
- [ ] Homebrew formula
- [ ] AUR package (Arch Linux)
- [ ] Snap/Flatpak packages
- [ ] Docker image
- [ ] Contributing guidelines
- [ ] Code of conduct

---

## Learning Checkpoints

As you build this, you'll learn:

### Week 1-2: Shell Basics
- [ ] How shells parse and execute commands
- [ ] Fork/exec process model
- [ ] Working with stdin/stdout/stderr
- [ ] Signal handling basics

### Week 3-4: Advanced I/O
- [ ] Process creation and management
- [ ] File descriptors and pipes
- [ ] Terminal control
- [ ] Environment variables

### Month 2-3: System Programming
- [ ] Inter-process communication
- [ ] Job control
- [ ] Session management
- [ ] Terminal modes and raw input

### Month 4+: Advanced Topics
- [ ] Concurrent programming
- [ ] Parser design
- [ ] State machines
- [ ] Performance optimization

---

## Testing Strategy

### Unit Tests
- [ ] Test parser with various inputs
- [ ] Test token counting
- [ ] Test session serialization
- [ ] Test config loading

### Integration Tests
- [ ] Test full command execution flow
- [ ] Test AI query with mocked Ollama
- [ ] Test session persistence

### Manual Testing
- [ ] Daily use as primary shell (dogfooding)
- [ ] Test edge cases as discovered
- [ ] Performance testing with large contexts

---

## Notes

- Start simple, iterate based on what you actually use
- Don't over-engineer early - get it working first
- Test frequently with real workflows
- Document learnings as you go
- Commit often with good messages (use your own tool!)

## Questions to Resolve

- [ ] Use Cobra for CLI or keep it minimal?
- [ ] Readline library or custom input handling?
- [ ] Exact token counting or approximation?
- [ ] Default models if config missing?
- [ ] Session cleanup: manual or automatic?
- [ ] Should `code` mode have separate history?
