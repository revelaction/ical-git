package config

import (
	"testing"

)


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
