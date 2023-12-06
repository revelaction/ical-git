package ical

import (
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/sosodev/duration"
	"time"
)

type Parser struct {
	notifications []notify.Notification
	conf          config.Config
}

func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return fmt.Errorf("calendar parse error: %w", err)
	}

	for _, event := range cal.Events() {

		et := newEventTime(event)
		et.parse()
		fmt.Printf("-------------------------rrule: %v\n", et.joinLines())
		eventTime, err := et.nextTime()
		if err != nil {
			if eventTime.IsZero() {
				fmt.Println("error:", err)
				continue
			}
		}

		for _, alarm := range config.DefaultAlarms {
			alarmTime, err := calculateAlarmTime(eventTime, alarm.When)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			// TODO format()
			fmt.Printf("📅%s duration %s ⏰%s \n\n", eventTime, alarm.When, alarmTime)

			tickDuration, _ := time.ParseDuration(p.conf.DaemonTick)

			if isInTickPeriod(alarmTime, tickDuration) {
				fmt.Println("in tick")
				n := buildNotification(event)
				n.Time = alarmTime
				n.EventTime = eventTime
				n.Type = alarm.Type
				p.notifications = append(p.notifications, n)
			}

			// if alarm in tick, (apply offset -3), build Notification
		}

	}

	return nil
}

func (p *Parser) Notifications() []notify.Notification {
	return p.notifications
}

func calculateAlarmTime(eventTime time.Time, iso8601Duration string) (time.Time, error) {

	d, err := duration.Parse(iso8601Duration)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing duration: %w", err)
	}

	alarmTime := eventTime.Add(-d.ToTimeDuration())
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

func buildNotification(event *ics.VEvent) notify.Notification {

	n := notify.Notification{}
	// TODO
	//n.EventTime, _ := event.GetStartAt()

	summaryProp := event.GetProperty(ics.ComponentPropertySummary)
	if nil != summaryProp {
		n.Summary = summaryProp.Value
	}

	descriptionProp := event.GetProperty(ics.ComponentPropertyDescription)
	if nil != descriptionProp {
		n.Description = descriptionProp.Value
	}

	return n
}

func NewParser(c config.Config) *Parser {
	return &Parser{
		notifications: []notify.Notification{},
		conf:          c,
	}
}

