package ical

import (
	"slices"
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

func TestParseComment(t *testing.T) {
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

	start := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=DAILY
SUMMARY:Monthly Rent Payment 
DESCRIPTION:Remember to pay the rent.
COMMENT:Stay organized and never miss a payment!
COMMENT:Keep your finances in check with timely reminders.
COMMENT:Effortless rent management starts here.
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
	if notification.Comment == "" {
		t.Errorf("Expected a non-empty Comment, got '%s'", notification.Comment)
	}
}

func TestParseCategories(t *testing.T) {
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

	start := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=DAILY
SUMMARY:Event with Categories
DESCRIPTION:Event with categories A, B, C, D
CATEGORIES:A
CATEGORIES:B
CATEGORIES:C
CATEGORIES:D
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
	expectedCategories := []string{"A", "B", "C", "D"}
	if !slices.Equal(notification.Categories, expectedCategories) {
		t.Errorf("Expected categories %v, got %v", expectedCategories, notification.Categories)
	}
}

func TestParseCategoriesWithNoDate(t *testing.T) {
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

	start := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=DAILY
SUMMARY:Event with Categories
DESCRIPTION:Event with categories A, B, show-no-date
CATEGORIES:A
CATEGORIES:B
CATEGORIES:loose
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
	expectedCategories := []string{"A", "B"}
	if !slices.Equal(notification.Categories, expectedCategories) {
		t.Errorf("Expected categories %v, got %v", expectedCategories, notification.Categories)
	}
}
func TestParseBase64Image(t *testing.T) {
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

	start := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=DAILY
SUMMARY:Event with Base64 Image
DESCRIPTION:Event with a base64 encoded image
ATTACH:data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==
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
	expectedImageData := []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 8, 6, 0, 0, 0, 31, 21, 196, 137, 0, 0, 0, 10, 73, 68, 65, 84, 120, 156, 99, 96, 0, 0, 0, 2, 0, 1, 244, 113, 100, 166, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
	if !bytes.Equal(notification.ImageData, expectedImageData) {
		t.Errorf("Expected image data %v, got %v", expectedImageData, notification.ImageData)
	}
}

func TestParseCategoriesWithNoAlarm(t *testing.T) {
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

	start := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	parser := NewParser(conf, start)

	// Test data
	icalData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=DAILY
SUMMARY:Event with Categories
DESCRIPTION:Event with categories A, B, show-no-alarm
CATEGORIES:A
CATEGORIES:B
CATEGORIES:show-alarm
CATEGORIES:loose
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
	expectedCategories := []string{"A", "B"}
	if !slices.Equal(notification.Categories, expectedCategories) {
		t.Errorf("Expected categories %v, got %v", expectedCategories, notification.Categories)
	}
}
