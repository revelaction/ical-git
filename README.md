<p align="center"><img alt="go-srs" src="logo.png"/></p>

[![GitHub Release](https://img.shields.io/badge/built_with-Go-00ADD8.svg?style=flat)]()
[![Test](https://github.com/revelaction/ical-git/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/test.yml)
[![Test](https://github.com/revelaction/ical-git/actions/workflows/build.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/ical-git)](https://goreportcard.com/report/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/ical-git?style=flat)]() 

**ical-git** is a simplistic calendar application written in Go. It reads a
directory for icalendar files and generates custom notifications defined in a
config file.

## Features

- **Version Control**: Utilizes iCal files under version control.
- **Notifications**: Supports Telegram bots and local Linux desktop notifications.

## Installation

To install ical-git, follow these steps:

## Usage

## Getting Started


## Examples


```console
crontab -e -u icalgit

```
```cron
* * * * * cd /home/icalgit/icalgit/mi-ical-files && GIT_SSH_COMMAND="ssh -i /home/icalgit/.ssh/id_icalgit_nopassphrase -o IdentitiesOnly=yes" git pull origin main 
```

