---
description: Create a classified GitHub backlog issue through an interactive questionnaire.
---

# Create GitHub Backlog Issue

Create a new GitHub Issue that follows the repository workflow, milestone numbering, label taxonomy, and body format.

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Turn a rough task description into a fully classified GitHub Issue by asking the user only the minimum questions needed to place it correctly.

## Required Behavior

1. Treat `$ARGUMENTS` as the initial task context.
2. If the input is empty, ask for a one-paragraph description of the task before doing anything else.
3. Read `docs/workflow-github-issues-y-specify.md` and `.github/ISSUE_TEMPLATE/backlog-item.yml` before classifying.
4. Inspect current milestones and labels with `gh` so you use the live repository taxonomy.
5. Ask the user a short interactive questionnaire that covers, at minimum:
   - milestone
   - area
   - work type
   - priority
   - dependency summary
   - objective
   - expected result
   - blocked vs ready-for-spec
6. Prefer multiple-choice questions when possible. Use concise labels and recommend the most likely option first when you can infer one safely.
7. Infer `stream:*` from the chosen milestone unless the user explicitly overrides it.
8. Compute the next ordered ID for the chosen milestone prefix by inspecting existing issue titles.
9. Build the issue title as `[PREFIX-##] Short title`.
10. Build the issue body using the repository template structure, even when creating through `gh` directly.
11. Show the user a short preview containing title, milestone, labels, dependencies, and status.
12. Create the issue with `gh issue create`.
13. Report the created issue URL and the chosen classification.

## Milestone Prefix Map

- `Run MVP Foundation` -> `RMF`
- `Zone Progression` -> `ZP`
- `Fish Data-Driven Expansion` -> `FDE`
- `Build and Services` -> `BS`
- `Economy and Meta Boundary` -> `EMB`
- `Architecture and Delivery` -> `AD`
- `Graphical Client Prototype` -> `GCP`

## Label Rules

Every created issue must include:

- exactly one `stream:*`
- exactly one `area:*`
- exactly one `type:*`
- exactly one `priority:*`
- exactly one `status:*`

Default status is `status:ready-for-spec` unless the user says the item is blocked.

## Questioning Style

- Ask only what is still ambiguous after reading the initial context.
- Group the classification into as few questions as practical.
- If the objective or expected result is weak, ask one follow-up to sharpen them.
- Do not ask the user to invent labels or prefixes manually.

## Implementation Notes

- Use `gh issue list --state all --json title,milestone` to inspect numbering in the selected milestone.
- You may use `python3 scripts/github_issue_next_id.py <PREFIX>` to compute the next ID if available.
- Use a temporary file or heredoc for the issue body so formatting is preserved.
- If an issue with the same computed title already exists, warn the user and stop instead of creating a duplicate.

## Output

When done, return:

- issue title
- milestone
- labels
- issue URL
- next recommended step, usually `/speckit.specify`
