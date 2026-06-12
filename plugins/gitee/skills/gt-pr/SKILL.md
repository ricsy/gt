---
name: gt-pr
description: Manage Gitee pull requests through gt, including listing, viewing, creating, commenting, closing, state changes, review and test gate handling, and merge operations. Use when the request is primarily about pull requests.
---

# GT PR

## Scope

Use this skill for:

- `gt pr list/view/create/comment/state/close/review/test/merge`

## Workflow

1. Resolve the repository and branch context first.
2. Prefer `gt pr create` and `gt pr merge` over raw API calls.
3. For repositories with PR gates, prefer explicit gate handling before merge:
   - `gt pr review <number> --repo owner/repo`
   - `gt pr test <number> --repo owner/repo`
   - `gt pr merge <number> --repo owner/repo --pass-gates`
4. Use `--force` only when the requester is an administrator and explicitly wants to force review or test pass:
   - `gt pr review <number> --repo owner/repo --force`
   - `gt pr test <number> --repo owner/repo --force`
   - `gt pr merge <number> --repo owner/repo --pass-gates --force`
5. When merge fails, separate platform policy failures from CLI defects.
6. If needed, use `gt api` only to confirm backend semantics.

## Guardrails

- Do not describe a repository policy failure as a broken merge command.
- If a repository requires review or test approval, do not tell the user to click the web UI first until you have checked whether `gt pr review`, `gt pr test`, or `gt pr merge --pass-gates` can satisfy the gate.
- Treat comments and state transitions as part of the PR workflow, not the issue workflow.
- Treat missing gate APIs as a CLI capability gap, not as proof that the merge endpoint itself is wrong.
