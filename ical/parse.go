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
	"math/rand"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// LooseCategory is the category used to indicate that the event is without a
// fixed date or deadline
const LooseCategory = "loose"
// ShowAlarmCategory is the category used to indicate that the alarm details
// should be displayed in the notification
const ShowAlarmCategory = "show-alarm"

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
		slog.Error("📅 Event", "📁", filepath.Base(f.Path), "📌", time.Time{}, "🚨", err)
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
			slog.Error("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime.Format("2006-01-02 15:04:05 MST"), "🚨", err)
			continue
		}

		// expired event
		if eventTime.IsZero() {
			slog.Info("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime.Format("2006-01-02 15:04:05 MST"), "💀️", "expired")
			continue
		}

		in := eventTime.Sub(p.start).Truncate(1 * time.Second)
		slog.Info("📅 Event", "📁", filepath.Base(f.Path), "📌", eventTime.Format("2006-01-02 15:04:05 MST"), "🔖", in)

		// Event Alarms
		eventAlarms := getEventAlarms(event, p.conf.NotifierTypes)

		for _, a := range eventAlarms {
			slog.Info("        : 🔔", "action", a.Action, "durIso", a.DurIso8601, "dur", a.Dur)
			if !a.InTickPeriod(eventTime, p.start, p.conf.DaemonTick) {
				continue
			}

			alarms = append(alarms, a)
		}

		// Only if there are no DEFINED event alarms (it does not matter if
		// they are not trigered in this tick period), only then consider
		// config alarms
		if len(eventAlarms) == 0 {
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
			n := p.buildNotificationBasic(event)
			n.Time = a.TriggerTime(eventTime)
			n.EventTime = eventTime
			n.EventPath = f.Path
			n.Type = a.Action
			n.DurIso8601 = a.DurIso8601
			n.Source = a.Source

            // we use the COMMENT property as a one of many DESCRIPTION 
	        n = p.buildNotificationCommentCategories(n, event, et)

            // we use the ATTACH property only as image file
            // if the image matches the conf, we get the url or base64 data defined in conf.
            n = p.buildNotificationImage(n, event)

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

func (p *Parser) buildNotificationBasic(event *ics.VEvent) notify.Notification {

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
func (p *Parser) buildNotificationImage(n notify.Notification, event *ics.VEvent) notify.Notification {
	var validImages []config.Image

	for _, prop := range event.Properties {
		if prop.IANAToken != string(ics.ComponentPropertyAttach) {
			continue
		}
		if image, ok := p.conf.Image(prop.Value); ok {
			validImages = append(validImages, image)
		} else {
			data, err := config.DecodeBase64URI(prop.Value)
			if err == nil {
				validImages = append(validImages, config.Image{Data: data, Type: config.ImageTypeBase64})
			} else {
				err := config.ValidateUrl(prop.Value)
				if err == nil {
					if seemsImageFile(prop.Value) {
						validImages = append(validImages, config.Image{Value: prop.Value, Type: config.ImageTypeUrl})
					}
				}
			}
		}
	}

	if len(validImages) == 0 {
		return n
	}

	randomImage := validImages[rand.Intn(len(validImages))]

	// only one of ImageUrl or ImageData is populated
	n.ImageUrl = randomImage.Value
	n.ImageData = randomImage.Data
	n.ImageName = randomImage.Name

	return n
}
// pickModuloProp calculates the modulo of the number of occurrences of an
// event since a fixed start time, based on the event's interval (only if
// RRULE). This is used to sequentially select properties like COMMENT or
// ATTACH, ensuring that we cycle through them in a predictable manner rather
// than randomly.
func pickModuloProp(eventInterval, modulo int, eventTime time.Time) int {
    fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    duration := eventTime.Sub(fixedTime)
    numberOfIntervals := int(duration.Hours() / 24 / float64(eventInterval))
    return numberOfIntervals % modulo
}

func (p *Parser) buildNotificationCommentCategories(n notify.Notification, event *ics.VEvent, et *EventTime) notify.Notification {
	// Collect all Comment properties
	var comments []string
	// Collect all Categories properties
	var categories []string

	for _, p := range event.Properties {
		if p.IANAToken == string(ics.ComponentPropertyComment) {
			comments = append(comments, p.Value)
		}
		if p.IANAToken == string(ics.ComponentPropertyCategories) {
			if p.Value == LooseCategory {
				n.Loose = true
				continue
			}
			if p.Value == ShowAlarmCategory {
				n.ShowAlarm = true
				continue
			}
			categories = append(categories, p.Value)
		}
	}

	// Select one Comment using pickModuloProp
    if len(comments) > 1 {
        if et.interval == 0 {
            commentIndex := rand.Intn(len(comments))
            n.Comment = comments[commentIndex]
        } else {
            commentIndex := pickModuloProp(et.interval, len(comments), n.EventTime)
            n.Comment = comments[commentIndex]
        }
    }


	// Assign collected categories to the Notification property
	n.Categories = categories

	return n
}
