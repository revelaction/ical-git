package config

import (
	"errors"
	"time"
)

type Config struct {
	TZ         string `toml:"timezone"`
	DaemonTick string `toml:"tick"`
	Icon       string `toml:"icon"`

	Alarms []Alarm

    Notifiers []string
    Telegram Telegram
    Desktop Desktop
}

type Alarm struct {
	Type string `toml:"type"`
	// ISO8601
	When string `toml:"when"`
}

type Telegram struct {
	Token string 
	ChatId int64 `toml:"chat_id"`

}

type Desktop struct {
    Icon string
}

var errConfNotDuration = errors.New("the value given can not be parsed to a Duration")

func (c *Config) Validate() error {
	_, err := time.ParseDuration(c.DaemonTick)
	if err != nil {
		return errConfNotDuration
	}

	return nil
}

func (c *Config) IsNotifierAllowed(notifier string) bool {
    for _, n := range c.Notifiers {
        if n == notifier {
            return true
        }
    }

    return false
}

type Alarms []Alarm

