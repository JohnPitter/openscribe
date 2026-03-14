.PHONY: build test test-coverage lint clean e2e

## build: Build the project
build:
	go build ./...

## test: Run all tests
test:
	go test ./... -v

## test-coverage: Run tests with coverage report
test-coverage:
	go test ./... -coverprofile=coverage.txt -covermode=atomic -v
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

## e2e: Run end-to-end tests
e2e:
	go test ./e2e/... -v -count=1

## lint: Run linter
lint:
	go vet ./...

## clean: Clean build artifacts
clean:
	rm -f coverage.txt coverage.html
	rm -rf testdata/output/
	go clean ./...

## fmt: Format code
fmt:
	gofmt -s -w .

## help: Show this help
help:
	@echo "Available commands:"
	@grep -E '^## ' Makefile | sed 's/## /  /'
