package config

import (
	"errors"
	"time"
)

type Config struct {
	TZ         string        `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	Alarms []Alarm

	FetcherFilesystem FetcherFilesystem `toml:"fetcher_filesystem"`

	Notifiers []string
	Telegram  Telegram `toml:"notifier_telegram"`
	Desktop   Desktop  `toml:"notifier_desktop"`
}

type Alarm struct {
	Type string `toml:"type"`
	// ISO8601 TODO Diff
	When string `toml:"when"`
}

type Telegram struct {
	Token  string
	ChatId int64 `toml:"chat_id"`
}

type Desktop struct {
	Icon string
}

type FetcherFilesystem struct {
	Directory string
}

var errConfNotDuration = errors.New("the value given can not be parsed to a Duration")

func (c *Config) AlarmsAllowed() []Alarm {
	als := []Alarm{}
	for _, al := range c.Alarms {
		for _, n := range c.Notifiers {
			if n == al.Type {
				als = append(als, al)
			}
		}
	}

	return als
}

type Alarms []Alarm
