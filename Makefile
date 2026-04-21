BINARY      := gittale
INSTALL     := /usr/local/bin/$(BINARY)
CONFIG_FILE ?= config/config.yaml

# configgen is a Go tool in this repo that reads the YAML and emits ldflags.
LDFLAGS := -ldflags "$(shell go run ./cmd/configgen $(CONFIG_FILE))"

.PHONY: build install uninstall

# Build with all config values baked in from CONFIG_FILE.
build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/cli

# Build and install. The binary is fully self-contained — no config file needed at runtime.
install: build
	sudo install -m 0755 $(BINARY) $(INSTALL)

uninstall:
	sudo rm -f $(INSTALL)
