<p align="center"><img alt="go-srs" src="logo.png"/></p>

[![GitHub Release](https://img.shields.io/badge/built_with-Go-00ADD8.svg?style=flat)]()
[![Test](https://github.com/revelaction/ical-git/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/test.yml)
[![Test](https://github.com/revelaction/ical-git/actions/workflows/build.yml/badge.svg)](https://github.com/revelaction/ical-git/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/ical-git)](https://goreportcard.com/report/github.com/revelaction/ical-git)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/ical-git?style=flat)]() 

**ical-git** is a simplistic calendar application daemon written in Go. It
reads a directory of iCalendar files and generates custom notifications
defined in a config file.

# Content

- [Features](#features)
- [Installation](#installation)
  - [Binary](#binary)
    - [Get the binary](#get-the-binary)
    - [Build Manually](#build-manually)
  - [systemd service file](#systemd-service-file)
    - [Modify the service file](#modify-the-service-file)
    - [Copy the Service File to the Systemd Directory](#copy-the-service-file-to-the-systemd-directory)
    - [Reload Systemd Daemon](#reload-systemd-daemon)
    - [Enable the Service](#enable-the-service)
    - [Start the Service](#start-the-service)
    - [Check the Service Status](#check-the-service-status)
    - [Reload Configuration on SIGHUP](#reload-configuration-on-sighup)
    - [See logs](#see-logs)
  - [ical files](#ical-files)
  - [ical configuration file](#ical-configuration-file)
- [Configuration File](#configuration-file)
- [Creating iCalendar Files](#creating-icalendar-files)
- [Managing iCal Files](#managing-ical-files)
- [Command line options](#command-line-options)

# Features

- **Conceptually Simple Design**: Offers a conceptually simple design that allows for a relatively simple private self-hosted calendar solution.
- **Low resources computers**: Supports installation on Raspberry Pi Zero and other cheap microcomputers.
- **Notifications**: Offers support for Telegram bots and local Linux desktop notifications to keep users informed.
- **Customizable Configuration**: Allows users to define custom notifications through a TOML configuration file.
- **Systemd Integration**: Facilitates seamless integration with systemd for service management and logging.

# Installation


## Binary
#### Get the binary

**ical-git** is a go binary.

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

#### Build Manually

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

## systemd service file

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

### Reload Configuration on SIGHUP
`ical-git` can reload its configuration file when it receives the SIGHUP signal. To reload the configuration without restarting the service, use the following command:

```console
sudo systemctl reload ical-git.service
```

### See logs

```console
sudo journalctl -u ical-git.service 
```

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






# Creating iCalendar Files

To create your own iCalendar files, you can use a private Language Model (LLM) on your computer. Here are some steps to guide you through the process:

1. **Install a Private LLM**: Ensure you have a private LLM installed on your computer. There are several open-source models available.
2. **Generate iCalendar Content**: Use the LLM to generate iCalendar content. You can provide prompts or templates to guide the generation process.
3. **Format the Output**: Ensure the generated content adheres to the iCalendar format. This includes proper syntax, event details, and other required fields.
4. **Save the File**: Save the generated content as a `.ics` file in the directory specified in your configuration.

Alternatively, you can copy and modify the existing iCalendar files from the `testdata` directory to suit your needs:

1. **Copy Files**: Copy the desired iCalendar files from the `testdata` directory to your working directory.
2. **Modify Content**: Open the copied files and modify the event details, dates, and other relevant information to fit your requirements.
3. **Save Changes**: Save the modified files in the directory specified in your configuration.

# Managing iCal Files

It is highly advisable to place your iCal files under revision control to ensure that changes are tracked and can be reverted if necessary. 
Additionally, setting up a cron job to periodically pull the latest content of these files can help keep your calendar up-to-date.

## Setting Up Revision Control

1. **Initialize a Git Repository**: Navigate to your iCal files directory and initialize a Git repository.

    ```console
    cd /home/icalgit/icalgit/my-ical-files
    git init
    ```

2. **Add and Commit Your Files**: Add your iCal files to the repository and commit them.

    ```console
    git add .
    git commit -m "Initial commit of iCal files"
    ```

3. **Push to a Remote Repository**: If you have a remote repository (e.g., on GitHub), push your local repository to the remote.

    ```console
    git remote add origin git@github.com:yourusername/your-repo.git
    git push -u origin main
    ```

## Setting Up a Cron Job

To ensure your iCal files are periodically updated, you can set up a cron job to pull the latest content from your revision control system.

1. **Edit Your Crontab**: Open your crontab file for editing.

    ```console
    crontab -e -u icalgit
    ```

2. **Add the Cron Job**: Add a cron job to pull the latest changes from your repository. It's advisable to use a passphraseless SSH key to avoid issues with cron executing in a non-interactive session.

    ```cron
    * * * * * cd /home/icalgit/icalgit/my-ical-files && GIT_SSH_COMMAND="ssh -i /home/icalgit/.ssh/id_icalgit_nopassphrase -o IdentitiesOnly=yes" git pull origin main
    ```

This cron job will run every minute, pulling the latest changes from the `main` branch of your repository. Adjust the schedule as needed.

# Configuration File

The `icalgit.toml` file is used to configure the behavior of the ical-git daemon. Below are the descriptions of the fields and their purposes.

## General Settings

- **timezone**: Specifies the timezone for the notifications. Both icalendar timezone and config timezone will be shown in the notification message
  - Example: `timezone = "Europe/Rome"`
- **tick**: Defines the interval at which the daemon checks for new events.
  - Example: `tick = "24h"`

## Alarms

- **alarms**: A list of alarm configurations. Each alarm specifies the type and when it should trigger. The `when` field uses the ISO 8601 duration format to specify the time before the event when the alarm should trigger.
  - Example:
    ```toml
    alarms = [
        {type = "desktop", when = "-P7D"},  # 7 days before the event
        {type = "desktop", when = "-P1D"},  # 1 day before the event
        {type = "desktop", when = "-PT1H"}, # 1 hour before the event
        {type = "telegram", when = "-P7D"}, # 7 days before the event
        {type = "telegram", when = "-P1D"}, # 1 day before the event
        {type = "telegram", when = "-PT1H"},# 1 hour before the event
    ]
    ```
  - ISO 8601 Duration Format Examples:
    - `-P7D`: 7 days before the event
    - `-P1D`: 1 day before the event
    - `-PT1H`: 1 hour before the event
    - `-PT30M`: 30 minutes before the event
    - `-P1DT12H`: 1 day and 12 hours before the event

## Notifiers

- **notifiers**: A list of notifier types that are enabled.
  - Example: `notifiers = ["desktop"]`

## Fetcher Filesystem

- **directory**: Specifies the directory where the iCal files are stored.
  - Example: `directory = "testdata"`

## Notifier Telegram

- **token**: The API token for the Telegram bot.
  - Example: `token = "yuu3b3k"`
- **chat_id**: The chat ID to which the notifications will be sent.
  - Example: `chat_id = 588488`

## Notifier Desktop

- **icon**: The path to the icon file used for desktop notifications.
  - Example: `icon = "/usr/share/icons/hicolor/48x48/apps/filezilla.png"`

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


