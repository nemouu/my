# Makefile for mu shell tests
# Run from the test/ directory

SHELL_BIN := ../my
BIN_DIR := bin

# All test binaries
TESTS := $(BIN_DIR)/test_jobs $(BIN_DIR)/test_signals

.PHONY: all test clean

# Build all test binaries
all: $(BIN_DIR) $(TESTS)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_DIR)/test_jobs: jobs/test_jobs.go
	go build -o $@ ./jobs/

$(BIN_DIR)/test_signals: signals/test_signals.go
	go build -o $@ ./signals/

# Run all tests (requires shell binary to be built first)
test: all
	@echo "Running all tests against $(SHELL_BIN)"
	@echo ""
	@$(BIN_DIR)/test_jobs $(SHELL_BIN)
	@echo ""
	@$(BIN_DIR)/test_signals $(SHELL_BIN)

# Remove compiled binaries
clean:
	rm -rf $(BIN_DIR)