# Contributing

[日本語](CONTRIBUTING.ja.md)

Thank you for your interest in contributing to readme-gen!

## Development Setup

### Prerequisites

- Go 1.23+
- [mise](https://mise.jdx.dev/) (optional, but recommended)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/hulk510/readme-gen.git
cd readme-gen

# Install dependencies
go mod download

# Verify setup
go build ./...
go test ./...
```

## Development Commands

### Using mise (recommended)

```bash
mise run build      # Build binary to bin/readme-gen
mise run dev        # Run in development mode
mise run test       # Run all tests
mise run lint       # Run linter
mise run install    # Install to $GOPATH/bin
mise run clean      # Clean build artifacts
```

### Using Go directly

```bash
# Build
go build -o bin/readme-gen ./cmd/readme-gen

# Run
go run ./cmd/readme-gen

# Test
go test ./...

# Test with verbose output
go test -v ./...

# Test specific package
go test ./internal/scanner/...

# Install
go install ./cmd/readme-gen
```

## Project Structure

```
cmd/readme-gen/     # CLI entry point
internal/
├── cmd/            # Cobra command definitions
├── i18n/           # Internationalization
├── marker/         # Marker-based updates
├── scanner/        # Directory scanning
├── template/       # README templates
└── ui/             # Terminal UI styles
```

## Making Changes

### Adding a New Feature

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Write tests first
3. Implement the feature
4. Ensure all tests pass: `go test ./...`
5. Update documentation if needed
6. Submit a pull request

### Adding Translations

1. Update `internal/i18n/i18n.go` with new messages
2. Add both English and Japanese translations
3. Update templates if needed

### Adding a New Template

1. Create template file in `internal/template/templates/`
2. Name format: `{name}.md.tmpl` (English) or `{name}_ja.md.tmpl` (Japanese)
3. Add tests in `internal/template/template_test.go`

## Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions

## Pull Request Process

1. Fork the repository
2. Create your feature branch
3. Commit your changes with clear messages
4. Push to your fork
5. Open a pull request

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new template option
fix: correct structure parsing
docs: update README
test: add scanner tests
refactor: simplify marker logic
```

## Testing

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Test both success and error cases

## Release Process

Releases are automated using [release-please](https://github.com/googleapis/release-please).

1. Merge PRs to `main`
2. release-please creates/updates a release PR
3. Merge the release PR to trigger a new release
4. GoReleaser builds and publishes binaries

## Questions?

Feel free to open an issue for any questions or suggestions!
