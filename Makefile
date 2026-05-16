.PHONY: build test vet clean install hook-bash hook-zsh

BINARY := xc
INSTALL_DIR := $(HOME)/.local/bin

build:
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
