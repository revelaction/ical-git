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
		return fmt.Errorf("calendar parse error: %w", err)
	}

	for _, event := range cal.Events() {

        et := newEventTime(event)
        et.parse()
        fmt.Printf("-------------------------rrule: %v\n", et.joinLines())
        eventTime, err := et.next()
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
    if nil != summaryProp {
        n.Description = descriptionProp.Value
    }

	return n 
}

func NewParser(c config.Config) *Parser {
	return &Parser{
		notifications: []notify.Notification{},
        conf: c,
	}
}

type EventTime struct {
	vEvent        *ics.VEvent
	dStart        string
	rRule         []string
	timeZone      *time.Location
	hasFloating   bool
}

func newEventTime(vEvent *ics.VEvent) *EventTime {
	return &EventTime{
		vEvent:      vEvent,
        rRule: []string{},
	}
}

func (et *EventTime) parse() {
    scanner := bufio.NewScanner(strings.NewReader(et.vEvent.Serialize()))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "DTSTART") {
            et.dStart = line
			continue
		}

		if strings.HasPrefix(line, "RRULE") {
			et.rRule = append(et.rRule, line)
            continue
		}
	}

	if err := scanner.Err(); err != nil {
        // TODO
	}
}

func (et *EventTime) hasRRule() bool {
    if len(et.rRule) > 0 {
        return true
    }

	return false
}

func (et *EventTime) hasDStart() bool {
    if et.dStart == "" {
        return false
    }

	return true
}

// TODO
func (et *EventTime) hasFloatingDStart() bool {
	return false 
}

func (et *EventTime) joinLines() string {

    s := []string{et.dStart}
    s = append(s, et.rRule...)
	return strings.Join(s, "\n")
}

// TODO
// if no RRULE get s.GetDTStart(), if parse error in DSTART, try goland-ical VALUE, check if TZ propeerty, check no UTC in value, apply config timezone if exist or machine localtion
// if no RRULE get s.GetDTStart(), if zero, return -> event is in the past
// if hasRRule, and parse error, return error, 
// if hasRRule, and zero value, return zero value, all events of set are in the past.
func (et *EventTime) next() (time.Time, error) {

    s, err := rrule.StrToRRuleSet(et.joinLines())
    if err != nil {
        return time.Time{}, fmt.Errorf("rrule parse error %s", err)
    }

    if !et.hasRRule() {
        dtStart := s.GetDTStart()

        if dtStart.After(time.Now()){
            return dtStart, nil
        }

        // expired
        return time.Time{}, nil
    }


    return s.After(time.Now(), false), nil
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

// if ruleLines has no RRULE ans DSTART, s.After provides Zero time 
// if ruleLines has RRULE, After provides Zero time when all recurrent events
// are older than DSTART
// 
// custom timezones referencing a VTIMEZONE in the VCALENDAR are not suported
// TODO strct EventTime with hasRRule for better logic
// if no RRULE get s.GetDTStart(), if parse error in DSTART, try goland-ical VALUE, check if TZ propeerty, check no UTC in value, apply config timezone if exist or machine localtion
// if no RRULE get s.GetDTStart(), if zero, return -> event is in the past
// if hasRRule, and parse error, return error, 
// if hasRRule, and zero value, return zero value, all events of set are in the past.
// next()
func nextEventTime(rruleLines string) (time.Time, error) {
    s, err := rrule.StrToRRuleSet(rruleLines)
    if err != nil {
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
