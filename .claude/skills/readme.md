---
description: Update README.md structure. Run with /readme or "update the README".
user_invocable: true
---

# README Update Skill

Automatically update README.md when project structure changes.

## Usage

- /readme command
- "Update the README" or similar request

## Steps

1. Run `readme-gen check` to detect structure changes
2. If out of sync, run `readme-gen structure --update`
3. After structure update, add/update description comments for new directories
4. Only modify other sections when explicitly requested

## Rules

- Only auto-update content between markers
- Add concise description comments (e.g., `# CLI entry point`)
- Preserve existing comments unless outdated
- Confirm with user for significant changes outside markers
