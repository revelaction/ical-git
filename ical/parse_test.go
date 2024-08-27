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
	{type = "desktop", when = "-P5D"}, 
	{type = "desktop", when = "-P6D"}, 
	{type = "telegram", when = "-P4D"}, 
	{type = "desktop", when = "-P2DT22H49M"}, 
	{type = "desktop", when = "-P3D"}, 
]

#notifiers = ["telegram", "desktop"]
notifiers = ["desktop"]

[fetcher_filesystem]
directory = "testdata"

[notifier_telegram]
token = "yuu3b3k"
chat_id = 588488

[notifier_desktop]
icon = "/usr/share/icons/hicolor/48x48/apps/filezilla.png"
`)
	conf, err := config.Load(configData)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
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
