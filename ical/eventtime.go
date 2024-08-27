package ical

import (
	"bufio"
	"fmt"
	"github.com/teambition/rrule-go"
	"regexp"
	"strings"
	"time"
)

// TODO to event
type EventTime struct {
	Event   string
	dtStart string
	rRule   []string
	rDate   []string
	guessed bool
}

func newEventTime(event string) *EventTime {
	return &EventTime{
		Event: event,
		rRule: []string{},
		rDate: []string{},
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
// golang-ical properly break lines when serialize() it adds a space after /r/n
// we remove the sequence "\r\n " to make sure our simple scanner works properly
func (et *EventTime) serialize() string {
	return strings.Replace(et.Event, "\r\n ", "", -1)
}

// In iCal format, floating events are represented by date-time values that do
// not include a "Z" suffix (which would indicate UTC) and do not have an
// associated TZID parameter.
func (et *EventTime) isFloating() bool {

	if !et.hasDtStart() {
		return false
	}

	if et.hasZSuffix() {
		return false
	}

	if et.hasTzId() {
		return false
	}

	return true
}

func (et *EventTime) hasZSuffix() bool {
	if matched, _ := regexp.MatchString(`\d{8}T\d{6}Z$`, et.dtStart); matched {
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

// nextTime returns the next Time of a event
// the timezone of the returned nextTime comes from the Event self, not the Config one
// now can be in a diferent tiemzone as nextTime
// It can return a zero time indicating that the event is in the past or that
// an error ocurres when parsing with rrule package.
//   - not existant DTSTART line
//   - a bad DTSTART line or RRULE or RDATE
//   - not parseable TZID in DTSTART
//
// floating DTSTART do not give error
// if the next eventTime is floating there is not error and should be interpreted as local by the caller
// RDATE does semmes to be properly supported by teambition. Custom logic to try to support
func (et *EventTime) nextTime(now time.Time) (time.Time, error) {

	s, err := rrule.StrToRRuleSet(et.joinLines())
	if err != nil {
		return time.Time{}, err
	}

	// If it also has RDATE this also works properly
	if et.hasRRule() {
		// time can be Zero
		return s.After(now, false), nil
	}

	// Seems to be BUG in github.com/teambition/rrule-go.
	// teambition does not consider the DTSTART when RDATE but no RRULE
	if et.hasRDate() {
		// According to the iCalendar specification (RFC 5545), DTSTART is a
		// required property for VEVENT components
		dtStart := s.GetDTStart()

		if dtStart.After(now) {
			return dtStart, nil
		}

		return s.After(now, false), nil
	}

	dtStart := s.GetDTStart()

	if dtStart.After(now) {
		return dtStart, nil
	}

	// expired
	return time.Time{}, nil
}

func (et *EventTime) joinLines() string {

	s := []string{et.dtStart}
	s = append(s, et.rRule...)
	s = append(s, et.rDate...)
	return strings.Join(s, "\n")
}

// golang-ical, rrule packages do not support custom timezones like:
// DTSTART;TZID=<a ref to a VTIMEZONE>:20231129T100000
// try to check map to config TZ Location
func (et *EventTime) guess(loc *time.Location) (time.Time, error) {
	if !et.hasDtStart() {
	    return time.Time{}, fmt.Errorf("Event without parseable DTSTART")
    }

    if !et.hasZSuffix() && et.hasTzId() {
        guessTime, errParse := et.parseDtStartInLocation(loc)
        if errParse != nil {
            return time.Time{}, fmt.Errorf("error: %w", errParse)
        }

        et.guessed = true
        return guessTime, nil
    }

	// no dstart TODO
	return time.Time{}, fmt.Errorf("Could not guess event time")
}

func (et *EventTime) parseDtStartInLocation(loc *time.Location) (time.Time, error) {
	// The layout for an iCalendar floating date-time value
	const layout = "20060102T150405"
	components := strings.Split(et.dtStart, ":")
	dateTime := components[1]

	t, err := time.ParseInLocation(layout, dateTime, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
