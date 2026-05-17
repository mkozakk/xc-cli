.PHONY: build test vet clean install hook-bash hook-zsh check-deps

BINARY := xc
INSTALL_DIR := $(HOME)/.local/bin

check-deps:
	@missing=""; \
	if [ -n "$$WAYLAND_DISPLAY" ]; then \
		command -v wl-copy >/dev/null 2>&1 || missing="$$missing wl-clipboard"; \
	fi; \
	command -v xclip >/dev/null 2>&1 || command -v xsel >/dev/null 2>&1 || \
		{ [ -z "$$WAYLAND_DISPLAY" ] && missing="$$missing xclip"; }; \
	if [ -n "$$missing" ]; then \
		echo "Warning: missing clipboard tools:$$missing"; \
		echo "xc will not be able to copy to clipboard until one is installed."; \
	fi

build: check-deps
	go build -o $(BINARY) .

test:
	go test -race ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

install: build
	install -m755 $(BINARY) $(INSTALL_DIR)/$(BINARY)

hook-bash: build
	./$(BINARY) init bash

hook-zsh: build
	./$(BINARY) init zsh

.DEFAULT_GOAL := build
