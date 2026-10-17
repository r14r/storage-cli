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

## Detailed help and recursive folder sizes (v0.2.0)

```sh
storage-cli help
storage-cli help folders
storage-cli folders ~/Downloads
storage-cli folders ~/Downloads --depth 2
storage-cli folders /System/Volumes/Data --depth 3 --json
```

`--depth 1` displays the direct children of the selected folder, including files and directories. Larger values additionally display nested descendants. Each directory reports its complete recursive size, independent of display depth. `--limit 0` shows all rows.
