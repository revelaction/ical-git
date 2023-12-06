package config

import (
	"errors"
	"time"
)

const (
	Telegram = "telegram"
	Email    = "email"
	Desktop  = "desktop"
)

type Config struct {
	TZ         string `toml:"timezone"`
	DaemonTick string `toml:"tick"`
	Icon       string `toml:"icon"`

	// git repo, credentials
	// or filesystem Path

	Alarms map[string]Alarm
}

type Alarm struct {
	Type string `toml:"type"`
	// ISO8601
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
	{Telegram, "P7D"},  // 1 week
	{Telegram, "P1D"},  // 1 day
	{Desktop, "PT15M"}, // 15 minutes
	{Desktop, "PT45M"},
	{Desktop, "PT13H30M"},
	{Desktop, "PT16H30M"},
	{Desktop, "PT2H30M"},
}
