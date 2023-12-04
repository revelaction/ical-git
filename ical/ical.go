package ical

import (
	"bytes"
	"fmt"
	"time"
	"strings"
    "bufio"
	"github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
    "github.com/sosodev/duration"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/config"

)

type Parser struct {
	notifications []notify.Notification
    conf config.Config
}


func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return err
	}

	for _, event := range cal.Events() {

        rruleLines := parseRRule(event)
        fmt.Printf("-------------------------rrule: %v\n", rruleLines)
        eventTime, err := nextEventTime(rruleLines)
        if err != nil {
            fmt.Println("error:", err)
            continue
        }

        for _, alarm := range config.DefaultAlarms {
            alarmTime, _ := calculateAlarmTime(eventTime, alarm.Duration)
            if err != nil {
                fmt.Println("error:", err)
                continue
            }

            //n := buildNotification(event) //debug
            fmt.Printf("📅%s duration %s ⏰%s \n\n", eventTime, alarm.Duration, alarmTime)

            tickDuration, _ := time.ParseDuration(p.conf.DaemonTick)

            if isInTickPeriod(alarmTime, tickDuration) {
                n := buildNotification(event)
                n.Time = alarmTime
                n.EventTime = eventTime
            }

            //calculate alarm time 
            // if alarm in tick, (apply offset -3), build Notification
        }

        //fmt.Println("The first recurrence after now is: ", nextEventTime(rruleLines))
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
    if nil != summaryProp {
        n.Description = descriptionProp.Value
    }

	return n 
}

func NewParser() *Parser {
	return &Parser{
		notifications: []notify.Notification{},
        // TODO toml
        conf: config.Config{
            TZ: "Europe/Paris",
            DaemonTick: "3s",
        },
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

func nextEventTime(rruleLines string) (time.Time, error) {
    s, err := rrule.StrToRRuleSet(rruleLines)
    if err != nil {
        // TODO maybe try to aplly config tz.
        return time.Time{}, err
    }

    next := s.After(time.Now(), false)

    // if no RRULE, After provides Zero time get 
    // Try DSTART > Now
    if next.IsZero(){
        // TODO location
        if s.GetDTStart().After(time.Now()){
            return s.GetDTStart(), nil
        }
    }


    return next, nil
}
