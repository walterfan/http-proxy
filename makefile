# Name of the binary output file
BINARY=http-proxy

# Main package path (adjust if needed)
MAIN_PKG=github.com/walterfan/http-proxy

# Flags for go build
GOBUILD=go build -o ${BINARY} ${MAIN_PKG}

# Default target
all: build

# Build the application
build:
	@echo "Building ${BINARY}..."
	${GOBUILD}

# Run the application
run:
	@echo "Running ${BINARY}..."
	go run ${MAIN_PKG}

# Run tests
test:
	@echo "Running tests..."
	go test ${MAIN_PKG}/...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f ${BINARY}

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ${MAIN_PKG}/...

# Vet Go code
vet:
	@echo "Vetting Go code..."
	go vet ${MAIN_PKG}/...

# Lint Go code (requires golangci-lint installed)
lint:
	@echo "Linting Go code..."
	golangci-lint run

.PHONY: all build run test clean fmt vet lint