# Contributing to comai

Thank you for your interest in contributing to **comai** — Here are few guidelines to get started.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Ollama](https://ollama.com) running locally with at least one model pulled (e.g. `ollama pull qwen3:8b`)
- Git

## Getting started

```bash
git clone https://github.com/JulienQNN/comai.git
cd comai
go mod download
go build -o comai .
```

## Making changes

1. Fork the repository and create a branch from `release/v0.1.0`
2. Keep changes focused — one concern per commit
3. Run `go build ./...` and `go vet ./...` before submitting

## Ideal issues

New features: issues proposing new features or enhancements to existing ones.
Bugs: issues where comai stops working or where it results in an unexpected error.
Refactoring: issues that improve code structure, readability, or maintainability without changing functionality.
Security: issues that could lead to a security vulnerability. As mentioned in SECURITY.md, please do not disclose security vulnerabilities publicly. Report them privately to the maintainers.
Performance: issues to make comai faster.
Documentation: issues that improve or add documentation, including examples and usage guides.

## Submitting a pull request

- Open a PR against `main`
- Describe **what** changed and **why**
- Link any related issue if applicable

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
