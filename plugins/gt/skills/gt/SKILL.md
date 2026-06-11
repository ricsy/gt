---
name: gt
description: Inspect and operate Gitee repositories, issues, pull requests, releases, notifications, gists, and related resources through the gt CLI. Use when the user asks for general Gitee help or wants Codex to choose the right gt workflow.
---

# GT

## Overview

Use this skill as the umbrella entrypoint for `gt`-powered Gitee workflows.

Prefer the local `gt` CLI for:

- repository inspection and creation
- issues, labels, and milestones
- pull requests
- releases and webhooks
- notifications
- gists
- user, organization, and activity queries

Use `gt api` only when:

- the user explicitly wants the raw API result
- a higher-level `gt` subcommand is missing
- you need to verify backend behavior against CLI behavior

## Routing

1. Identify the resource type first:
   - `repo`, `branch`, `collaborator`, `fork`
   - `issue`, `label`, `milestone`
   - `pr`
   - `release`, `webhook`, `check`
   - `notification`, `gist`, `user`, `org`, `activity`
2. Route to a narrower GT skill when the request is clearly domain-specific.
3. Keep raw `gt api` out of the default path unless it is needed for verification or fallback.

## Output Expectations

- Say which `gt` domain you are using.
- Prefer concrete commands or command execution over abstract advice.
- When a command fails, distinguish CLI defects from platform/resource constraints.
