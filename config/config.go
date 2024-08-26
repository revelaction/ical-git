package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sosodev/duration"
)

type Config struct {
	Location   Location      `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	Alarms []Alarm

	FetcherFilesystem FetcherFilesystem `toml:"fetcher_filesystem"`

	Notifiers []string
	Telegram  Telegram `toml:"notifier_telegram"`
	Desktop   Desktop  `toml:"notifier_desktop"`
}

type Location struct {
	*time.Location
}

func (l *Location) UnmarshalText(text []byte) error {
	loc, err := time.LoadLocation(string(text))
	if err != nil {
		return err
	}

	l.Location = loc
	return nil
}

type Whend struct {
	time.Duration
}

type Alarm struct {
	Type  string        `toml:"type"`
	When  string        `toml:"when"`
	Duration time.Duration
}

type Alarms []Alarm

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

func LoadConfig(path string) (Config, error) {
	// Config file
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
