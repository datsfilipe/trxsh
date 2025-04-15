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
  --fzf, -f        : Restore files using fzf
  --list, -l       : List files in trash
  --restore, -r ID : Restore file by ID
  --cleanup, -c    : Empty all trash directories
  --dir-sizes, -s  : Show directory sizes
  --help, -h       : Show this help
```

### Trash Specification

I've tried implementing the basics from [freedesktop.org trash specification](https://specifications.freedesktop.org/trash-spec/1.0/), but I haven't checked if eveything meets the spec. If you feel like this is an important matter, please feel free to address any issue you find with a proper PR, or to point out and I can work on that whenever I have some spare time.

### License

This project is licensed under the [MIT License](./LICENSE).
