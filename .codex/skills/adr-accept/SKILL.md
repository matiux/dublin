---
name: adr-accept
description: Set a specific ADR to accepted and update the acceptance date, including support for Structurizr.
---

## Required input
- ADR file path (e.g. docs/decisions/0002-something.md)

## Rules
1. If a line exists:
   - `Status: ...`
   - update it to `Status: Accepted`
   - update (or add if missing) the line:
      - `Acceptance date: YYYY-MM-DD HH:MM:SS` (current date and time)

2. If the status is already `Accepted`:
   - make no changes

3. Do not change anything else in the file.

4. Always show the diff.

## Notes
- Do not rename files.
- Do not modify templates or other documents.
