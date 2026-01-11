# README Update Skill

Update README.md when project structure or functionality changes significantly.

## Trigger

- Directory structure changed
- New features added
- User requests "update README"

## Steps

1. Run `readme-gen check` to detect structure changes
2. If out of sync, run `readme-gen structure --update`
3. Review and update other sections (Usage, Description) based on code changes
4. Keep updates concise and accurate

## Rules

- Don't modify content outside markers unless explicitly needed
- Keep README concise - avoid over-documentation
- Ask user if unsure about significant changes
