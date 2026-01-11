# readme-gen

[日本語](README.ja.md)

A CLI tool for managing README.md with automatic structure sync.

## Features

- 📝 Generate README from templates (oss / personal / team)
- 🔄 Marker-based directory structure auto-sync
- 🤖 AI description generation with Claude Code
- ✅ CI-friendly check command
- 🌐 English / Japanese support

## Structure

<!-- readme-gen:structure:start -->
```
├── .claude/           # Claude Code skills
│   └── skills/
├── .github/           # GitHub Actions
│   └── workflows/
├── cmd/               # CLI entry point
│   └── readme-gen/
├── extras/            # Distributable skills for users
│   └── skills/
└── internal/          # Internal packages
    ├── cmd/           # Cobra command definitions
    ├── i18n/          # Internationalization (en/ja)
    ├── marker/        # Marker update processing
    ├── scanner/       # Directory scanning
    ├── template/      # Template processing
    │   └── templates/
    └── ui/            # Charm UI styles
```
<!-- readme-gen:structure:end -->

## Installation

```bash
# Go
go install github.com/hulk510/readme-gen@latest

# or curl
curl -fsSL https://raw.githubusercontent.com/hulk510/readme-gen/main/install.sh | bash
```

## Usage

### Initialize

```bash
# Interactive mode
readme-gen init

# Specify template
readme-gen init --template oss

# Non-interactive mode (all defaults)
readme-gen init --yes

# With AI generation
readme-gen init --yes --with-ai
```

### Update Structure

```bash
# Show current structure
readme-gen structure

# Update structure in README.md
readme-gen structure --update
```

### Check Diff

```bash
# Check if structure is up to date (for CI)
readme-gen check
# → exits with code 1 if out of sync
```

## Command Options

### `readme-gen init`

| Option | Description |
|--------|-------------|
| `-t, --template` | Template selection (oss, personal, team) |
| `-y, --yes` | Non-interactive mode |
| `--with-skills` | Add Claude Code skills |
| `--with-ai` | Generate descriptions with AI |
| `--no-skills` | Skip adding skills |
| `--no-ai` | Skip AI generation |
| `--lang` | Language (en, ja) |

### `readme-gen structure`

| Option | Description |
|--------|-------------|
| `--update` | Update structure in README.md |

### `readme-gen check`

| Option | Description |
|--------|-------------|
| `--lang` | Language (en, ja) |

## Claude Code Integration

When you add Claude Code skills with `readme-gen init`, `.claude/skills/readme-update.md` is created.

When you open the project with Claude Code, it will suggest README updates when structure changes.

### AI Description Generation

With the `--with-ai` option or interactive selection, Claude Code automatically generates descriptions for each directory.

```bash
readme-gen init --yes --with-ai
```

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

### Quick Start

```bash
# Clone
git clone https://github.com/hulk510/readme-gen.git
cd readme-gen

# Install dependencies
go mod download

# Build
mise run build
# or
go build -o bin/readme-gen ./cmd/readme-gen

# Run tests
mise run test
# or
go test ./...

# Run locally
mise run dev
# or
go run ./cmd/readme-gen
```

## License

MIT License - see [LICENSE](LICENSE) for details.
