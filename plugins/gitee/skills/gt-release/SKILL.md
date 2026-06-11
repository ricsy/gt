---
name: gt-release
description: Manage Gitee releases, webhooks, check runs, git data, and repository statistics through gt. Use when the request is about release-adjacent repository operations rather than issues or pull requests.
---

# GT Release

## Scope

Use this skill for:

- `gt release *`
- `gt webhook *`
- `gt check *`
- `gt gitdata *`
- `gt repostats *`

## Workflow

1. Resolve the repository first.
2. Use the matching high-level `gt` command family.
3. When statistics or check-run data looks inconsistent, confirm backend payload shape with `gt api`.

## Guardrails

- Treat `gitdata`, `check`, and `repostats` as verification-heavy commands.
- Be explicit when the backend returns empty data instead of assuming the CLI is broken.
