# storage-cli

macOS-focused disk usage and cautious cache cleanup utility, written in Go.

## Setup

```sh
git clone https://github.com/r14r/storage-cli.git
cd storage-cli
just bootstrap
just precommit
just build
just install
```

Installation uses `INSTALL_DIR` (default `/usr/local/bin`).

## Commands

```sh
storage-cli disks
storage-cli disks --apfs
storage-cli largest --path "$HOME" --limit 30 --min-mb 100
storage-cli folders --path "$HOME"
storage-cli targets
storage-cli cleanup --category homebrew --older-than 30
storage-cli cleanup --category homebrew --older-than 30 --delete --yes
```

Deletion is disabled by default. Explicit category, --delete and --yes are required. System-wide cache and temporary folders are inspection-only. Cache files may be in use: close related applications before deleting. This utility does not clean APFS snapshots.

## Developer workflow

`just bootstrap`, `just fmt`, `just precommit`, `just build`, `just test`, `just release`, `just install`.

`update-cli.yaml` follows the schema-v2 update-cli workflow conventions.
