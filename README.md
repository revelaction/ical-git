<p align="center"><img alt="go-srs" src="logo.png"/></p>

[![Test](https://github.com/revelaction/ical-git/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/test.yml)
[![Test](https://github.com/revelaction/ical-git/actions/workflows/build.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/ical-git)](https://goreportcard.com/report/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/ical-git?style=flat)]() 
[![Go Reference](https://pkg.go.dev/badge/github.com/revelaction/ical-git)](https://pkg.go.dev/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/badge/built_with-Go-00ADD8.svg?style=flat)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![GitHub repo size](https://img.shields.io/github/repo-size/revelaction/ical-git)
![GitHub stars](https://img.shields.io/github/stars/revelaction/ical-git?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/revelaction/ical-git?color=red)

**ical-git** is a minimalistic calendar application daemon written in Go. It
reads [a directory of iCalendar files](https://github.com/revelaction/ical-git/tree/master/testdata) (normally the files are directly fetched
from your private git repository) and generates custom notifications based on the icalendar alarm definitions or default alarms
defined in the config file.

# Content

- [Usage](#usage)
- [Features](#features)
- [Installation](#installation)
  - [Binary](#binary)
    - [Get the binary](#get-the-binary)
    - [Build Manually](#build-manually)
  - [systemd Service File](#systemd-service-file)
  - [Ical files](#ical-files)
  - [Configuration file](#configuration-file)
- [Managing iCal Files](#managing-ical-files)
- [Command line options](#command-line-options)

# Usage
<p align="center"><img alt="pixelchanged" src="logo.png"/></p>

The basic usage involves having a private repository containing iCalendar
[`.ics` files](testdata/event-recurrent.ics). You need to provide the SSH key and repository address in the
configuration TOML file. Alternatively, you can provide a path in the current
filesystem. For more details, refer to the [Configuration
File](configuration.md). Once configured, you can run the daemon to start
processing the calendar events.

# Features

- **Focusing on Private Hosted Solution**: Designed to provide a simple and effective private self-hosted calendar solution.
- **Low resources computers**: Supports installation on Raspberry Pi Zero and other cheap microcomputers.
- **Notifications**: Offers support for Telegram bots and local Linux desktop notifications.
- **Direct Git Fetching**: Can fetch iCal files directly from a git repository without saving them locally.
- **Alarm Support**: Supports alarms defined in the calendar `.ics` files and defined in the config. Alarms defined in the calendar `.ics` files have priority.
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

For instructions on setting up and managing the systemd service file, see the [systemd.md](systemd.md) file.

## Configuration file


Copy the TOML configuration file to the working directory specified in `WorkingDirectory`.

```console
cp icalgit.toml /home/icalgit/icalgit

```
If you prefer to store the TOML file in a different location, specify the path in the `ExecStart` line of the systemd service file:

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

For a description of the configuration `icalgit.toml` file, see the [Configuration File](configuration.md) section.

## ical files

The preferred method for managing iCal files is to use a private Git repository. Provide the SSH key and repository address in the TOML file under `fetcher_git`. 

```toml
[fetcher_git]
private_key_path = "/path/to/ssh/key"
url = "git@github.com:yourusername/your-repo.git"
```

For instructions on managing and creating iCal files, see the [Managing iCal Files](ical.md#managing-ical-files) section.

# Managing iCal Files

For instructions on managing iCal files, see the [Managing iCal Files](ical.md) section.

# Command line options


```console
    ical-git [-c CONF_FILE] 

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


