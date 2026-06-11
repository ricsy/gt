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

## Install In Codex

From the repository root:

```bash
codex plugin marketplace add .
codex plugin add gt@gt-local
```

Useful checks:

```bash
codex plugin marketplace list --json
codex plugin list --json
gh to codex plugin view gt@gt-local
```

## Update During Local Development

When the plugin changes, bump the Codex cachebuster and reinstall:

```bash
python C:\Users\Administrator\.codex\skills\.system\plugin-creator\scripts\update_plugin_cachebuster.py plugins/gt
codex plugin add gt@gt-local
```

After reinstalling, start a new Codex thread so the updated plugin skills are
picked up reliably.
