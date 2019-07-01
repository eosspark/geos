# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

GOBIN = $(shell pwd)/build/bin
GO ?= latest

eosgo:
	build/env.sh go run build/ci.go install ./programs/eosgo
	@echo "Done building."
	@echo "Run \"$(GOBIN)/eosgo\" to launch eosgo."

all:
	build/env.sh go run build/ci.go install

test:all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

