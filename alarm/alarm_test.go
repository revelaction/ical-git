package alarm

import (
	"testing"
	"time"

)

func TestGetInTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 24 hours
	// -> tick period [00:00, 23:59]
	// event 1 sept 20:00
	// alarm 15 hours before
	// => alarm 1 sept 5:00, IN tick period
	// Create a literal Alarm instance
	alarm := Alarm{
		Action:     "desktop",
		DurIso8601: "-PT15H",
		Dur:        -15 * time.Hour,
	}

	// Set tick start time to 1 September 2024, 00:00
	tickStart := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	Tick := 24 * time.Hour

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Check if the alarm is in the tick period
	if !alarm.InTickPeriod(eventTime, tickStart, Tick) {
		t.Errorf("Expected alarm to be in tick period")
	}
}

func TestGetAfterTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 4 hours
	// -> tick period [00:00, 03:59]
	// event 1 sept 20:00
	// alarm 15 hours before
	// => alarm 1 sept 5:00, after the current tick period
	// Create a literal Alarm instance
	alarm := Alarm{
		Action:     "desktop",
		DurIso8601: "-PT15H",
		Dur:        -15 * time.Hour,
	}

	// Set tick start time to 1 September 2024, 00:00
	tickStart := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	Tick := 4 * time.Hour

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Check if the alarm is in the tick period
	if alarm.InTickPeriod(eventTime, tickStart, Tick) {
		t.Errorf("Expected alarm to be after tick period")
	}
}

func TestGetBeforeTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 12 hours
	// -> tick period [00:00, 11:59:59]
	// event 1 sept 20:00
	// alarm 21 hours before
	// => alarm 31 august 23:00, before the current tick period
	// Create a literal Alarm instance
	alarm := Alarm{
		Action:     "desktop",
		DurIso8601: "-PT21H",
		Dur:        -21 * time.Hour,
	}

	// Set tick start time to 1 September 2024, 00:00
	tickStart := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	Tick := 12 * time.Hour

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Check if the alarm is in the tick period
	if alarm.InTickPeriod(eventTime, tickStart, Tick) {
		t.Errorf("Expected alarm to be before tick period")
	}
}
