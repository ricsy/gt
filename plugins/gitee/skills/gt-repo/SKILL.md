---
name: gt-repo
description: Manage Gitee repositories through gt, including repository inspection, creation, clone, branches, collaborators, and forks. Use when the request is primarily about repository state or repository-level operations.
---

# GT Repo

## Scope

Use this skill for:

- `gt repo list/view/create/delete/clone`
- `gt repo branch *`
- `gt repo collaborator *`
- `gt repo fork *`

Route to `$gt-mode` instead when the primary request is about repository scope switching, org-vs-personal execution context, or env-file driven temporary setup.

## Workflow

1. Resolve the repository identifier first.
2. Prefer the highest-level `gt repo ...` command that matches the user request.
3. If the user is asking about repository capabilities and behavior, compare with `gt api` only when needed.
4. When visibility or permission semantics are involved, call out platform limits explicitly.

## Guardrails

- Do not jump to `gt api` if `gt repo` already covers the operation.
- Treat collaborator, fork, and branch commands as repository-domain work even when they look like separate features.
- For repository deletion, prefer the guarded `gt repo delete` workflow instead of raw API deletion, and preserve the interactive confirmation path for repositories with commit history.
