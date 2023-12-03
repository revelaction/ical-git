package ical

import (
	"bytes"
	//"fmt"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/notify"
)

type Parser struct{
    notifications []notify.Notification
}

func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return err
	}

    for _, event := range cal.Events()  {

        if isEventRecurrent(event) {
            // calculate next
            //fmt.Printf("is recurrent %v#", event)
        }

        //fmt.Printf("is NOT recurrent %v#", event)
        p.notifications = append(p.notifications, buildNotification(event))
        // simple
    }

	return nil
}

func (p *Parser) Notifications() []notify.Notification {
	return p.notifications
}

func buildNotification(event *ics.VEvent) notify.Notification {

    eventTime, _ := event.GetStartAt()
    // TODO check
	summary := event.GetProperty(ics.ComponentPropertySummary).Value

    return notify.Notification{
        EventTime: eventTime,
        Summary: summary,
    }
}

func isEventRecurrent(event *ics.VEvent) bool {

	rule := event.GetProperty(ics.ComponentPropertyRrule)
    if rule == nil {
        return false
    }

    return true
}

func NewParser() *Parser {
	return &Parser{
        notifications: []notify.Notification{},
    }
}
