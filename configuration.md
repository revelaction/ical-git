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

## Images

- **images**: A map of image names to their URL paths. If the key is present in the ATTACH line of the calendar, the URL value here is used. Local filesystem paths are not supported.
  - Example:
    ```toml
    images = { "icon1" = "https://example.com/icon1.png", "icon2" = "https://example.com/icon2.png" }
    ```

## Notifiers

- **notifiers**: A list of notifier types that are enabled.
  - Example: `notifiers = ["desktop"]`

## Fetcher Filesystem

- **directory**: Specifies the directory where the iCal files are stored.
  - Example: `directory = "testdata"`

## Fetcher Git

- **url**: Specifies the URL of the git repository containing the iCal files.
  - Example: `url = "git@mygit-repo.com:me/myrepo.git"`
- **private_key_path**: Specifies the path to the private SSH key used to access the git repository.
  - Example: `private_key_path = "/home/path/to/key"`

## Notifier Telegram

- **token**: The API token for the Telegram bot.
  - Example: `token = "yuu3b3k"`
- **chat_id**: The chat ID to which the notifications will be sent.
  - Example: `chat_id = 588488`

## Notifier Desktop

- **icon**: The path to the icon file used for desktop notifications.
  - Example: `icon = "/usr/share/icons/hicolor/48x48/apps/filezilla.png"`
