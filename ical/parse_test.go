package ical

import (
	"testing"
	"time"

	"github.com/revelaction/ical-git/config"
)

func TestParser_Parse(t *testing.T) {
	// Setup
	conf := config.Config{
		Location: config.Location{
			Location: time.UTC,
		},
	}
	start := time.Now()
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:Simple Event
DTSTART:20231211T100000
DTEND:20231211T103000
END:VEVENT
END:VCALENDAR
	`)

	// Parse the iCal data
	err := parser.Parse(icalData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the notifications
	notifications := parser.Notifications()
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	notification := notifications[0]
	if notification.Summary != "Simple Event" {
		t.Errorf("Expected summary 'Simple Event', got '%s'", notification.Summary)
	}
}
