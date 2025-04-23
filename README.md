# Trxsh

[![build](https://github.com/datsfilipe/trxsh/actions/workflows/build.yml/badge.svg)](https://github.com/datsfilipe/trxsh/actions/workflows/build.yml)

Trxsh is a simple and efficient command-line trash manager written in Go. It allows you to safely delete, list, restore, and permanently remove files from your system's trash directory.

### Table of Contents

- [Installation](#installation)
- [Features](#features)
- [Usage](#usage)
- [Trash Specification](#trash-specification)
- [License](#license)

### Installation

- Arch:

The package is now in [AUR](https://aur.archlinux.org/packages/trxsh), so to install it you can use an AUR helper like `yay`:

```bash
yay -Syu trxsh
```

or go through the process manually:

```bash
sudo pacman -S --needed git base-devel
git clone https://aur.archlinux.org/trxsh.git
cd yay
makepkg -si
```

- NixOs

For now you can add a custom derivation to your NixOS config, like so:

```nix
{
  stdenv,
  lib,
  fetchurl,
}: let
  source = builtins.fromJSON (builtins.readFile ./conf/source.json);
in
  stdenv.mkDerivation rec {
    pname = "trxsh";
    version = source.version;

    src = fetchurl {
      url = "https://github.com/datsfilipe/trxsh/releases/download/${version}/trxsh-${version}-linux-amd64.tar.gz";
      sha256 = source.sha256;
    };

    installPhase = ''
      mkdir -p $out/bin
      tar -xzf $src -C $out/bin
      chmod +x $out/bin/trxsh
    '';

    meta = with lib; {
      description = "trxsh is a terminal-based trash manager";
      homepage = "https://github.com/datsfilipe/trxsh";
      license = licenses.mit;
      platforms = ["x86_64-linux"];
    };

    unpackPhase = ":";
    dontStrip = true;
  }
```

In future I'll try adding it to `nixpkgs`.

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
