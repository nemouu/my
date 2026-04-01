// Package test contains integration tests for the mu shell.
//
// # Structure
//
// Tests are organized by feature area:
//
//	test/
//	├── doc.go         # this file
//	├── jobs/          # job control tests (fg, bg, jobs, Ctrl+Z)
//	│   └── test_jobs.go
//	├── signals/       # signal handling tests (SIGINT, SIGTSTP, SIGCHLD)
//	│   └── test_signals.go
//	└── bin/           # compiled test binaries (not committed to git)
//
// # Building
//
// Use the Makefile to build all tests at once:
//
//	cd test && make
//
// Or build individually:
//
//	go build -o test/bin/test_jobs ./test/jobs/
//	go build -o test/bin/test_signals ./test/signals/
//
// # Running
//
// Tests require the mu shell binary to be built first:
//
//	go build -o my .
//
// Then run from the project root:
//
//	./test/bin/test_jobs ./my
//	./test/bin/test_signals ./my
//
// Or run all tests with:
//
//	cd test && make test
//
// # Adding New Tests
//
// To add a new test suite:
//  1. Create a new folder under test/ named after the feature (e.g. test/pipes/)
//  2. Add a single Go file with package main and a main() function
//  3. Add a build target to the Makefile
//  4. Tests should spawn the shell with the -p=false flag to suppress the prompt
//  5. Use the spawnShell() and readUntil() helpers as a pattern
package test
