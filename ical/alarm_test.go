package ical

import (
	"bytes"
	"testing"
	"alarm"
	"github.com/arran4/golang-ical"
)

func TestGetEventAlarm(t *testing.T) {
	// Create an iCal literal with a VEVENT containing an alarm
	icalContent := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20240902T100000Z
DTEND:20240902T110000Z
SUMMARY:Test Event
BEGIN:VALARM
ACTION:DISPLAY
TRIGGER:-PT10M
END:VALARM
END:VEVENT
END:VCALENDAR`

	// Parse the iCal literal
	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	// Get the first event
	event := cal.Events()[0]

	// Call getEventAlarm with the parsed event
	result := getEventAlarm(event)

	// Verify the result struct
	expected := alarm.Alarm{
		Action:     "DISPLAY",
		DurIso8601: "-PT10M",
	}

	if result.Action != expected.Action {
		t.Errorf("Expected Action %s, but got %s", expected.Action, result.Action)
	}

	if result.DurIso8601 != expected.DurIso8601 {
		t.Errorf("Expected DurIso8601 %s, but got %s", expected.DurIso8601, result.DurIso8601)
	}
}
