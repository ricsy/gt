# GT Codex Plugin

This plugin packages skill-based Codex workflows for the `gt` CLI.

The first version is intentionally skills-first:

- no connector dependency
- no bundled app mapping
- local `gt` CLI as the execution path

It is designed to route Codex to the right Gitee workflow domain and then use
the local `gt` command or `gt api` fallback where needed.

## Layout

- plugin root: `plugins/gt`
- repo marketplace: `.agents/plugins/marketplace.json`

The marketplace entry uses `./plugins/gt`, which is resolved relative to the
repository root.

## Manage In Codex

Install from the GitHub repository source:

```bash
codex plugin marketplace add ricsy/gt
codex plugin add gitee@ricsy-curated
codex plugin remove gitee@ricsy-curated
```

Useful checks:

```bash
codex plugin marketplace list --json
codex plugin list --json
```

After reinstalling, start a new Codex thread so the updated plugin skills are
picked up reliably.
