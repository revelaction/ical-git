package ical

import (
	"github.com/revelaction/ical-git/config"
	"time"
)

type Alarms struct {
	//alarms []Alarm
	conf config.Config

	//start is the of the start of the program or after SIGHUP
	start time.Time
}

type Alarm struct {
	Type string
	Time time.Time
	Diff string
}

func NewAlarms(c config.Config, start time.Time) *Alarms {
	return &Alarms{
		conf:  c,
		start: start,
	}
}

// Get returns the alarms to be trigger in this tick
func (s *Alarms) Get(eventTime time.Time) []Alarm {

	tickAlarms := []Alarm{}
	for _, alarm := range s.conf.AlarmsAllowed() {

		alTime := eventTime.Add(alarm.Duration)
		//slog.Info("🔔 Alarm", "diff", alarm.When, "type", alarm.Type, "alarm_time", alTime)

		if s.isInTickPeriod(alTime) {
			tickAlarms = append(tickAlarms, Alarm{Type: alarm.Type, Time: alTime, Diff: alarm.When})
		}

		// TODO if alarm in tick, (apply offset -3), build Notification
	}

	return tickAlarms
}

func (s *Alarms) isInTickPeriod(t time.Time) bool {

	if t.Before(s.start) {
		return false
	}

	if t.After(s.start.Add(s.conf.DaemonTick)) {
		return false
	}

	return true
}
