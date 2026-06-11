---
name: gt-issue
description: Manage Gitee issues, labels, and milestones through gt. Use when the request is about issue creation, triage, comments, labels, or milestone workflows.
---

# GT Issue

## Scope

Use this skill for:

- `gt issue *`
- `gt label *`
- `gt milestone *`

## Workflow

1. Resolve the target repository before acting.
2. Use `gt issue` for issue state and comments.
3. Use `gt label` for repository, issue, or project label operations.
4. Use `gt milestone` for milestone lifecycle operations.
5. If the request is about missing or surprising fields, verify with `gt api` only when needed.

## Guardrails

- Keep repository labels, issue labels, and project labels distinct in both language and command choice.
- Call out platform-required fields when they differ from older CLI behavior.
