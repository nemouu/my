# μ (My) Shell

A minimal shell with integrated local LLM capabilities via Ollama. Unlike bolt-on AI tools, μ treats AI as a first-class shell feature with automatic context awareness and seamless workflow integration.

**Name Origin:** μ (my) — the Greek letter, representing "micro" — a minimal, focused shell that does one thing well: blend traditional shell commands with AI assistance.

> **Status:** Design phase — project skeleton in place, actively being developed as a learning project.

## Core Ideas

- AI features feel native, not bolted on — `ask`, `commit`, and `code` are built-in commands just like `cd` or `exit`
- Standard commands just work via fork+exec (gcc, git, flutter, etc.)
- Per-directory context automatically maintained across sessions
- Privacy-first: everything runs locally through Ollama
- A long-term systems programming learning project built from the ground up

## Usage

**Standard shell commands:**
```bash
μ ~/project> ls -la
μ ~/project> gcc main.c -o main
μ ~/project> git status
μ ~/project> cd src
```

**Ask — quick AI queries with automatic context:**
```bash
μ ~/project> ask what does the divide() function do in main.c?
# Automatically loads main.c, queries local LLM

μ ~/project> ask compare auth.go and middleware.go
# Both files auto-loaded into context
```

**Commit — AI-powered git commit messages:**
```bash
μ ~/project> git add .
μ ~/project> commit
# Generates conventional commit message from staged diff
# [C]ommit [E]dit [R]egenerate [A]bort
```

**Code — interactive mode for code queries and generation:**
```bash
μ ~/project> code
code> explain divide() in utils.c
code> how are utils.c and main.c connected?
code> create helper.c with file I/O functions
code> exit
```

## Requirements

- Go 1.21+
- Ollama installed locally with models pulled (e.g. qwen2.5:7b, qwen2.5-coder:7b)
- 16GB RAM minimum (32GB recommended)

## Getting Started

```bash
git clone https://github.com/nemouu/my
cd my
go build -o my
./my --init   # creates ~/.config/my/config.yaml with defaults
./my           # start the shell
```

## Project Structure

```
my/
├── main.go              # Entry point
├── shell/
│   ├── repl.go          # Main shell loop
│   ├── parser.go        # Input parsing
│   ├── executor.go      # Command execution
│   └── builtins.go      # Built-in commands (cd, exit, ask, commit, code)
├── ai/
│   ├── client.go        # Ollama API wrapper
│   ├── modes.go         # ask/code/commit modes
│   └── session.go       # Per-directory session & context management
├── git/
│   └── integration.go   # Git helpers for commit command
├── config/
│   └── config.go        # YAML config loading
└── go.mod
```

See [TODO.md](TODO.md) for the full development roadmap and detailed implementation plan.

## License

MIT — see [LICENSE](LICENSE) for details.
