package alarm

import (
	"time"
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

// Get returns the alarms to be trigger in this tick
// the timezone of the eventTime comes from the event timezone, or (automatically) local if floating event
// The alarms time returned by Get are also the timezone of the event
//func (s *Alarms) Get(eventTime time.Time) []Alarm {
//
//	tickAlarms := []Alarm{}
//	for _, alarm := range s.conf.AlarmsAllowed() {
//
//		alarmTime := eventTime.Add(alarm.Dur)
//		//slog.Info("🔔 Alarm", "diff", alarm.When, "type", alarm.Type, "alarm_time", alTime)
//
//		if s.isInTickPeriod(alarmTime) {
//			tickAlarms = append(tickAlarms, Alarm{Alarm: alarm, time: alarmTime})
//		}
//
//	}
//
//	return tickAlarms
//}
//
//func (s *Alarms) isInTickPeriod(t time.Time) bool {
//
//	if t.Before(s.start) {
//		return false
//	}
//
//	if t.After(s.start.Add(s.conf.DaemonTick)) {
//		return false
//	}
//
//	return true
//}
