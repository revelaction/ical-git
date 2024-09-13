<p align="center"><img alt="go-srs" src="logo.png"/></p>

[![GitHub Release](https://img.shields.io/badge/built_with-Go-00ADD8.svg?style=flat)]()
[![Test](https://github.com/revelaction/ical-git/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/test.yml)
[![Test](https://github.com/revelaction/ical-git/actions/workflows/build.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/ical-git)](https://goreportcard.com/report/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/ical-git?style=flat)]() 
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**ical-git** is a minimalistic calendar application daemon written in Go. It
reads a directory of iCalendar files (normally the files are directly fetched from your private git repository) and generates custom notifications
based on the icalendar alarm definitions or default alarms defined in the config file.

# Content

- [Usage](#usage)
- [Features](#features)
- [Installation](#installation)
  - [Binary](#binary)
    - [Get the binary](#get-the-binary)
    - [Build Manually](#build-manually)
  - [systemd Service File](#systemd-service-file)
  - [ical files](#ical-files)
  - [ical configuration file](#ical-configuration-file)
- [Configuration File](#configuration-file)
- [Creating iCalendar Files](#creating-icalendar-files)
- [Managing iCal Files](#managing-ical-files)
- [Command line options](#command-line-options)

# Usage

The basic usage involves having a private repository containing iCalendar (`.ics`) files. You need to provide the SSH key and repository address in the configuration TOML file. Once configured, you can run the daemon to start processing the calendar events.

# Features

- **Conceptually Simple Design**: Offers a conceptually simple design that allows for a relatively simple private self-hosted calendar solution.
- **Low resources computers**: Supports installation on Raspberry Pi Zero and other cheap microcomputers.
- **Notifications**: Offers support for Telegram bots and local Linux desktop notifications to keep users informed.
- **Customizable Configuration**: Allows users to define custom notifications through a TOML configuration file.
- **Systemd Integration**: Facilitates seamless integration with systemd for service management and logging.

# Installation

## Binary
### Get the binary

On Linux, macOS, FreeBSD you can use the [pre-built binaries](https://github.com/revelaction/ical-git/releases/) 

If your system has a supported version of Go, you can build from source

```console
go install github.com/revelaction/ical-git/cmd/ical-git@latest
```

Move the binary to a suitable path

```console
mv incal-git /home/icalgit/bin/ical-git
chmod +x /home/icalgit/bin/ical-git
```

### Build Manually

To build `ical-git` manually from the source code, follow these steps:

1. **Clone the Repository**: Clone the `ical-git` repository to your local machine.

    ```console
    git clone https://github.com/revelaction/ical-git.git
    cd ical-git
    ```

2. **Build the Binary**: Use `go build` with `ldflags` to include the Git tag in the binary.

    ```console
    go build -ldflags "-X main.BuildTag=$(git describe --tags)" ./cmd/ical-git
    ```

3. **Move the Binary**: Move the built binary to a suitable path and set the executable permission.

    ```console
    mv ical-git /home/icalgit/bin/ical-git
    chmod +x /home/icalgit/bin/ical-git
    ```

## systemd Service File

For detailed instructions on setting up and managing the systemd service file, please refer to the [systemd.md](systemd.md) file.

## ical configuration file

Copy the TOML configuration file to the working directory specified in `WorkingDirectory`.

```console
cp icalgit.toml /home/icalgit/icalgit
```

Update the `directory` path in the TOML file to point to your iCal files directory.

```toml
...

[fetcher_filesystem]
directory = "/home/ical/path/to/my-ical-files"

...
```

If you prefer to store the TOML file in a different location, specify the path in the `ExecStart` line of the systemd service file.

``` 
[Service]
User=icalgit
Group=icalgit

Type=simple
WorkingDirectory=/home/icalgit/icalgit
ExecStart=/home/icalgit/bin/ical-git --config /path/to/my-file.toml
Restart=on-failure
TimeoutSec=10
```

## ical files

Place your iCal files in a directory of your choice, preferably under revision
control. Ensure that these files are located in the working directory of the
service. For instructions on creating iCalendar files, refer to the [Creating iCalendar Files](#creating-icalendar-files) section.

```console
mkdir /home/icalgit/icalgit/my-ical-files
```

If you prefer to store them in a different directory, specify the path in the
TOML configuration file. To specify a different directory for your iCal files,
update the TOML configuration file as follows:

```toml
[fetcher_filesystem]
directory = "/home/icalgit/path/to/my-cal-files"
```

# Command line options


```console
    ical=git [-c CONF_FILE] 

Options:
    -c, --config                load the configuration file at CONF_FILE instead of default
    -v, --version               Print the version 
    -h, --help                  Show this

CONF_FILE is the toml configuration file 

ical-git will react to a SIGHUP signal reloading the configuration file.

Examples:
    $ ical-git --config /path/to/config/file.toml # start the daemon with the configuration file
    $ ical-git -v  # print version`
```


