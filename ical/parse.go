package ical

import (
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"time"
	"log/slog"
)

type Parser struct {
	notifications []notify.Notification
	conf          config.Config
	start         time.Time
}

func NewParser(c config.Config, start time.Time) *Parser {
	return &Parser{
		notifications: []notify.Notification{},
		conf:          c,
		start:         start,
	}
}

func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return fmt.Errorf("calendar parse error: %w", err)
	}

	alarms := NewAlarms(p.conf, p.start)

	for _, event := range cal.Events() {

		et := newEventTime(event)
		et.parse()

		eventTime, err := et.nextTime()
        slog.Info("Next Event time", "event_time", eventTime, "has_rrule", et.hasRRule(), "has_rdate", et.hasRDate(), "is_guessed", et.isGuessed())
		if err != nil {
			if eventTime.IsZero() {
				fmt.Println("error:", err)
				continue
			}
		}

		als := alarms.Get(eventTime)

		for _, alarm := range als {

			n := buildNotification(event)
			n.Time = alarm.Time
			n.EventTime = eventTime
			n.Type = alarm.Type
			p.notifications = append(p.notifications, n)
		}
	}

	return nil
}

func (p *Parser) Notifications() []notify.Notification {
	return p.notifications
}

func buildNotification(event *ics.VEvent) notify.Notification {

	n := notify.Notification{}

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
