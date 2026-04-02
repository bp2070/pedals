.PHONY: build run test clean verify

# Build the application
build:
	go build -o pedals ./cmd/pedals

# Run the application
run: build
	./pedals

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f pedals pedals.exe pedals-* verify
	rm -rf dist/

# Verify compilation
verify:
	go run verify.go

# Build for multiple platforms
build-all: build-linux build-macos build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/pedals-linux-amd64 ./cmd/pedals

build-macos:
	GOOS=darwin GOARCH=arm64 go build -o dist/pedals-darwin-arm64 ./cmd/pedals
	GOOS=darwin GOARCH=amd64 go build -o dist/pedals-darwin-amd64 ./cmd/pedals

build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/pedals-windows-amd64.exe ./cmd/pedals

# Install dependencies
deps:
	go mod tidy
	go mod download

# Development server (example)
dev:
	@echo "Starting development..."
	@echo "Run './pedals' to start the TUI"