---
name: adr-new
description: Create a new ADR in docs/decisions using the official template and (optionally) update the documentation.
---

## Required input (ask the user if missing)
- ADR title (short)
- Context and problem statement
- Decision drivers
- Considered options
- Selected option (must be one of the considered options)
- Rationale

## Rules (mandatory)
- Use the template: docs/decisions/_templates/adr-template.md
- Create the file at: docs/decisions/NNNN-<slug>.md
- NNNN = next available number (4 digits, sequential)
- Slug = kebab-case
- Set:
  - Status: Proposed
  - Date: current date and time (YYYY-MM-DD HH:MM:SS)
  - Acceptance date: -
- DO NOT modify the template.

## Template cleanup (MANDATORY)
- Do not leave any template placeholder text in the file, such as lines containing:
  `[ad esempio`, `[etc.]`, `[opzione`, `<!-- opzionale -->`.
- If there is no real content for an optional section, remove the entire section (including its heading) instead of leaving placeholders.

## Output
1) New ADR file created and populated.
2) If requested by the user, add a reference in docs/documentation/0002-architettura.md under a "Decisions" section (create it if missing), as a list of links.
3) Show the diff.
