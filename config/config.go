package config

import (
    "time"
    "errors"
)


type Type int
const (
	Telegram Type = iota
	Email
	LinuxDesktop
)


type Config struct {
    TZ string `toml:"timezone"`
    DaemonTick string `toml:"tick"`
    Icon string `toml:"icon"`

    // git repo, credentials
    // or filesystem Path

	Alarms     map[string]Alarm 
}

type Alarm struct {
	Type      string `toml:"type"`
	Duration string `toml:"alarm_duration_before"` 
}

var errConfNotDuration = errors.New("the value given can not be parsed to a Duration")

func (c *Config) Validate() error {
	_, err := time.ParseDuration(c.DaemonTick)
	if err != nil {
		return errConfNotDuration
	}

	return nil
}


type Alarms []Alarm

// if not config given
var DefaultAlarms = Alarms{
	{"telegram", "7d"},     // 1 week
	{"telegram", "1d"},     // 1 day
	{"desktop", "15m"},    // 15 minutes
}
