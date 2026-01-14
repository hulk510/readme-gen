# readme-gen

[日本語](README.ja.md)

A CLI tool for managing README.md with automatic structure sync.

## Features

- Generate README from templates (oss / general)
- Marker-based directory structure auto-sync
- AI description generation with Claude Code
- CI-friendly check command
- English / Japanese support

## Use Cases

**Perfect for:**
- Small to medium projects that need a quick README
- Projects created with `create-next-app`, `go mod init`, etc. where the initial README is outdated
- Personal projects where README maintenance is neglected
- First draft generation when you just want something to start with

**Not designed for:**
- Large monorepos with complex structures
- Projects requiring detailed documentation (use dedicated docs tools)
- Generating perfect, production-ready documentation

readme-gen focuses on **structure sync** and **initial scaffolding**, not comprehensive documentation generation.

## Structure

<!-- readme-gen:structure:start -->
```
├── .claude/
│   └── skills/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   └── workflows/
├── cmd/
│   └── readme-gen/
└── internal/
    ├── cmd/
    ├── config/
    ├── i18n/
    ├── marker/
    ├── scanner/
    ├── template/
    │   └── templates/
    └── ui/
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
| `-t, --template` | Template selection (oss, general) |
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
