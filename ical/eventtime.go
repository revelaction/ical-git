package ical

import (
	"bufio"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
	"regexp"
	"strings"
	"time"
)

// TODO to event
type EventTime struct {
	vEvent  *ics.VEvent
	dtStart string
	rRule   []string
	rDate   []string
	// TODO deprecate this is config
	timeZone *time.Location
	guessed  bool
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

	scanner := bufio.NewScanner(strings.NewReader(et.serialize()))

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

// content line (icalendar spec) should not be longer thant 75 chars.
// golang-ical properly break lines when serialize()
// we remove the space to make sure simple scanner works properly
func (et *EventTime) serialize() string {
	return strings.Replace(et.vEvent.Serialize(), "\r\n ", "", -1)
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

// nextTime returns the next ocurrence of a event
// It can return a zero time indicating that the event is in the past or that
// an error ocurred.
// it tries to guess the time of a event with custom VTIMEZONE. (TODO remove the guessed and return bool gor the guess)
// TODO nextTime(now) and bring the start time
func (et *EventTime) nextTime() (time.Time, error) {

	now := time.Now()

	s, err := rrule.StrToRRuleSet(et.joinLines())
	if err != nil {
		t, err := et.guessEventTimeForError(err)
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
