# Worklog

![Go](https://github.com/PossibleLlama/worklog/workflows/Go/badge.svg)

Intented to be a quick productivity CLI tool, to track what you
have done each day.

## Supported versions

- v0.3.4

## Installation

To install worklog, you can download the latest version from
[github][GithubReleases], choose version you want, unzip the
folder, and add the approperiate binary for your machine, to
a location on your machines PATH.

There is both a windows executeable, as well as a linux binary.

### Linux

```bash
curl -L -O https://github.com/PossibleLlama/worklog/releases/download/v0.3.0/worklog-binary.zip
unzip worklog-binary.zip
mv worklog /usr/bin/worklog
rm worklog-binary.zip
rm worklog.exe
```

[GithubReleases]: https://github.com/PossibleLlama/worklog/releases

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

- `--description "bar"` A longer description of the work done. This
  gives the details of the work done.
- `--duration 30` How long the work took. This can be any unit of
  measurement that suits you.
- `--tags "buzz, bang"` A comma seperate list of tags to describe
  the work.
- `--when "2000/12/31"` Timestamp of when the work was done. This
  should be in [RFC3339] format, either just the date, or datetime.

[RFC3339]: https://tools.ietf.org/html/rfc3339

### Example create

``` bash
worklog create \
--title "Wake up" \
--description "A detailed description of my morning routine."
--duration 60
--tags "morning"
```

## Reading worklogs

``` bash
worklog print <DATE> <FILTERS> <FORMAT>
```

Printing all worklogs to the console that match your criteria.

### Date

One of the following must be provided within the flags.

- `--startDate "2000/12/31"` Print all work completed after this
  date.
- `--endDate "2001/01/01"` Print all work between `startDate` and
  `endDate`. Must be used in conjunction with `--startDate`.
- `--today` Print all work complete since midnight. This is not the
  last 24 hours.
- `--thisWeek` Print all work completed since midnight on the last
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

- `--pretty` Output format is text. (Default)
- `--yaml` Output format is yaml.
- `--json` Output format is json.

### Example print

``` bash
worklog print --today --json --tags "morning"
```

## Configuration

``` bash
worklog configure
```

Configuration for the application.

For basic setup you can run `worklog configure`.
This will provide an empty string for the author, and a
duration of 15.

For more advanced setup, `worklog configure defaults <FLAGS>`
will add provided flags into the configuration.

- `--author "Alice"` String of the author's name.
- `--duration 15` Default duration that a task takes.
- `--format "json"` Default format to print output.
  Accepts `"pretty"`, `"yaml"` or `"json"`.

### Example configure

``` bash
worklog configure defaults --author "Alice" --duration 30 --format "pretty"
```

### Example file

``` yml
author: "Alice"
default:
  duration: 15
```
