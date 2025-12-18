# μ (My) Shell

A minimal shell with integrated local LLM capabilities via Ollama. Unlike bolt-on AI tools, μ treats AI as a first-class shell feature with automatic context awareness and seamless workflow integration.

**Name Origin:** μ (my) - the Greek letter, representing "micro" - a minimal, focused shell that does one thing well: blend traditional shell commands with AI assistance.

## Philosophy

Build a real shell from the ground up as a long-term systems programming learning project, while delivering immediate practical value through AI integration. Start simple, grow organically.

**Core Principles:**
- AI features feel native, not bolted on
- Standard commands just work (gcc, git, flutter, etc.)
- Per-directory context automatically maintained
- Privacy-first: everything runs locally
- Learn shell internals by building one

## What Makes μ Different

**Existing AI shell tools:**
- Shell plugins or CLI wrappers (not actual shells)
- Require hotkeys and explicit invocation
- Context switching between shell and AI
- Usually cloud-dependent

**μ approach:**
- Actual shell that you run as your primary environment
- AI commands are built-in (like `cd` or `exit`)
- Automatic context from current directory and session
- Local-first by design (Ollama)

## Core Features

### Standard Shell
```bash
μ ~/project> ls -la
μ ~/project> gcc main.c -o main
μ ~/project> git status
μ ~/project> cd src
μ ~/project/src> flutter build
```

All standard commands work via fork+exec - just like bash/zsh.

### Built-in AI Commands

**Quick Queries:**
```bash
μ ~/project> ask what does the divide() function do in main.c?
[Automatically loads main.c, queries local LLM]

μ ~/project> ask how to find files modified today
[AI suggests command, you confirm execution]
```

**Git Integration:**
```bash
μ ~/project> git add .
μ ~/project> commit
Generated commit message:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
feat: Add authentication middleware

- Implement JWT token validation
- Add login/logout endpoints
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[C]ommit [E]dit [R]egenerate [A]bort >
```

**Code Mode:**
```bash
μ ~/project> code
code> explain divide() in utils.c
[AI explains function with context]
code> how are utils.c and main.c connected?
[AI analyzes multiple files]
code> create helper.c with file I/O functions
[AI generates code, shows preview, asks to save]
code> exit
μ ~/project>
```

### Automatic Context Management

**Per-directory sessions:**
- Context automatically saved/loaded based on current directory
- No manual session management required
- Ask follow-up questions naturally - previous context remembered
- Switch directories = switch contexts automatically

**Smart file loading:**
```bash
μ ~/project> ask compare auth.go and middleware.go
[Both files auto-loaded into context]
```

**Token management:**
- Warning at 80% context capacity
- Auto-truncation of old messages when full
- File contents preserved (loaded files never truncated)

## How It Works

### Command Execution Model

```
User input
    ↓
Built-in? (cd, exit, ask, commit, code)
    ↓ YES → Handle in Go
    ↓ NO
    ↓
External command (gcc, git, ls, etc.)
    ↓ → fork() + exec() → Just works!
```

**Built-ins:** `cd`, `exit`, `pwd`, `ask`, `commit`, `code`, `history`, `session`

**Everything else:** Executed as external process (standard Unix model)

### Session Storage

```
~/.config/my/
├── config.yaml
├── sessions/
│   ├── abc123.json    # ~/project1
│   ├── def456.json    # ~/project2
│   └── ...
└── history
```

Each session contains:
- Conversation history
- Loaded files and their contents
- Token usage tracking
- Timestamp

### Model Selection

Different models optimized for different tasks:
- **ask**: General language model (qwen2.5:7b)
- **code**: Code-specialized model (qwen2.5-coder:7b)
- **commit**: Concise commit generation (deepseek-coder:6.7b)

## Configuration

YAML config at `~/.config/my/config.yaml`:

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

## Hardware Requirements

- 16GB RAM minimum (32GB recommended)
- Local Ollama installation
- Works completely offline after initial model download

## Technology Stack

**Language:** Go
- Easy process management (exec package)
- Single binary distribution
- Good learning path from C
- Fast startup, compiled performance

**Dependencies:**
- Cobra (CLI framework - maybe, TBD)
- Ollama Go API client
- YAML parsing (gopkg.in/yaml.v3)

**Architecture:**
```
my/
├── main.go              # Entry point
├── shell/
│   ├── repl.go         # Main shell loop
│   ├── parser.go       # Input parsing
│   ├── executor.go     # Command execution
│   └── builtins.go     # Built-in commands
├── ai/
│   ├── client.go       # Ollama wrapper
│   ├── modes.go        # ask/code/commit modes
│   └── session.go      # Context management
├── git/
│   └── integration.go  # Git helpers
├── config/
│   └── config.go       # Config loading
└── go.mod
```

## Project Goals

### Short Term (Weeks 1-4)
1. Working shell that executes any command
2. `ask` command with auto-context
3. `commit` command with AI generation
4. Basic `code` mode

### Medium Term (Months 2-3)
1. Robust session management
2. Context optimization (summarization)
3. Shell history and completion
4. Git hooks generation

### Long Term (Months 4-12+)
1. Pipes and redirection
2. Job control (background processes, Ctrl+Z)
3. Shell scripting support
4. Advanced features (aliases, functions)

## Learning Journey

This project teaches:
- **Shell internals**: How shells actually work (fork/exec model)
- **Process management**: Creating, controlling, waiting for processes
- **Systems programming**: File descriptors, signals, environment
- **Go language**: Building real software in Go
- **LLM integration**: Working with local AI APIs
- **Unix philosophy**: Small tools that do one thing well

## Non-Goals

- Not trying to replace bash/zsh for power users (initially)
- Not a full IDE or development environment
- Not trying to compete with cloud AI quality
- Not implementing bash/zsh compatibility (fresh start)

## Use Cases

**Perfect for:**
- Learning codebases (explain functions, trace connections)
- Quick git commits (stop writing "fix bug")
- Exploring new projects (ask about architecture)
- Solo development with AI pair programming
- Systems programming practice project

**Not ideal for:**
- Complex shell scripts (use bash/zsh for now)
- Production deployment automation
- Heavy piping workflows (will improve over time)

## Installation (Future)

```bash
# Install from source
git clone https://github.com/yourusername/my
cd my
go build -o my
sudo mv my /usr/local/bin/

# Initialize config
my --init

# Start using
my
```

## Inspiration

**Concept:** Like having Claude Code's capabilities built into your shell, but using local models and maintaining privacy.

**Technical:** Standard Unix shell architecture (similar to dash/sh) with AI as native feature rather than plugin.

## Why "μ" (My)?

- **Minimal**: Small, focused shell - not trying to do everything
- **Micro**: Fits the "small but powerful" philosophy
- **Memorable**: Single character, easy to type
- **Greek letter**: Programming tradition (λ calculus, μ-recursive functions)
- **Unique**: Not competing with bash/zsh namespace

---

*Status: Design phase - actively being developed as a learning project*
