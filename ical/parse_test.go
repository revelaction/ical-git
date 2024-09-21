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

	notification := notifications[0]
	if notification.Summary != "Simple Event" {
		t.Errorf("Expected summary 'Simple Event', got '%s'", notification.Summary)
	}
}

func TestParseEventAlarmTriggered(t *testing.T) {

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
TRIGGER:-P2D
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

	notification := notifications[0]
	if notification.DurIso8601 != "-PT1H" {
		t.Errorf("Unexpected duration', got '%s'", notification.DurIso8601)
	}

	if notification.Source != "event" {
		t.Errorf("Unexpected Source', got '%s'", notification.Source)
	}
}

func TestParseConfigAlarmTriggered(t *testing.T) {

	configData := []byte(`
timezone = "Europe/Berlin"
tick = "24h"

notifiers = ["desktop"]

alarms = [
	{type = "desktop", when = "-PT1H"},  
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
TRIGGER:-P2D
ACTION:DISPLAY
DESCRIPTION:Reminder 1 day before
END:VALARM
BEGIN:VALARM
TRIGGER:-PT1D
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

	// only event alarms considered, neither will trigger
	notifications := parser.Notifications()
	if len(notifications) != 0 {
		t.Fatalf("Expected 0 notification, got %d", len(notifications))
	}
}

func TestParseConfigAlarmOrEventAlarmTriggered(t *testing.T) {

	configData := []byte(`
timezone = "Europe/Madrid"
tick = "24h"

notifiers = ["desktop"]

alarms = [
	{type = "desktop", when = "-PT2H"},  
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
TRIGGER:-P2D
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

	//  event first
	notification := notifications[0]
	if notification.DurIso8601 != "-PT1H" {
		t.Errorf("Expected duration', got '%s'", notification.DurIso8601)
	}

	if notification.Source != "event" {
		t.Errorf("Unexpected Source', got '%s'", notification.Source)
	}
}

func TestParseTwoEventAlarmTriggered(t *testing.T) {

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
TRIGGER:-PT2H
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
	if len(notifications) != 2 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	notification := notifications[0]
	if notification.DurIso8601 != "-PT2H" {
		t.Errorf("Expected duration', got '%s'", notification.DurIso8601)
	}

	if notification.Source != "event" {
		t.Errorf("Unexpected Source', got '%s'", notification.Source)
	}

	notification = notifications[1]
	if notification.DurIso8601 != "-PT1H" {
		t.Errorf("Expected duration', got '%s'", notification.DurIso8601)
	}

	if notification.Source != "event" {
		t.Errorf("Unexpected Source', got '%s'", notification.Source)
	}
}
