# Contributing to chatwork-go

First off, thank you for considering contributing to chatwork-go! It's people like you that make chatwork-go such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible.
* **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples.
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**
* **Include the version of chatwork-go you're using** and the version of Go.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
* **Provide specific examples to demonstrate the steps**.
* **Describe the current behavior** and **explain which behavior you expected to see instead** and why.
* **Explain why this enhancement would be useful** to most chatwork-go users.

### Pull Requests

Please follow these steps to have your contribution considered by the maintainers:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Development Setup

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/chatwork-go.git
   cd chatwork-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a branch for your feature or fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Guidelines

### Code Style

* Follow the standard Go code style guidelines
* Run `gofmt` and `goimports` on your code
* Use `golangci-lint` to check for common issues:
  ```bash
  golangci-lint run
  ```

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

### Testing

* Write unit tests for new functionality
* Ensure all tests pass before submitting PR:
  ```bash
  go test -v ./...
  ```
* Add integration tests for new API endpoints (use build tags to separate from unit tests)
* Aim for good test coverage, but prioritize meaningful tests over coverage percentage

### Documentation

* Add GoDoc comments to all exported types, functions, and methods
* Update the README.md if you're adding new features or changing behavior
* Include examples in your documentation when appropriate
* Follow the Go documentation conventions

### API Design

* Follow RESTful conventions where applicable
* Keep the API intuitive and consistent with existing patterns
* Provide convenience methods where they add value
* Return appropriate error types with meaningful messages

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable.
2. Update the CHANGELOG.md with a note describing your changes.
3. The PR will be merged once you have the sign-off of at least one maintainer.

## Release Process

Releases are managed by the maintainers. The process is:

1. Update version numbers
2. Update CHANGELOG.md
3. Create a git tag
4. Push the tag to trigger the release workflow

## Questions?

Feel free to open an issue with your question or contact the maintainers directly.

Thank you for contributing! 🎉