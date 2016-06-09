# gocal
Simple cli tool for google calendar api.

## Usage
```
usage: gocal [<flags>] <command> [<args> ...]

google calendar events api

Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
  -c, --conf="$HOME/.config/gocal/calendar.toml"  
                 config file
      --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  events list [<flags>]
    insert google calendar events

  events insert --name=NAME [<flags>]
    insert google calendar events

```

### Insert events from JSON file

`gocal` supports to insert from JSON file.

```
$ gocal events insert --json-file /path/to/jsonfile.json
```

example JSON:

```json
[
  {
      "title": "input title",
      "detail": "input a description in this field"
  },
  ...
]
```

## Config

The format is toml, and default path is `$HOME/.config/gocal/calendar.toml`.

```
calendar_id = "your calendar id"
credential_file = "/path/to/credential_file.json"
```

You can create `credential_file.json` on Google API Console at Service Account Manager.
