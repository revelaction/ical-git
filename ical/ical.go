package ical

import (
	"bytes"
	"fmt"
	"time"
	"strings"
    "bufio"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/notify"
	"github.com/teambition/rrule-go"

)

type Parser struct {
	notifications []notify.Notification
}

func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return err
	}

	for _, event := range cal.Events() {

        rruleLines := parseRRule(event)
        fmt.Println("rrule lines", rruleLines)

        fmt.Println("The first recurrence after now is: ", nextEventTime(rruleLines))

		//p.notifications = append(p.notifications, buildNotification(event))
	}

	return nil
}

func (p *Parser) Notifications() []notify.Notification {
	return p.notifications
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
    if nil != summaryProp {
        n.Description = descriptionProp.Value
    }

	return n 
}

func NewParser() *Parser {
	return &Parser{
		notifications: []notify.Notification{},
	}
}

func parseRRule(event *ics.VEvent) string {

    //fmt.Printf("event: %v\n", event.Serialize())
    scanner := bufio.NewScanner(strings.NewReader(event.Serialize()))
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "DTSTART") {
			lines = append([]string{line}, lines...)
			continue
		}

		if strings.HasPrefix(line, "RRULE") {
			lines = append(lines, line)
			return strings.Join(lines, "\n")
		}
	}

	if err := scanner.Err(); err != nil {
        // TODO
	}

	return strings.Join(lines, "\n")
}

func nextEventTime(rruleLines string) time.Time {
    s, err := rrule.StrToRRuleSet(rruleLines)
    if err != nil {
        fmt.Printf("rrlue could not find next: %s\n", err)
        return time.Time{} 
    }
    next := s.After(time.Now(), false)

    // if no RRULE, After provides Zero time get 
    // Try DSTART > Now
    if next.IsZero(){
        // TODO location
        if s.GetDTStart().After(time.Now()){
            return s.GetDTStart()
        }
    }

    return next
}
