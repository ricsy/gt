---
name: gt-mode
description: Manage gt repository scope mode and temporary env-file driven execution context. Use when the request is about switching between personal and organization repository context, avoiding accidental operations in the wrong namespace, or running commands/tests with a temporary account or org setup.
---

# GT Mode

## Scope

Use this skill for:

- `gt repo mode show`
- `gt repo mode personal`
- `gt repo mode org <namespace>`
- `gt repo mode clear`
- `gt --env-file <path> ...`
- temporary environment-driven repo scope such as `GT_REPO_SCOPE_MODE`, `GT_REPO_SCOPE_NAMESPACE`, and `GT_ENV_FILE`
- live test setup that needs a temporary owner, repo, or base branch without rewriting local persistent config

## Workflow

1. Determine whether the user wants a persistent scope change or a temporary execution context.
2. For persistent repository context, prefer `gt repo mode ...`.
3. For temporary command or test execution, prefer `--env-file` or `GT_ENV_FILE`.
4. When the request involves live integration testing, check whether `GT_INTEGRATION_TESTS`, `GT_TEST_OWNER`, `GT_TEST_REPO`, and `GT_TEST_PR_BASE` should be supplied through the env file.
5. When the request is about avoiding accidental personal-vs-org operations, explicitly confirm whether `personal`, `org`, or `none` is the correct target mode.

## Guardrails

- Do not hardcode a developer's personal owner, repository, or organization into env examples, tests, fixtures, or workflow docs.
- Prefer official public sample identifiers such as `gitee` when a stable example namespace is needed.
- Use `gt repo mode clear` when the user wants to leave scoped mode and return to neutral behavior.
- Prefer env-file driven temporary context for testing, CI, and one-off live validation instead of mutating long-lived local config.
