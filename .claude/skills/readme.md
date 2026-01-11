---
description: Update README.md structure. Run with /readme or "update the README".
user_invocable: true
---

# README Update Skill

Automatically update README.md when project structure changes.

## Usage

When user runs:
- `/readme` command
- "Update the README" or similar request

## Steps

1. Run `readme-gen check` to detect structure changes
2. If out of sync, run `readme-gen structure --update`
3. Only update structure section (between markers)
4. Modify other sections only when explicitly requested

## Rules

- Only auto-update content between markers (`<!-- readme-gen:structure:start -->` / `<!-- readme-gen:structure:end -->`)
- Preserve comments (`# description`) if present
- Confirm with user for significant changes
