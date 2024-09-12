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

func TestNumAlarmsDesktop(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "telegram", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedLen := 3
	actualAlarms := conf.Alarms

	if len(actualAlarms) != expectedLen {
		t.Fatalf("Expected %d alarms, got %d", expectedLen, len(actualAlarms))
	}
}

func TestPositiveDurationError(t *testing.T) {
	confData := []byte(`
tick = "1h"
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "P1D"},  
]
	`)
	_, err := Load(confData)
	if err == nil {
		t.Errorf("Expected error for positive duration, got nil")
	}
}

func TestNegativeTickDurationError(t *testing.T) {
	confData := []byte(`
tick = "-1h"
notifiers = ["desktop"]
	`)
	_, err := Load(confData)
	if err == nil {
		t.Errorf("Expected error for negative tick duration, got nil")
	}
}

func TestAlarmDuration(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
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
	if conf.Alarms[0].Dur != expectedDuration {
		t.Fatalf("Expected duration %v, got %v", expectedDuration, conf.Alarms[0].Dur)
	}
}

func TestAlarmDurationParseError(t *testing.T) {
	var testToml = []byte(`
alarms = [
    notifiers = ["desktop"]
	{type = "desktop", when = "invalid-duration"},  
]
`)
	_, err := Load(testToml)
	if err == nil {
		t.Fatalf("Expected an error due to invalid duration, but got none")
	}
}

func TestExistentImagesWithoutValues(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
]

[images]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if conf.Images == nil {
		t.Fatalf("Expected images property to be initialized, but it was nil")
	}

	if len(conf.Images) != 0 {
		t.Fatalf("Expected images property to be empty, but got %d entries", len(conf.Images))
	}
}

func TestEmptyImagesProperty(t *testing.T) {
	var testToml = []byte(`
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
]
`)
	conf, err := Load(testToml)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if conf.Images == nil {
		t.Fatalf("Expected images property to be initialized, but it was nil")
	}

	if len(conf.Images) != 0 {
		t.Fatalf("Expected images property to be empty, but got %d entries", len(conf.Images))
	}
}
