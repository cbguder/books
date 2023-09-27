# Books

Use [Libby](https://libbyapp.com/) and [Goodreads](https://www.goodreads.com/)
from the command line.

## Installation

The easiest way to install `books` is to download a pre-built binary from
the [releases page](https://github.com/cbguder/books/releases).

Pre-built binaries are available for macOS and Linux.

### Dependencies

The `repackage` command depends on [ffmpeg](https://ffmpeg.org/) for audiobooks
and [html-tidy](http://www.html-tidy.org/) for ebooks.

You can install these using Homebrew on the Mac:

    $ brew install ffmpeg tidy-html5

or a package manager on Linux:

    $ sudo apt install ffmpeg tidy

### Installing From Source

If you have Go installed, you can install `books` from source:

    $ go install github.com/cbguder/books@latest

## Usage

```
Usage:
  books [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  goodreads   Goodreads commands
  help        Help about any command
  libby       Libby commands
  repackage   Repackage a book into an epub or MP3
  version     Print version information

Flags:
      --config string   config file (default is $HOME/.books.yml)
  -h, --help            help for books

Use "books [command] --help" for more information about a command.
```

## Example

```
$ books libby auth
Go to Menu > Settings > Copy To Another Device. You will see a setup code. Enter it below.
Setup code: 12345678
Syncing...
Clone successful

$ books libby search -f audiobook "Patrick Radden Keefe Say Nothing"
╭────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ San Francisco Public Library                                                                       │
├─────────┬──────────────────────┬─────────────┬──────┬───────────┬──────────┬───────────┬───────────┤
│ ID      │ Author               │ Title       │ Year │ Type      │ Language │ Available │ Est. Wait │
├─────────┼──────────────────────┼─────────────┼──────┼───────────┼──────────┼───────────┼───────────┤
│ 3984153 │ Patrick Radden Keefe │ Say Nothing │ 2019 │ Audiobook │ English  │ true      │         2 │
╰─────────┴──────────────────────┴─────────────┴──────┴───────────┴──────────┴───────────┴───────────╯

$ books libby borrow 3984153

$ books libby loans
Syncing...
╭─────────┬──────────────────────┬─────────────┬───────────┬─────────┬───────────────┬────────────╮
│ ID      │ Author               │ Title       │ Type      │ Library │ Checkout Date │ Due Date   │
├─────────┼──────────────────────┼─────────────┼───────────┼─────────┼───────────────┼────────────┤
│ 3984153 │ Patrick Radden Keefe │ Say Nothing │ Audiobook │ sfpl    │ 2023-08-22    │ 2023-09-12 │
╰─────────┴──────────────────────┴─────────────┴───────────┴─────────┴───────────────┴────────────╯

$ books libby download 3984153
Syncing...
Opening audiobook...
Downloading 17 files to "Patrick Radden Keefe - Say Nothing"...

$ books libby return 3984153
Syncing...

$ books repackage "Patrick Radden Keefe - Say Nothing"
```

## Caveats

* The `repackage` command is experimental, and is likely to fail or produce
  invalid epub files. Validating the output with
  [EPUBCheck](https://www.w3.org/publishing/epubcheck/) is recommended.

## See Also

- [ping/odmpy](https://github.com/ping/odmpy) is a full-featured alternative
  written in Python
