package ical

import (
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/alarm"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch"
	"github.com/revelaction/ical-git/notify"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"
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

func (p *Parser) Parse(f fetch.File) error {

	reader := bytes.NewReader(f.Content)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return fmt.Errorf("calendar parse error: %w", err)
	}

	for _, event := range cal.Events() {

		if event == nil {
			slog.Error("📅 Event", "📁", filepath.Base(f.Path), "📌", time.Time{}, "🚨", "unparseable event (nil event)")
			continue
		}

		var alarms []alarm.Alarm = []alarm.Alarm{}

		et := newEventTime(event.Serialize())
		et.parse()

		eventTime, err := et.nextTime(p.start)
		if err != nil {
			slog.Error("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime, "🚨", err)
			continue
		}

		// expired event
		if eventTime.IsZero() {
			slog.Info("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime, "💀️", "expired")
			continue
		}

		in := eventTime.Sub(p.start).Truncate(1 * time.Second)
		slog.Info("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime, "🔖", in)

		// Event Alarms
		for _, a := range getEventAlarms(event, p.conf.NotifierTypes) {
			slog.Info("        : 🔔", "action", a.Action, "durIso", a.DurIso8601, "dur", a.Dur)
			if !a.InTickPeriod(eventTime, p.start, p.conf.DaemonTick) {
				continue
			}

			alarms = append(alarms, a)
		}

		// Only if there are no event alarms, consider config alarms
		if len(alarms) == 0 {
			// Config Alarms
			for _, a := range p.conf.Alarms {

				if !a.HasAllowedAction(p.conf.NotifierTypes) {
					continue
				}

				if !a.InTickPeriod(eventTime, p.start, p.conf.DaemonTick) {
					continue
				}

				alarms = append(alarms, a)
			}
		}

		for _, a := range alarms {

			// To notification
			n := p.buildNotification(event)
			n.Time = a.TriggerTime(eventTime)
			n.EventTime = eventTime
			n.EventPath = f.Path
			n.Type = a.Action
			n.DurIso8601 = a.DurIso8601
			n.Source = a.Source

			p.notifications = append(p.notifications, n)
		}
	}

	return nil
}

func (p *Parser) Notifications() []notify.Notification {
	return p.notifications
}

func buildEventDuration(event *ics.VEvent) time.Duration {
	// duration
	start, err := event.GetStartAt()
	if err != nil {
		return 0
	}

	end, err := event.GetEndAt()
	if err != nil {
		return 0
	}

	return end.Sub(start)
}

func (p *Parser) buildNotification(event *ics.VEvent) notify.Notification {

	n := notify.Notification{}

	n.Duration = buildEventDuration(event)

	summaryProp := event.GetProperty(ics.ComponentPropertySummary)
	if nil != summaryProp {
		n.Summary = summaryProp.Value
	}

	descriptionProp := event.GetProperty(ics.ComponentPropertyDescription)
	if nil != descriptionProp {
		n.Description = descriptionProp.Value
	}

	// we use the ATTACH property only as image file
	// if the image matches the conf, we get the url defined in conf first.
	imageUrlProp := event.GetProperty(ics.ComponentPropertyAttach)
	if nil != imageUrlProp {
		if url, ok := p.conf.Images[imageUrlProp.Value]; ok {
			n.ImageUrl = url
		} else {
			if seemsImageFile(imageUrlProp.Value) {
				n.ImageUrl = imageUrlProp.Value
			}
		}
	}

	locationProp := event.GetProperty(ics.ComponentPropertyLocation)
	if nil != locationProp {
		n.Location = locationProp.Value
	}

	statusProp := event.GetProperty(ics.ComponentPropertyStatus)
	if nil != statusProp {
		n.Status = statusProp.Value
	}

	attendees := event.Attendees()
	if len(attendees) > 0 {
		for _, attendee := range attendees {
			n.Attendees = append(n.Attendees, attendee.Email())
		}
	}

	return n
}

func seemsImageFile(path string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png"}

	ext := strings.ToLower(filepath.Ext(path))

	return slices.Contains(imageExtensions, ext)
}
