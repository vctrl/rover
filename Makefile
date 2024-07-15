.PHONY: all build run console file interactive docker-build docker-run docker-console docker-file docker-interactive test lint install-lint install-goimports format

# Default target
all: build

# Build the Go binary
build:
	go build -o rover ./cmd/rover

# Run the Go binary with no flags (prompts for mode)
run: build
	./rover

# Run the Go binary in console mode
console: build
	./rover --mode=console

# Run the Go binary in file mode with a specified file
file: build
	./rover --mode=file --file=$(FILE)

# Run the Go binary in interactive mode
interactive: build
	./rover --mode=interactive

# Build the Docker image
docker-build:
	docker-compose build

# Run the Docker container (prompts for mode)
docker-run:
	docker-compose run rover ./rover

# Run the Docker container in console mode
docker-console:
	docker-compose run rover ./rover --mode=console

# Run the Docker container in file mode with a specified file
docker-file:
	docker-compose run rover ./rover --mode=file --file=$(FILE)

# Run the Docker container in interactive mode
docker-interactive:
	docker-compose run rover ./rover --mode=interactive

# Run tests
test:
	go test ./...

# Check and install golangci-lint if not present
install-lint:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint не найден. Установка..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Check and install goimports if not present
install-goimports:
	@if ! [ -x "$$(command -v goimports)" ]; then \
		echo "goimports не найден. Установка..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi

# Format code
format: install-goimports
	goimports -w .

# Run golangci-lint
lint: install-lint format
	golangci-lint run
