package ical

import (
	"testing"
	"time"

	"github.com/revelaction/ical-git/config"
)

func TestGetInTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 24 hours
	// -> tick period [00:00, 23:59]
	// event 1 sept 20:00
	// alarm 15 hours before
	// => alarm 1 sept 5:00, IN tick period
	tomlConfig := []byte(`
	timezone = "UTC"
	tick = "24h"

	alarms = [
		{type = "desktop", when = "-PT15H"},  
	]

	notifiers = ["desktop"]
	`)

	conf, err := config.Load(tomlConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Set tick start time to 1 September 2024
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)

	alarms := NewAlarms(conf, start)

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Get alarms
	result := alarms.Get(eventTime)

	// Check if the alarm is in the result
	if len(result) != 1 {
		t.Errorf("Expected 1 alarm, got %d", len(result))
	}

	// 1 September 2024, 5:00
	expectedAlarmTime := time.Date(2024, time.September, 1, 5, 0, 0, 0, time.UTC)
	if result[0].Time != expectedAlarmTime {
		t.Errorf("Expected alarm time %v, got %v", expectedAlarmTime, result[0].Time)
	}
}

func TestGetAfterTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 4 hours
	// -> tick period [00:00, 03:59]
	// event 1 sept 20:00
	// alarm 15 hours before
	// => alarm 1 sept 5:00, after the current tick period
	tomlConfig := []byte(`
	timezone = "UTC"
	tick = "4h"

	alarms = [
		{type = "desktop", when = "-PT15H"},  
	]

	notifiers = ["desktop"]
	`)

	conf, err := config.Load(tomlConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Set tick start time to 1 September 2024
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)

	alarms := NewAlarms(conf, start)

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Get alarms
	result := alarms.Get(eventTime)

	// Check if the alarm is in the result
	if len(result) != 0 {
		t.Errorf("Expected 1 alarm, got %d", len(result))
	}
}

func TestGetBeforeTick(t *testing.T) {
	// start 1 sept 00:00
	// tick 12 hours
	// -> tick period [00:00, 11:59:59]
	// event 1 sept 20:00
	// alarm 21 hours before
	// => alarm 31 august 23:00, before the current tick period
	tomlConfig := []byte(`
	timezone = "UTC"
	tick = "4h"

	alarms = [
		{type = "desktop", when = "-PT21H"},  
	]

	notifiers = ["desktop"]
	`)

	conf, err := config.Load(tomlConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Set tick start time to 1 September 2024
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)

	alarms := NewAlarms(conf, start)

	// Set event time to 1 September 2024, 20:00
	eventTime := time.Date(2024, time.September, 1, 20, 0, 0, 0, time.UTC)

	// Get alarms
	result := alarms.Get(eventTime)

	// Check if the alarm is in the result
	if len(result) != 0 {
		t.Errorf("Expected 1 alarm, got %d", len(result))
	}
}
