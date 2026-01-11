# CLAUDE.md

readme-gen - A CLI tool for managing README.md with structure auto-sync.

## Project Structure

```
cmd/readme-gen/     # Entry point
internal/
├── cmd/            # CLI commands (cobra)
├── scanner/        # Directory structure scanning
├── template/       # README template processing
├── marker/         # Marker-based partial updates
└── ui/             # Charm-based UI styles
templates/          # Embedded README templates
extras/skills/      # Claude Code skills for users
```

## Commands

```bash
mise run build      # Build binary
mise run dev        # Run in dev mode
mise run test       # Run tests
mise run install    # Install locally
```

## Architecture

- Uses `spf13/cobra` for CLI
- Uses `charmbracelet/huh` for interactive prompts
- Uses `charmbracelet/lipgloss` for styling
- Templates embedded via `embed` package

## Key Conventions

- Markers format: `<!-- readme-gen:structure:start -->` / `<!-- readme-gen:structure:end -->`
- Detect Go projects via `go.mod`, TS via `package.json`
- Default excludes: .git, node_modules, vendor, dist, build
