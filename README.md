# Worklog

![Go](https://github.com/PossibleLlama/worklog/workflows/Go/badge.svg)

Intented to be a quick productivity CLI tool, to track what you
have done each day.

## Supported versions

- v0.3.1

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

## Reading worklogs

Printing worklogs to the console. You can use `worklog print <FLAGS>`.
One of the following must be provided within the flags.

- `--startDate "2000/12/31"` Print all work completed after this
  date.
- `--today` Print all work complete since midnight. This is not the
  last 24 hours.
- `--thisWeek` Print all work completed since midnight on the last
  Monday. This is not the last 168 hours.

## Configuration

Configuration for the application.

For basic setup you can run `worklog configure`.

For more advanced setup, `worklog configure defaults <FLAGS>`
will add provided flags into the configuration.

- `--author "Alice"` String of the author's name.
- `--duration 15` Default duration that a task takes.

### Example file

``` yml
author: "Alice"
default:
  duration: 15
```
