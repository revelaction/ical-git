package config

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/revelaction/ical-git/alarm"
)

const (
	NotifierTypeTelegram = "telegram"
	NotifierTypeDesktop  = "desktop"
)

var validTypes = []string{NotifierTypeTelegram, NotifierTypeDesktop}

type Config struct {
	Location   Location      `toml:"timezone"`
	DaemonTick time.Duration `toml:"tick"`
	Icon       string        `toml:"icon"`

	Alarms []alarm.Alarm

	Images []Image `toml:"images"`

	FetcherGit FetcherGit `toml:"fetcher_git"`

	FetcherFilesystem FetcherFilesystem `toml:"fetcher_filesystem"`

	NotifierTypes []string `toml:"notifiers"`
	Telegram      Telegram `toml:"notifier_telegram"`
	Desktop       Desktop  `toml:"notifier_desktop"`
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

type FetcherGit struct {
	Url            string
	PrivateKeyPath string `toml:"private_key_path"`
}

type Image struct {
	Name string `toml:"name"`
	Uri  string `toml:"uri"`

	Type string `json:"-"`
	Data []byte `json:"-"`
}

func (c *Config) Image(name string) (Image, bool) {
	for _, img := range c.Images {
		if img.Name == name {
			return img, true
		}
	}
	return Image{}, false
}

// Load loads the configuration. Only alarms compatible with the notifiers are
// present in conf.Alarms
func Load(data []byte) (Config, error) {
	var conf Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return Config{}, err
	}

	// initialize location to UTC if not specified
	if conf.Location.Location == nil {
		conf.Location.Location = time.UTC
	}

	if err := validatePositiveDuration(conf.DaemonTick); err != nil {
		return Config{}, fmt.Errorf("invalid duration for tick %w", err)
	}

	if err := conf.validateNotifierTypes(); err != nil {
		return Config{}, err
	}

	// initialize if alarms not present
	if conf.Alarms == nil {
		conf.Alarms = make([]alarm.Alarm, 0)
	}

	for i, a := range conf.Alarms {
		if err := validateNotifierType(a.Action); err != nil {
			return Config{}, fmt.Errorf("invalid alarm action for alarm %d: %w", i, err)
		}

		dur, err := alarm.ParseIso8601(a.DurIso8601)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing duration for alarm %d: %w", i, err)
		}
		if err := validateNegativeDuration(dur); err != nil {
			return Config{}, fmt.Errorf("invalid duration for alarm %d: %w", i, err)
		}
		conf.Alarms[i].Dur = dur
		conf.Alarms[i].Source = "config"
	}

	// initialize if images not present
	if conf.Images == nil {
		conf.Images = make([]Image, 0)
	}

	for i, im := range conf.Images {
		// check if uri is base64
		data, err := decodeBase64URI(im.Uri)
		if err != nil {
			// try external URL
			errUrl := validateUrl(im.Uri)
			if errUrl != nil {
                return Config{}, fmt.Errorf("Image %d error: bad base64 format[%w], bad Url format [%w]", i, err, errUrl)
			}

			conf.Images[i].Data = nil
			// TODO const
			conf.Images[i].Type = "url"

		} else {
			conf.Images[i].Data = data
			conf.Images[i].Type = "base64"
		}

	}

	return conf, nil
}

func validateNegativeDuration(d time.Duration) error {
	if d > 0 {
		return fmt.Errorf("duration must be positive: %s", d)
	}
	return nil
}

func validatePositiveDuration(d time.Duration) error {
	if d < 0 {
		return fmt.Errorf("duration must be negative: %s", d)
	}
	return nil
}

func validateNotifierType(nt string) error {
	for _, vt := range validTypes {
		if nt == vt {
			return nil
		}
	}
	return fmt.Errorf("invalid notifier type: %s", nt)
}

func (c *Config) validateNotifierTypes() error {
	for _, nt := range c.NotifierTypes {
		if err := validateNotifierType(nt); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) IsFetcherGit() bool {
	return c.FetcherGit.PrivateKeyPath != ""
}

func (c *Config) Fetcher() string {
	if c.IsFetcherGit() {
		return "git"
	}
	return "filesystem"
}

func decodeBase64URI(s string) ([]byte, error) {
	var decoder = base64.StdEncoding
	decodedData, err := decoder.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 data: %w", err)
	}

	return decodedData, nil
}

func validateUrl(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s. Only http and https are allowed", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("missing host in URL")
	}

	if strings.ToLower(u.Hostname()) == "localhost" {
		return fmt.Errorf("localhost is not allowed as an external URL")
	}

	return nil
}
