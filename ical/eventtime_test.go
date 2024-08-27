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

	if !et.hasDtStart() {
		t.Errorf("hasDtStart() = false; want true")
	}

	if et.hasZSuffix() {
		t.Errorf("hasZSuffix() = true; want false")
	}

	if !et.hasTzId() {
		t.Errorf("hasTzId() = false; want true")
	}

	if et.isFloating() {
		t.Errorf("isFloating() = true; want false")
	}
}

func TestNextTimeUTC(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;20240401T000000Z
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, _ := time.LoadLocation("UTC")
	expectedTime := time.Date(2024, 4, 1, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}

	if !et.hasDtStart() {
		t.Errorf("hasDtStart() = false; want true")
	}

	if !et.hasZSuffix() {
		t.Errorf("hasZSuffix() = false want true")
	}

	if et.hasTzId() {
		t.Errorf("hasTzId() = true want false")
	}

	if et.isFloating() {
		t.Errorf("isFloating() = true; want false")
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
}

func TestNextTimeUnparsableTZError(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=Invalid/Timezone:20240401T000000
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	_, err := et.nextTime(now)
	t.Logf("Error message is: %s", err)
	if err == nil {
		t.Fatalf("nextTime should have failed due to unparsable timezone")
	}

	if !et.hasDtStart() {
		t.Errorf("hasDtStart() = false; want true")
	}

	if et.hasZSuffix() {
		t.Errorf("hasZSuffix() = true; want false")
	}

	if !et.hasTzId() {
		t.Errorf("hasTzId() = false; want true")
	}

	if et.isFloating() {
		t.Errorf("isFloating() = true; want false")
	}
}

func TestNextTimeDtStartEmptyError(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
SUMMARY:April 1st Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	_, err := et.nextTime(now)
	t.Logf("Error message is: %s", err)
	if err == nil {
		t.Fatalf("nextTime should have failed due to DTSTART empty line")
	}

	if et.hasDtStart() {
		t.Errorf("hasDtStart() = true; want false")
	}

	if et.hasZSuffix() {
		t.Errorf("hasZSuffix() = true; want false")
	}

	if et.hasTzId() {
		t.Errorf("hasTzId() = true; want false")
	}

	if et.isFloating() {
		t.Errorf("isFloating() = true; want false")
	}
}

func TestHasDtStartEmpty(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
SUMMARY:Empty DTSTART Event
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	if et.hasDtStart() {
		t.Errorf("hasDtStart() = true; want false")
	}
}

func TestHasDtStartEmptyStruct(t *testing.T) {
	et := &EventTime{}

	if et.hasDtStart() {
		t.Errorf("hasDtStart() = true; want false")
	}
}

func TestNextTimeRRule(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RRULE:FREQ=MONTHLY;BYMONTHDAY=-6
SUMMARY:Monthly Event on the 6th last day
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2024, 4, 25, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}

func TestNextTimeRRule1YearLater(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RRULE:FREQ=MONTHLY;BYMONTHDAY=-6
SUMMARY:Monthly Event on the 6th last day
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2025, 4, 25, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}

func TestNextTimeBadRRuleError(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RRULE:FREQ=INVALID;BYMONTHDAY=-6
SUMMARY:Event with invalid RRULE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	_, err := et.nextTime(now)
	if err == nil {
		t.Fatalf("nextTime should have failed due to invalid RRULE")
	}
	t.Logf("Error message is: %s", err)
}

func TestNextTimeRDate(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20250401T000000
RDATE;TZID=America/New_York:20260402T000000
RDATE;TZID=America/New_York:20270403T000000
RDATE;TZID=America/New_York:20280404T000000
SUMMARY:Event with multiple RDATE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2025, 4, 1, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}

func TestNextTimeRDateUTC(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART:20250401T000000Z
RDATE:20260402T000000Z
RDATE:20270402T000000Z
RDATE:20280402T000000Z
SUMMARY:Event with multiple RDATE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2025, 4, 1, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}

// github.com/teambition/rrule-go works after the DTSTART time
func TestNextTimeRDateAfterFirst(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RDATE;TZID=America/New_York:20250402T000000
RDATE;TZID=America/New_York:20260403T000000
RDATE;TZID=America/New_York:20270404T000000
SUMMARY:Event with multiple RDATE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

    // After first
	now := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2025, 4, 2, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}

func TestNextTimeBadRDateError(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RDATE:INVALID_DATE
SUMMARY:Event with invalid RDATE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	_, err := et.nextTime(now)
	if err == nil {
		t.Fatalf("nextTime should have failed due to invalid RDATE")
	}
	t.Logf("Error message is: %s", err)
}
func TestNextTimeRRuleRdate(t *testing.T) {
	event := `BEGIN:VEVENT
UID:123456789
DTSTAMP:20240109T090000Z
DTSTART;TZID=America/New_York:20240401T000000
RRULE:FREQ=MONTHLY;BYMONTHDAY=-6
RDATE;TZID=America/New_York:20240501T000000
SUMMARY:Event with RRULE and RDATE
END:VEVENT`

	et := newEventTime(event)
	et.parse()

	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	nextTime, err := et.nextTime(now)
	if err != nil {
		t.Fatalf("nextTime failed: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}
	expectedTime := time.Date(2024, 4, 25, 0, 0, 0, 0, loc)
	if !nextTime.Equal(expectedTime) {
		t.Errorf("nextTime() = %v; want %v", nextTime, expectedTime)
	}
}
