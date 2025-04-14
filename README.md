# Trxsh

[![build](https://github.com/datsfilipe/trxsh/actions/workflows/build.yml/badge.svg)](https://github.com/datsfilipe/trxsh/actions/workflows/build.yml)

Trxsh is a simple and efficient command-line trash manager written in Go. It allows you to safely delete, list, restore, and permanently remove files from your system's trash directory.

### Features

- deletion: move files to a designated trash directory instead of permanently deleting them.​
- list: biew all files currently in the trash.​
- restore: recover files from the trash by ID or using an interactive fzf interface.​
- cleanup: permanently delete all files in the trash.​

### Usage

```bash
Usage: ./dist/trxsh [OPTIONS] [FILES]

Options:
  --fzf, -f         : Restore files using fzf
  --list, -l        : List files in trash
  --restore, -r ID  : Restore file by ID
  --cleanup, -c     : Empty all trash directories
  --help, -h        : Show this help
```

### Config

No configuration is available right now. By default, trxsh uses the `~/.Trash` directory to store trashed files and important information.

### License

This project is licensed under the [MIT License](./LICENSE).
