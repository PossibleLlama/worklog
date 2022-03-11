# Worklog

![Build](https://github.com/PossibleLlama/worklog/workflows/Go/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/PossibleLlama/worklog)](https://goreportcard.com/report/github.com/PossibleLlama/worklog)

Intented to be a quick productivity CLI tool, to track what you
have done each day.

## Supported versions

- v0.6.x

## Installation

To install worklog, you can download the latest version from
[github][GithubReleases], choose version you want, unzip the
folder, and add the approperiate binary for your machine, to
a location on your machines PATH.

There is both a windows executeable, as well as a linux binary.

### Linux

```bash
curl -L -O $(curl -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/possiblellama/worklog/releases | jq -r .[0].assets[0].browser_download_url)
unzip worklog-binary.zip
mv 64bit/worklog-linux /usr/bin/worklog
rm -rf 32bit 64bit
rm worklog-binary.zip
```

[GithubReleases]: https://github.com/PossibleLlama/worklog/releases

## Global flags

```bash
worklog --config "/path/to/config" <command>
```

If you store the configuration file somewhere other than
`$HOME/.worklog/config.yaml`, you will need to specify where, via
this flag.
It will have to be provided everytime you use this, as such, if
you want to permanently use another location, setting this as an
alias will likely be most useful.

```bash
worklog --repo "bolt"
```

To use a different way of storing the worklogs, you can change it
via this flag.
The currently available options are listed below, although more
may be available in the future.

This will be most useful when changing between storage types.
Currently the default is `bolt`, however up until 0.6.1, the
default repo was `legacy`, and as such to search those, you'll
need to specify this repo type when printing.
You can also specify the default repo type via the `configure`
command.

Valid options:

- `bolt`
- `legacy`

```bash
worklog --repoPath "/path/to/repo"
```

This will be to specify where the store of worklogs is.
If it isn't in the default location of `$HOME/.worklog/worklog.db`
you will need to specify this.
You can also specify the default repo type via the `configure`
command.

This is not supported by the `legacy` repo type.

## Creating worklogs

``` bash
worklog create <FLAGS>
```

To create a basic worklog, you can use `worklog create --title "foo"`.
This will create an item of work, with the current timestamp, and the
name `"foo"`.

Author and duration will automatically be pulled from the
configuration file if provided, and createdAt will be the current
time.
All other fields will be left blank.

You can specify further fields as you want to.

- `--title "foo"` The title of the work done. This is a summary of
  the work.
- `--description "bar"` A longer description of the work done. This
  gives the details of the work done.
- `--author "Alice"` The name of the person doing the work. This
  will override the default value in the config file.
- `--duration 30` How long the work took. This can be any unit of
  measurement that suits you. This will override the default value
  in the config file.
- `--tags "buzz, bang"` A comma seperate list of tags to describe
  the work.
- `--when "2000/12/31"` Timestamp of when the work was done. This
  should be in [RFC3339] format, either just the date, or datetime.

[RFC3339]: https://tools.ietf.org/html/rfc3339

### Example create

``` bash
worklog create \
--title "Wake up" \
--description "A detailed description of my morning routine." \
--author "Alice" \
--duration 60 \
--tags "morning"
```

## Editing worklogs

``` bash
worklog edit <ID> <FLAGS>
```

To edit an existing worklog, you will need to know the ID, then you can
use the command `worklog edit "abc" --description "New description"`.
To find this ID, you will need to have printed that worklog previously,
and enough of the ID must be used to make it unique amongst all other
worklogs.

The revision and createdAt fields will automatically be updated when
editing and can't be specified.

Any specified fields will be merged into the existing work overwritting
previous fields.

Any blank `--title ""` fields will be ignored, and all strings have
whitespace removed from both ends.
This does mean that you will be unable to update a field to an empty
state.

You can specify further fields as you want to.

- `--title "foo"` The title of the work done. This is a summary of
  the work.
- `--description "bar"` A longer description of the work done. This
  gives the details of the work done.
- `--author "Alice"` The name of the person doing the work. This
  will override the default value in the config file.
- `--duration 30` How long the work took. This can be any unit of
  measurement that suits you. This will override the default value
  in the config file.
- `--tags "buzz, bang"` A comma seperate list of tags to describe
  the work.
- `--when "2000/12/31"` Timestamp of when the work was done. This
  should be in [RFC3339] format, either just the date, or datetime.

### Example edit

``` bash
worklog edit \
"abc"
--title "Better summary" \
--description "Made my coffee" \
--author "Bob" \
--duration 10 \
--tags "morning"
```

## Reading worklogs

``` bash
worklog print <DATE> <FILTERS> <FORMAT> <IDS>
```

Printing all worklogs to the console that match your criteria.

By default a subset of fields will be displayed, however the `--all`
flag can be used to show all of them.

### Date

One of the following must be provided within the flags.

- `--startDate "2000/12/31"` Print all work completed after this
  date.
- `--endDate "2001/01/01"` Print all work between `startDate` and
  `endDate`. Must be used in conjunction with `--startDate`.
- `--today`, `-t` Print all work complete since midnight. This is not the
  last 24 hours.
- `--thisWeek`, `-w` Print all work completed since midnight on the last
  Monday. This is not the last 168 hours.

### Filters

You can specify as many of these as you need to search for the correct
worklogs.

All filters:

- Are case insensitive.
- Will include partial matches (`--title "a"` will return all titles
  that include an "a" anywhere within the title.)
- Any returned Work must satisfy all filters.

Arguments match the names used when creating a worklog.

Valid arguments are:

- `--title`
- `--description`
- `--author`
- `--tags`

### Format

Optionally, you can provide 1 of the following which will change how
the output is formatted.

- `--pretty`, `-p` Output format is text. (Default)
- `--yaml`, `-y` Output format is yaml.
- `--json`, `-j` Output format is json.

### Misc

Other fields that can additionally be used.

- `--all`, `-a` Output all fields.

### IDs

Id's of the specific worklogs that should be printed.
This overrides the DATE and FILTERS flags.

If any of the IDs is not unique, an error will be thrown.
If the ID does not exist, a message will be printed to that affect.

### Example print

``` bash
worklog print --today --json --tags "morning"
worklog print "abc" "def"
```

## Configuration

``` bash
worklog configure
```

Configuration for the application.

For basic setup you can run `worklog configure`.
This will provide an empty string for the author, and a
duration of 15.

For more advanced setup, `worklog configure overrideDefaults <FLAGS>`
will add provided flags into the configuration.

- `--author "Alice"` String of the author's name.
- `--duration 15` Default duration that a task takes.
- `--format "json"` Default format to print output.
  Accepts `"pretty"`, `"yaml"` or `"json"`.

### Example configure

``` bash
worklog configure overrideDefaults --author "Alice" --duration 30 --format "pretty"
```

### Example file

``` yml
default:
  author: "Alice"
  duration: 15
  format: "pretty"
```
