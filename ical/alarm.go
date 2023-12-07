package ical

import (
	"time"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/sosodev/duration"
)

type Alarms struct {
    alarms []Alarm 
    conf config.Config
    
    //start is the of the start of the program
    start time.Time
}

type Alarm struct {
	Type string 
	Time time.Time  
}

func NewAlarms(c config.Config, start time.Time) *Alarms {
	return &Alarms{
		conf: c,
        start: start, 
	}
}

//Get returns the alarms to be trigger in this tick
func (s *Alarms) Get(eventTime time.Time) []Alarm {
    for _, alarm := range s.conf.Alarms {
        alarmTime, err := calculateAlarmTime(eventTime, alarm.When)
        if err != nil {
            fmt.Println("error:", err)
            continue
        }

        // TODO format()
        fmt.Printf("📅%s duration %s ⏰%s \n\n", eventTime, alarm.When, alarmTime)

        tickDuration, _ := time.ParseDuration(s.conf.DaemonTick)

        if isInTickPeriod(alarmTime, tickDuration) {

        }

        // if alarm in tick, (apply offset -3), build Notification
    }
    return nil
}

func calculateAlarmTime(eventTime time.Time, iso8601Duration string) (time.Time, error) {

	d, err := duration.Parse(iso8601Duration)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing duration: %w", err)
	}

	alarmTime := eventTime.Add(d.ToTimeDuration())
	return alarmTime, nil
}

func isInTickPeriod(t time.Time, duration time.Duration) bool {
	now := time.Now()

	if t.Before(now) {
		return false
	}

	if t.After(now.Add(duration)) {
		return false
	}

	return true
}


