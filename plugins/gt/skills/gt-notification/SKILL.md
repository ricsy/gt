---
name: gt-notification
description: Manage Gitee notifications, messages, gists, user, org, and activity queries through gt. Use when the request is about inbox state, user-facing account activity, or Gist workflows.
---

# GT Notification

## Scope

Use this skill for:

- `gt notification *`
- `gt gist *`
- `gt user *`
- `gt org *`
- `gt activity *`
- `gt email list`

## Workflow

1. Decide whether the request is about notifications, gists, account profile, org membership, or activity streams.
2. Prefer the domain-specific `gt` command family directly.
3. When a read path requires existing external state, say that explicitly instead of implying a CLI defect.

## Guardrails

- Distinguish platform visibility limits from CLI bugs, especially for `gist` and notification read flows.
- Use `gt api` only when verifying whether a missing item is absent from the backend or just hidden by CLI formatting.
