package alarm

import (
	"testing"
	"time"

)

func TestGetInTick(t *testing.T) {
	// Test that an alarm 15 hours before an event at 20:00 on 1 September 2024
	// is within the tick period starting at 00:00 on 1 September 2024 and lasting 24 hours.
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
	// Test that an alarm 15 hours before an event at 20:00 on 1 September 2024
	// is after the tick period starting at 00:00 on 1 September 2024 and lasting 4 hours.
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
	// Test that an alarm 21 hours before an event at 20:00 on 1 September 2024
	// is before the tick period starting at 00:00 on 1 September 2024 and lasting 12 hours.
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
