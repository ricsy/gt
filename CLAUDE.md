# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Tool Preference

**Prefer Goland MCP (`mcp__goland__*`) over direct file editing for code operations.** Use it for reading, searching, replacing, and navigating code. Only use direct file tools (Read, Edit, Write) when MCP tools are insufficient.

## API Reference

**MCP `swagger-gitee`**: Query Gitee Open API v5 endpoints directly.
- `mcp__swagger-gitee__list_api_groups` — list all API groups
- `mcp__swagger-gitee__search_apis` — search APIs by keyword
- `mcp__swagger-gitee__get_api_detail` — get full API definition (path, method, parameters, schema)

`data/api-1.json` is the local copy of the same specification (sourced from https://help.gitee.com/openapi/v5).

## Project Overview

`gt` is a CLI tool for Gitee (gitee.com). It manages repos, issues, PRs, releases, and organizations via the Gitee API v5.

## Common Commands

```bash
go build                          # Build binary (outputs gt/gt.exe)
go test ./...                     # Run all tests
go test -v ./internal/cmd        # Run cmd package tests
go test -v ./pkg/api             # Run API package tests
golangci-lint run ./...           # Lint all code
pnpm run test:integration         # Run integration tests (requires built binary)
pnpm run release:local            # Local release build
```

## Architecture

```
internal/cmd/    # CLI command implementations (Cobra commands)
pkg/api/         # Gitee API client (HTTP calls + response types)
pkg/auth/        # Authentication (token storage, cache)
pkg/config/      # Config file and hosts file management
pkg/util/        # Shared utilities
```

### API Client (`pkg/api/client.go`)

`Client` struct wraps HTTP calls to Gitee API v5. All API paths use `config.ApiUrl(host)` as base. Auth header format is `Authorization: token <token>`.

Each domain (issue, pr, repo, org, release) has its own file in `pkg/api/`.

### Auth Flow (`pkg/auth/auth.go`)

`GetToken(host)` checks `GITEE_TOKEN` env var first, then falls back to `config.LoadHosts()` (YAML file). Results are cached in-memory with double-checked locking.

### Repo Resolution (`internal/cmd/repo_helper.go`)

Use `getClient()` to get an authenticated API client. Use `resolveRepoFlag(flag)` to resolve repo from flag -> `GT_REPO` env var -> `default_repo` config. Use `ResolveRepo(owner/repo)` to parse and validate `owner/repo` strings.

## Key Patterns

- All command handlers use `getClient()` + `resolveRepoFlag()` — never call `auth.GetToken` + `api.NewClient` directly
- URL construction always goes through `config.ApiUrl(host)` — no hardcoded `https://gitee.com`
- `api.StateOpen` constant exists for issue/PR state filtering (value: "open")
