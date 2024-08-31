<p align="center"><img alt="go-srs" src="logo.png"/></p>

[![GitHub Release](https://img.shields.io/badge/built_with-Go-00ADD8.svg?style=flat)]()
[![Test](https://github.com/revelaction/ical-git/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/test.yml)
[![Test](https://github.com/revelaction/ical-git/actions/workflows/build.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/ical-git)](https://goreportcard.com/report/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/ical-git?style=flat)]() 

**ical-git** is a simplistic calendar application daemon written in Go. It reads a
directory for icalendar files and generates custom notifications defined in a
config file.

# Content

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Getting started](#getting-started)
- [Examples](#examples)

## Features

- **Version Control**: Utilizes iCal files under version control.
- **Notifications**: Supports Telegram bots and local Linux desktop notifications.

## Installation

### Get the binary

**ical-git** is a go binary.

On Linux, macOS, FreeBSD you can use the [pre-built binaries](https://github.com/revelaction/ical-git/releases/) 

If your system has a supported version of Go, you can build from source

```console
go install github.com/revelaction/ical-git/cmd/mankidown@latest
```

Move the binary to a suitable path

```console
mv incal-git /home/icalgit/bin/ical-git
chmod +x /home/icalgit/bin/ical-git
```

### Modify the service file

``` 
[Service]
User=icalgit
Group=icalgit

Type=simple
WorkingDirectory=/home/icalgit/icalgit
ExecStart=/home/icalgit/bin/ical-git
Restart=on-failure
TimeoutSec=10
```

Change the `user` and `group` to your own user.
Make sure `WorkingDirectory` and `ExecStart` points to existing own paths.


### Copy the Service File to the Systemd Directory

You need to copy your .service file to the `/etc/systemd/system`
directory. This is where systemd looks for service files. You'll need superuser
privileges to do this. 

```console
sudo cp ical-git.service /etc/systemd/system/
```

### Reload Systemd Daemon: 

After adding or modifying any service file, you need to reload the systemd daemon so it can recognize the new or changed service file.
```console
sudo systemctl daemon-reload
```

### Enable the Service

i8f you want your service to start automatically on boot, you should enable it.
This step is optional but recommended for most services that you want to run
continuously.

```console
sudo systemctl enable ical-git.service 
```

### Start the Service

Now, you can start your service. If you've enabled it, this step is technically
optional since it will start on the next boot, but you'll likely want to start
it immediately for testing purposes.

```console
sudo systemctl start ical-git.service 
```

### Check the Service Status 
To verify that your service is running as expected, you can check its status.

```console
sudo systemctl status ical-git.service 
```

### See logs

```console
sudo journalctl -u ical-git.service 
```


## Usage




## Getting Started


## Examples


```console
crontab -e -u icalgit

```
```cron
* * * * * cd /home/icalgit/icalgit/mi-ical-files && GIT_SSH_COMMAND="ssh -i /home/icalgit/.ssh/id_icalgit_nopassphrase -o IdentitiesOnly=yes" git pull origin main 
```

