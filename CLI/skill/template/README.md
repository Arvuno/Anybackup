# Foundation Skill Template

This directory provides reusable templates for Foundation `foundation-cli` skills.

## Contents

- `domain-entry-rules.md`
  - Reusable rules for domain entry skills
  - Use it when standardizing `SKILL.md -> commands.md -> commands/*.md`
- `domain-entry-skill-template.md`
  - Reusable entry-skill template
  - Use it when creating or refactoring a domain `SKILL.md`
- `domains/domain-template.md`
  - Detailed reference template
  - Use it for intent-to-command mappings, write-operation notes, array parameter conventions, and fallback guidance

## Recommended Usage

1. When a domain skill should behave as a lightweight entry page, follow `domain-entry-rules.md` and start from `domain-entry-skill-template.md`.
2. For the detailed command reference format, copy `domains/domain-template.md` and fill in the command-level details.
3. After editing skill docs, run `python3 skill/scripts/validate_skill_docs.py` to catch BOM, broken links, and structure drift.

## Fixed Conventions

1. Prefer `foundation-cli <domain> <resource> <action>` standard commands.
2. Use `foundation-cli api` only when no business command exists.
3. Entry skills should follow `SKILL.md -> references/commands.md -> references/commands/*.md`.
4. Keep `SKILL.md` as an entry page; push examples, request-body details, and troubleshooting into command docs when needed.
