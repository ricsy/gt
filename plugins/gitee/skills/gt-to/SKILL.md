---
name: gt-to
description: Open Gitee repository or namespace pages in the default browser through gt. Use when the request is about jumping from the current git worktree, a repo identifier, or a user/org identifier to the corresponding Gitee web page.
---

# GT To

## Scope

Use this skill for:

- `gt to`
- `gt to <owner/repo>`
- `gt to <user-or-org>`
- `gt to <repo> --repo`
- `gt to <repo> -r`

Route to `$gt-repo` instead when the primary request is about repository management rather than browser navigation.

## Workflow

1. Decide whether the user wants a repository page or a namespace homepage.
2. If no target is provided, prefer the current git worktree remote on the active Gitee host.
3. If the target contains `/`, resolve it as `owner/repo`.
4. If the target is a single segment:
   - treat it as a namespace homepage by default
   - treat it as a repository name only when `--repo` or `-r` is explicit
5. When `--repo` is used with a single-segment target, resolve the owner from the authenticated user on the active host.

## Guardrails

- Do not infer a single-segment target as a repository unless the explicit repo mode flag is present.
- Prefer the configured host and current authentication context; do not hardcode user or organization names in examples.
- Keep browser-opening behavior delegated to the shared URL opener instead of reimplementing platform-specific shell commands.
