package ical

import (
	"testing"
	"time"

	"github.com/revelaction/ical-git/config"
)

func TestParse(t *testing.T) {
	// Setup
	configData := []byte(`
timezone = "Europe/Berlin"
tick = "24h"

alarms = [
	{type = "telegram", when = "-P7D"},  
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]

notifiers = ["desktop"]
`)
	conf, err := config.Load(configData)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	start := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

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
	err = parser.Parse(icalData)
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
