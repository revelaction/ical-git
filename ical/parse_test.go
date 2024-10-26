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
ATTACH;ENCODING=BASE64;FMTTYPE=image/jpeg:iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAACklEQVR4nGMAAQAABQABDQottAAAAABJRU5ErkJggg==
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
	if notification.ImageData == nil {
		t.Errorf("Expected non-nil ImageData, got nil")
	}
}

func TestParseTwoValidAttachLines(t *testing.T) {
	configData := []byte(`
timezone = "Europe/Berlin"
tick = "24h"

notifiers = ["desktop"]

alarms = [
	{type = "desktop", when = "-P1D"},  
]

images = [
	{name = "testImage", value = "http://example.com/image.jpg"}
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
SUMMARY:Event with Two Valid Attach Lines
DESCRIPTION:Event with one base64 encoded image and one URL image
ATTACH;ENCODING=BASE64;FMTTYPE=image/jpeg:iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAACklEQVR4nGMAAQAABQABDQottAAAAABJRU5ErkJggg==
ATTACH:testImage
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
	if notification.ImageData == nil && notification.ImageUrl == "" {
		t.Errorf("Both ImageData and ImageUrl should not be empty")
	}
}

func TestPickModuloProp(t *testing.T) {
	testCases := []struct {
		eventInterval int
		modulo        int
		eventTime     time.Time
		expected      int
	}{
		{
			eventInterval: 1,
			modulo:        10,
			eventTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:      1461 % 10,
		},
		{
			eventInterval: 7,
			modulo:        5,
			eventTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:      (1461 / 7) % 5,
		},
		{
			eventInterval: 30,
			modulo:        3,
			eventTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:      (1461 / 30) % 3,
		},
		{
			eventInterval: 1,
			modulo:        1,
			eventTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:      0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("eventInterval=%d, modulo=%d, eventTime=%s", tc.eventInterval, tc.modulo, tc.eventTime), func(t *testing.T) {
			result := pickModuloProp(tc.eventInterval, tc.modulo, tc.eventTime)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}
