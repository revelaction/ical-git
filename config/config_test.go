package config


import (
	"testing"
	"time"
)

func TestTomlParseError(t *testing.T) {
	var badToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
	invalid-toml-content
`)
	_, err := Load(badToml)
	if err == nil {
		t.Fatalf("Expected an error due to bad TOML content, but got none")
	}
}

func TestAlarmsAllowedDesktop(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLen := 3
	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != expectedLen {
		t.Fatalf("Expected %d alarms, got %d", expectedLen, len(actualAlarms))
	}
}

func TestAlarmsAllowedTelegram(t *testing.T) {
	var testToml = []byte(`
notifiers = ["telegram"]
alarms = [
	{type = "telegram", when = "-P1D"},  
	{type = "telegram", when = "-PT15M"},  
	{type = "telegram", when = "-PT1H"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLen := 3
	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != expectedLen {
		t.Fatalf("Expected %d alarms, got %d", expectedLen, len(actualAlarms))
	}

}

func TestAlarmsAllowedNoTelegram(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "telegram", when = "-P1D"},  
	{type = "telegram", when = "-PT15M"},  
	{type = "telegram", when = "-PT1H"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLen := 0
	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != expectedLen {
		t.Fatalf("Expected %d alarms, got %d", expectedLen, len(actualAlarms))
	}

}

func TestAlarmDuration(t *testing.T) {
	var testToml = []byte(`
alarms = [
	{type = "desktop", when = "-P1D"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(conf.Alarms) != 1 {
		t.Fatalf("Expected 1 alarm, got %d", len(conf.Alarms))
	}

	expectedDuration := -24 * time.Hour
	if conf.Alarms[0].Duration != expectedDuration {
		t.Fatalf("Expected duration %v, got %v", expectedDuration, conf.Alarms[0].Duration)
	}
}

func TestAlarmDurationParseError(t *testing.T) {
	var testToml = []byte(`
alarms = [
	{type = "desktop", when = "invalid-duration"},  
]
`)
	_, err := Load(testToml)
	if err == nil {
		t.Fatalf("Expected an error due to invalid duration, but got none")
	}
}
