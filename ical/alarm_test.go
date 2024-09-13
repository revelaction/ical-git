package ical

import (
	"bytes"
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/alarm"
	"testing"
)

func TestGetEventAlarms(t *testing.T) {
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

	t.Logf("notifications: %#v", event)
	// Call getEventAlarm with the parsed event
	result := getEventAlarms(event, []string{"desktop"})

	// Verify the result struct
	expected := []alarm.Alarm{
		{
			Action:     "desktop",
			DurIso8601: "-PT10M",
		},
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d alarms, but got %d", len(expected), len(result))
	}

	for i, exp := range expected {
		if result[i].Action != exp.Action {
			t.Errorf("Expected Action %s, but got %s", exp.Action, result[i].Action)
		}

		if result[i].DurIso8601 != exp.DurIso8601 {
			t.Errorf("Expected DurIso8601 %s, but got %s", exp.DurIso8601, result[i].DurIso8601)
		}
	}
}

func TestGetEventWithoutBeginVCALENDAR(t *testing.T) {
	// Create an iCal literal without the BEGIN:VCALENDAR tag
	icalContent := `VERSION:2.0
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
	_, err := ics.ParseCalendar(reader)

	// Verify that an error is returned
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestGetEventWithoutEndVCALENDAR(t *testing.T) {
	// Create an iCal literal with a VEVENT tag but without the END:VCALENDAR tag
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
END:VEVENT`

	// Parse the iCal literal
	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	// Verify that the calendar is nil
	if cal != nil {
		t.Errorf("Expected calendar to be nil, but got %#v", cal)
	}
}

func TestParseInvalidICalContent(t *testing.T) {
	// Create a string that does not resemble a proper iCal file
	invalidICalContent := `INVALID CONTENT`

	// Parse the invalid iCal content
	reader := bytes.NewReader([]byte(invalidICalContent))
	_, err := ics.ParseCalendar(reader)

	// Verify that an error is returned
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestGetEventWithoutEndVEVENT(t *testing.T) {
	// Create an iCal literal with a VEVENT tag but without the END:VEVENT tag
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
END:VCALENDAR`

	// Parse the iCal literal
	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	// Get the first event
	event := cal.Events()[0]

	// Verify that the event is nil
	if event != nil {
		t.Errorf("Expected event to be nil, but got %#v", event)
	}
}

func TestGetEventAlarmsWithoutEndingVALARM(t *testing.T) {
	// Create an iCal literal with a VEVENT containing an alarm but without the ending VALARM tag
	icalContent := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20240902T100000Z
DTEND:20240902T110000Z
SUMMARY:Test Event
BEGIN:VALARM
ACTION:DISPLAY
TRIGGER:-PT10M
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

	// Verify that the event is nil
	if event != nil {
		t.Errorf("Expected event to be nil, but got %#v", event)
	}
}
