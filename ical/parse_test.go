package ical

import (
	"testing"
	"time"

	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch"
)

func TestParse(t *testing.T) {

	configData := []byte(`
timezone = "Europe/Berlin"
tick = "24h"
notifiers = ["desktop"]

alarms = [
	{type = "desktop", when = "-P1D"},  
]
`)
	conf, err := config.Load(configData)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}


	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:Simple Event
DTSTART:20241201T100000Z
END:VEVENT
END:VCALENDAR
`)

	// Parse the iCal data
	file := fetch.File{
		Path:    "",
		Content: icalData,
		Error:   nil,
	}

	start := time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	err = parser.Parse(file)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the notifications
	notifications := parser.Notifications()
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	t.Logf("notifications: %#v", notifications)
	notification := notifications[0]
	if notification.Summary != "Simple Event" {
		t.Errorf("Expected summary 'Simple Event', got '%s'", notification.Summary)
	}
}

// TestParseEvents tests the parsing of an iCal file containing an event with multiple alarms.
// It sets up a configuration with a specific timezone, tick duration, and notifiers.
// The test data includes an event titled "Event with Alarms" with alarms triggering 1 day and 1 hour before the event.
// The test verifies that the parser correctly generates a notification for this event.
func TestParseEvents(t *testing.T) {
	// Setup
	configData := []byte(`
timezone = "Europe/Berlin"
tick = "24h"

notifiers = ["desktop"]

alarms = [
	{type = "desktop", when = "-P1D"},  
]

`)
	conf, err := config.Load(configData)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	start := time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:Event with Alarms
DTSTART:20241201T100000Z
BEGIN:VALARM
TRIGGER:-P1D
ACTION:DISPLAY
DESCRIPTION:Reminder 1 day before
END:VALARM
BEGIN:VALARM
TRIGGER:-PT1H
ACTION:DISPLAY
DESCRIPTION:Reminder 1 hour before
END:VALARM
END:VEVENT
END:VCALENDAR
`)

	// Parse the iCal data
	file := fetch.File{
		Path:    "",
		Content: icalData,
		Error:   nil,
	}
	err = parser.Parse(file)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the notifications
	notifications := parser.Notifications()
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	t.Logf("notifications: %#v", notifications)
	notification := notifications[0]
	if notification.Summary != "Event with Alarms" {
		t.Errorf("Expected summary 'Event with Alarms', got '%s'", notification.Summary)
	}
}
