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

	expectedTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.FixedZone("America/New_York", -4*3600))
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}
