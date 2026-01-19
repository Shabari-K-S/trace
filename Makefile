BINARY_NAME=trace
GO=go

.PHONY: all build install clean test

all: build

build:
	$(GO) build -o $(BINARY_NAME) cmd/trace/main.go

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -rf .trace

test:
	$(GO) test ./...
