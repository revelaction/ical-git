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
	config.Alarm
	time time.Time
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

		alarmTime := eventTime.Add(alarm.Dur)
		//slog.Info("🔔 Alarm", "diff", alarm.When, "type", alarm.Type, "alarm_time", alTime)

		if s.isInTickPeriod(alarmTime) {
			tickAlarms = append(tickAlarms, Alarm{Alarm: alarm, time: alarmTime})
		}

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
