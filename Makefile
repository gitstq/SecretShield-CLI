# Build the binary
build:
	go build -o bin/secretshield .

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f secretshield secretshield.exe

# Run tests
test:
	go test ./... -v

# Install to GOPATH/bin
install:
	go install .

# Cross-compile release binaries
release:
	GOOS=linux GOARCH=amd64 go build -o bin/secretshield-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/secretshield-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o bin/secretshield-windows-amd64.exe .

# Lint with go vet
lint:
	go vet ./...

.PHONY: build clean test install release lint
