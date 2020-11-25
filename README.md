# Worklog

![Go](https://github.com/PossibleLlama/worklog/workflows/Go/badge.svg)

Intented to be a quick productivity CLI tool, to track what you
have done each day.

## Supported versions

- v0.2.0

## Installation

To install worklog, you can download the latest version from
[github][GithubReleases], choose version you want, unzip the
folder, and add the approperiate binary for your machine, to
a location on your machines PATH.

There is both a windows executeable, as well as a linux binary.

### Linux

```bash
curl -L -O https://github.com/PossibleLlama/worklog/releases/download/v0.2.0/worklog-binary.zip
unzip worklog-binary.zip
mv worklog /usr/bin/worklog
rm worklog-binary.zip
```

[GithubReleases]: https://github.com/PossibleLlama/worklog/releases

## Configuration

Configuration for the application.

For basic setup you can run `worklog configure`.

For more advanced setup, `worklog configure defaults <FLAGS>`
will add provided flags into the configuration.

### Example

``` yml
author: "Alice"
default:
  duration: 15
```
