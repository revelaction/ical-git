package alarm

import (
	"time"
	"fmt"
	"github.com/sosodev/duration"
)

type Alarm struct {
	Action     string `toml:"type"`
	DurIso8601 string `toml:"when"`

	Triggertime time.Time     `toml:"-"`
	Dur         time.Duration `toml:"-"`
	// For VLARM maybe
	Description string `toml:"-"`
	// event or config
	Source string `toml:"-"`
}

func (a *Alarm) TriggerTime(eventTime time.Time) time.Time {
	return eventTime.Add(a.Dur)
}

func (a *Alarm) InTickPeriod(eventTime, tickStart time.Time, tick time.Duration) bool {

	t := a.TriggerTime(eventTime)

	if t.Before(tickStart) {
		return false
	}

	if t.After(tickStart.Add(tick)) {
		return false
	}

	return true
}

func ParseIso8601(durIso8601 string) (time.Duration, error) {
	d, err := duration.Parse(durIso8601)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}
	return d.ToTimeDuration(), nil
}
