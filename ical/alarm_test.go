package ical

import (
	"testing"
	"time"

	"github.com/revelaction/ical-git/config"
)

func TestGetAlarms(t *testing.T) {
	// Define a TOML configuration with one alarm 15 hours ago
	tomlConfig := []byte(`
	timezone = "UTC"
	tick = "1h"

	alarms = [
		{type = "desktop", when = "PT15H"},  
	]
	`)

	// Load the configuration
	conf, err := config.Load([]byte(tomlConfig))
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Set start time to 1 September 2024
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)

	// Create Alarms instance
	alarms := NewAlarms(conf, start)

	// Set event time to 1 September 2024, 12:00
	eventTime := time.Date(2024, time.September, 1, 12, 0, 0, 0, time.UTC)

	// Get alarms
	result := alarms.Get(eventTime)

	// Check if the alarm is in the result
	if len(result) != 1 {
		t.Errorf("Expected 1 alarm, got %d", len(result))
	}

	expectedAlarmTime := eventTime.Add(-15 * time.Hour)
	if result[0].Time != expectedAlarmTime {
		t.Errorf("Expected alarm time %v, got %v", expectedAlarmTime, result[0].Time)
	}
}
