package ical

import (
	"testing"
	"time"
)

func TestNextTime(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	expectedTime := time.Date(2024, 4, 1, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}

	if !et.isFloating() {
		t.Errorf("isFloating() = false; want true")
	}

	if !et.hasDtStart() {
		t.Errorf("hasDtStart() = false; want true")
	}

	if !et.hasTzId() {
		t.Errorf("hasTzId() = false; want true")
	}
}

func TestNextTimeInPast(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	if !nextTime.IsZero() {
		t.Errorf("nextTime() = %v; want empty time", nextTime)
	}
}

func TestNextTimeFloating(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART:20240401T000000
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	expectedTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}

	if !et.isFloating() {
		t.Errorf("isFloating() = false; want true")
	}

	if !et.hasDtStart() {
		t.Errorf("hasDtStart() = false; want true")
	}

	if et.hasTzId() {
		t.Errorf("hasTzId() = true; want false")
	}
}
