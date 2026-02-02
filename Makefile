.PHONY: build test acceptance clean install

GO := CGO_ENABLED=0 go
BINARY := claudit

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./internal/...

acceptance: build
	$(GO) test ./tests/acceptance/... -v

all: test acceptance

clean:
	rm -f $(BINARY)

install: build
	cp $(BINARY) $(GOPATH)/bin/

# Development helpers
.PHONY: fmt lint

fmt:
	go fmt ./...

lint:
	go vet ./...
