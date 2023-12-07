package config

import (
	"errors"
	"time"
)

type Config struct {
	TZ         string        `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	alarms []Alarm

	Notifiers []string
	Telegram  Telegram
	Desktop   Desktop
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

var errConfNotDuration = errors.New("the value given can not be parsed to a Duration")

func (c *Config) AlarmsAllowed() []Alarm {
	als := []Alarm{}
	for _, al := range c.alarms {
		for _, n := range c.Notifiers {
			if n == al.Type {
				als = append(als, al)
			}
		}
	}

	return als
}

type Alarms []Alarm
