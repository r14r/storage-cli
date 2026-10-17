package storage

import (
 "fmt"
 "io"
 "strings"
)

func help(w io.Writer,args []string)error{
 command:=""
 if len(args)>0 {command=args[0]}
 topics:=map[string]string{
 "disks":`DISK USAGE
  storage-cli disks [--path PATH] [--apfs]
    --path PATH      Filesystem path; default /
    --apfs           Show all macOS APFS containers and volumes (macOS only)
  Examples:
    storage-cli disks
    storage-cli disks --path /System/Volumes/Data
    storage-cli disks --apfs`,
 "folders":`DIRECTORY SIZES
  storage-cli folders [FOLDER] [--depth LEVEL] [--limit COUNT] [--json]
    FOLDER           Directory to inspect; default .
    --path PATH      Alternative to positional FOLDER
    --depth LEVEL    0 = no children, 1 = immediate files/folders (default),
                     2 = also their children, N = through level N
    --limit COUNT    Maximum output rows; 0 = unlimited (default)
    --json           Structured JSON instead of text
  Folder sizes include all descendants, even beyond --depth.
  Symbolic links are not followed. Entries are sorted by depth and size.
  Examples:
    storage-cli folders ~
    storage-cli folders /Users --depth 2
    storage-cli folders /System/Volumes/Data --depth 3 --limit 50
    storage-cli folders --path . --depth 2 --json`,
 "largest":`LARGEST FILES
  storage-cli largest [--path PATH] [--limit COUNT] [--min-mb SIZE] [--json]
    --path PATH      Directory to scan; default .
    --limit COUNT    Maximum results; positive integer (default 25)
    --min-mb SIZE    Minimum file size in MiB; non-negative (default 1)
    --json           Structured JSON output
  Examples:
    storage-cli largest --path ~ --limit 30 --min-mb 100
    storage-cli largest --path /Users --json`,
 "targets":`CLEANUP TARGETS
  storage-cli targets
  Lists detected cache locations and whether deletion is supported.
  Example:
    storage-cli targets`,
 "cleanup":`CACHE INSPECTION AND CLEANUP
  storage-cli cleanup [--category VALUE] [--older-than DAYS] [--limit COUNT]
                      [--json] [--delete --yes]
    --category VALUE  all (default), app-cache, homebrew, system-cache, system-temp
    --older-than DAYS Minimum age in days (default 7; non-negative)
    --limit COUNT     Maximum preview rows (default 40; positive)
    --json            JSON preview (cannot be combined with --delete)
    --delete          Delete eligible regular files
    --yes             Required with --delete, plus an explicit category
  Only app-cache and homebrew permit deletion; system locations are report-only.
  Examples:
    storage-cli cleanup
    storage-cli cleanup --category homebrew --older-than 30
    storage-cli cleanup --category app-cache --delete --yes`,
 "version":`VERSION
  storage-cli version
  storage-cli --version`,
 }
 if command!=""{
  if v,ok:=topics[command];ok{fmt.Fprintln(w,v);return nil}
  return fmt.Errorf("unknown help topic %q",command)
 }
 fmt.Fprintln(w,"storage-cli",Version)
 fmt.Fprintln(w,strings.Join([]string{
  "",
  "DISK INSPECTION",
  "  disks       Filesystem usage and macOS APFS details",
  "  largest     Largest regular files in a directory",
  "  folders     File and directory sizes, with --depth",
  "",
  "CLEANUP",
  "  targets     Available cache locations and cleanup safety",
  "  cleanup     Preview or explicitly delete eligible cache files",
  "",
  "INFORMATION",
  "  help        Show this help or help for one command",
  "  version     Print application version",
  "",
  "For all parameters, valid values and examples:",
  "  storage-cli help folders",
  "  storage-cli help cleanup",
  "  storage-cli <command> --help",
  "",
  "Quick examples:",
  "  storage-cli folders ~ --depth 2",
  "  storage-cli disks --apfs",
  "  storage-cli largest --path ~ --limit 20",
  "  storage-cli cleanup --category homebrew --older-than 30",
 },"\n"))
 return nil
}
