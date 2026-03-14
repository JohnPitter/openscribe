# Contributing to OpenScribe

Thank you for your interest in contributing to OpenScribe!

## Development Setup

1. Fork and clone the repository
2. Ensure you have Go 1.22+ installed
3. Run `make test` to verify your setup

## Code Guidelines

- Follow standard Go conventions (`gofmt`, `go vet`)
- Write tests for new functionality
- Keep the library free of external dependencies
- Use the `common` package for shared types

## Pull Request Process

1. Create a feature branch from `main`
2. Write tests first (TDD preferred)
3. Ensure `make test` passes
4. Update CHANGELOG.md if applicable
5. Submit your PR with a clear description

## Reporting Issues

Use GitHub Issues to report bugs or request features. Include:
- Go version
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
