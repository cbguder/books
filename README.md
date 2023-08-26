# Books

Use the [Libby](https://libbyapp.com/) API from the command line.

## Usage

```
Usage:
  books [command]

Available Commands:
  borrow      Borrow a media item
  clone       Clone Libby account from another device
  completion  Generate the autocompletion script for the specified shell
  download    Download a borrowed media item
  help        Help about any command
  holds       Show current holds
  library     Search for a library
  loans       Show current loans
  repackage   repackage an ebook into an epub
  return      Return a loaned item
  search      Search for media
  sync        Sync data manually

Flags:
      --config string   config file (default is $HOME/.books.yml)
  -h, --help            help for books

Use "books [command] --help" for more information about a command.
```

## Example

```
$ books clone
Go to Menu > Settings > Copy To Another Device. You will see a setup code. Enter it below.
Setup code: 12345678
Syncing...
Clone successful

$ books search -f audiobook "Patrick Radden Keefe Say Nothing"
╭────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ San Francisco Public Library                                                                       │
├─────────┬──────────────────────┬─────────────┬──────┬───────────┬──────────┬───────────┬───────────┤
│ ID      │ Author               │ Title       │ Year │ Type      │ Language │ Available │ Est. Wait │
├─────────┼──────────────────────┼─────────────┼──────┼───────────┼──────────┼───────────┼───────────┤
│ 3984153 │ Patrick Radden Keefe │ Say Nothing │ 2019 │ Audiobook │ English  │ true      │         2 │
╰─────────┴──────────────────────┴─────────────┴──────┴───────────┴──────────┴───────────┴───────────╯

$ books borrow 3984153

$ books loans
Syncing...
╭─────────┬──────────────────────┬─────────────┬───────────┬───────────────┬────────────╮
│ ID      │ Author               │ Title       │ Type      │ Checkout Date │ Due Date   │
├─────────┼──────────────────────┼─────────────┼───────────┼───────────────┼────────────┤
│ 3984153 │ Patrick Radden Keefe │ Say Nothing │ Audiobook │ 2023-08-22    │ 2023-09-12 │
╰─────────┴──────────────────────┴─────────────┴───────────┴───────────────┴────────────╯

$ books download 3984153
Syncing...
Opening audiobook...
Downloading 17 files to "Patrick Radden Keefe - Say Nothing"...

$ books return 3984153
Syncing...
```

## Caveats

Only supports borrowing audiobooks at the moment.

## See Also

- [ping/odmpy](https://github.com/ping/odmpy) is a full-featured alternative written in Python
