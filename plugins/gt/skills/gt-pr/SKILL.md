---
name: gt-pr
description: Manage Gitee pull requests through gt, including listing, viewing, creating, commenting, closing, state changes, and merge operations. Use when the request is primarily about pull requests.
---

# GT PR

## Scope

Use this skill for:

- `gt pr list/view/create/comment/state/close/merge`

## Workflow

1. Resolve the repository and branch context first.
2. Prefer `gt pr create` and `gt pr merge` over raw API calls.
3. When merge fails, separate platform policy failures from CLI defects.
4. If needed, use `gt api` only to confirm backend semantics.

## Guardrails

- Do not describe a repository policy failure as a broken merge command.
- Treat comments and state transitions as part of the PR workflow, not the issue workflow.
