package ical

import (
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/sosodev/duration"
	"time"
)

type Alarms struct {
	//alarms []Alarm
	conf config.Config

	//start is the of the start of the program
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
		fmt.Printf("alarm %#v\n", alarm)

		alTime, err := alarmTime(eventTime, alarm.When)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		fmt.Printf("📅%s duration %s ⏰%s \n\n", eventTime, alarm.When, alTime)

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

func alarmTime(eventTime time.Time, iso8601Duration string) (time.Time, error) {

	d, err := duration.Parse(iso8601Duration)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing duration: %w", err)
	}

	alarmTime := eventTime.Add(d.ToTimeDuration())
	return alarmTime, nil
}
