package ical

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/sosodev/duration"
	"github.com/teambition/rrule-go"
	"regexp"
	"strings"
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
			alarmTime, err := calculateAlarmTime(eventTime, alarm.Duration)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			// TODO format()
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

type EventTime struct {
	vEvent      *ics.VEvent
	dtStart     string
	rRule       []string
	rDate       []string
	timeZone    *time.Location
	hasFloating bool
	guessed     bool
}

func newEventTime(vEvent *ics.VEvent) *EventTime {
	//validate in config // TODO
	loc, _ := time.LoadLocation("Europe/Berlin")
	return &EventTime{
		vEvent:   vEvent,
		rRule:    []string{},
		rDate:    []string{},
		timeZone: loc,
	}
}

func (et *EventTime) parse() {

	// content line (icalendar spec) should not be longer thant 75 chars.
	// golang-ical properly break lines when serialize()
	// we remove the space to make sure simple scanner works properly
	eventCleaned := strings.Replace(et.vEvent.Serialize(), "\n ", "", -1)
	scanner := bufio.NewScanner(strings.NewReader(eventCleaned))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "DTSTART") {
			et.dtStart = line
			continue
		}

		if strings.HasPrefix(line, "RRULE") {
			et.rRule = append(et.rRule, line)
			continue
		}

		if strings.HasPrefix(line, "RDATE") {
			et.rDate = append(et.rDate, line)
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

func (et *EventTime) hasRDate() bool {
	if len(et.rDate) > 0 {
		return true
	}

	return false
}

func (et *EventTime) hasDtStart() bool {
	if et.dtStart == "" {
		return false
	}

	return true
}

func (et *EventTime) isGuessed() bool {
	return et.guessed
}

// TODO err
func (et *EventTime) isFloating() bool {
	if matched, _ := regexp.MatchString(`\d{8}T\d{6}$`, et.dtStart); matched {
		return true
	}

	return false
}

// hasTzId parses DTSTART;TZID=Some/Timezone:20231129T100000
func (et *EventTime) hasTzId() bool {

	components := strings.Split(et.dtStart, ":")

	if len(components) != 2 {
		return false
	}

	parameters := strings.Split(components[0], ";")

	for _, p := range parameters {
		if strings.HasPrefix(p, "TZID=") {
			return true
		}
	}

	return false
}

func (et *EventTime) parseDtStartInLocation() (time.Time, error) {
	// The layout for an iCalendar floating date-time value
	const layout = "20060102T150405"
	components := strings.Split(et.dtStart, ":")
	dateTime := components[1]

	t, err := time.ParseInLocation(layout, dateTime, et.timeZone)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

//func (et *EventTime) format() string {
//    return fmt.Printf("📅%s duration %s ⏰%s \n\n", eventTime, alarm.Duration, alarmTime)
//}

func (et *EventTime) joinLines() string {

	s := []string{et.dtStart}
	s = append(s, et.rRule...)
	s = append(s, et.rDate...)
	return strings.Join(s, "\n")
}

// golang-ical, rrule  do not support custom timezones
// try to find one
// DTSTART;TZID=<a ref to a VTIMEZONE>:20231129T100000
func (et *EventTime) guessEventTimeForError(err error) (time.Time, error) {
	if !et.hasRRule() && et.hasDtStart() {
		if et.isFloating() && et.hasTzId() {
			guessTime, errParse := et.parseDtStartInLocation()
			if errParse != nil {
				return time.Time{}, fmt.Errorf("error %w: error %w ", err, errParse)
			}

			return guessTime, fmt.Errorf("error %w: guess event time ok", err)
		}
	}

	return time.Time{}, fmt.Errorf("Could not guess event time: %w", err)
}

func (et *EventTime) nextTime() (time.Time, error) {

	now := time.Now()

	s, err := rrule.StrToRRuleSet(et.joinLines())
	if err != nil {
		t, err := et.guessEventTimeForError(err)
		fmt.Println("-----------", t)
		if !t.IsZero() {
			et.guessed = true
			// check if after now
			if t.After(now) {
				return t, err
			}
		}
		return time.Time{}, nil
	}

	if !et.hasRRule() && !et.hasRDate() {
		dtStart := s.GetDTStart()

		if dtStart.After(now) {
			return dtStart, nil
		}

		// expired
		return time.Time{}, nil
	}

	return s.After(now, false), nil
}
