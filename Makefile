.PHONY: all test clean build install

PROG = gameroom
GOFLAGS ?= $(GOFLAGS:)

all: install test

build:
	@go build $(GOFLAGS) ./...

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

kill:
	-@killall -9 $(PROG) 2>/dev/null || true

run:
	@make kill
	@make build; (if [ "$$?" -eq 0 ]; then (./${PROG}); fi)