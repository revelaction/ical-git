package config

import (
	"fmt"
	//"slices"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/revelaction/ical-git/alarm"
)

const (
	NotifierTypeTelegram = "telegram"
	NotifierTypeDesktop  = "desktop"
)


type Config struct {
	Location   Location      `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	Alarms []alarm.Alarm

	FetcherFilesystem FetcherFilesystem `toml:"fetcher_filesystem"`

	NotifierTypes []string
	Telegram  Telegram `toml:"notifier_telegram"`
	Desktop   Desktop  `toml:"notifier_desktop"`
}

func (c *Config) validateNotifierTypes() error {
	validTypes := []string{NotifierTypeTelegram, NotifierTypeDesktop}

	for _, nt := range c.NotifierTypes {
		valid := false
		for _, vt := range validTypes {
			if nt == vt {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid notifier type: %s", nt)
		}
	}

	return nil
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

// Load loads the configuration. Only alarms compatible with the notifiers are
// present in conf.Alarms
func Load(data []byte) (Config, error) {
	var conf Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return Config{}, err
	}

	if err := conf.validateNotifierTypes(); err != nil {
		return Config{}, err
	}

	for i, a := range conf.Alarms {

		dur, err := alarm.ParseIso8601(a.DurIso8601)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing duration for alarm %d: %w", i, err)
		}
		conf.Alarms[i].Dur = dur
		conf.Alarms[i].Source = "config"
	}
	return conf, nil
}

