package config

import (
	"fmt"
	"slices"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/revelaction/ical-git/alarm"
	"github.com/sosodev/duration"
)

type Config struct {
	Location   Location      `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	Alarms []alarm.Alarm

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

	conf.Alarms = slices.DeleteFunc(conf.Alarms, func(a alarm.Alarm) bool {
		for _, n := range conf.Notifiers {
			if n == a.Action {
				return false
			}
		}

		return true

	})

	for i, alarm := range conf.Alarms {

		dur, err := parseDurIso8601(alarm.DurIso8601)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing duration for alarm %d: %w", i, err)
		}
		conf.Alarms[i].Dur = dur
	}
	return conf, nil
}

func parseDurIso8601(isoDur string) (time.Duration, error) {
	d, err := duration.Parse(isoDur)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}
	return d.ToTimeDuration(), nil
}
