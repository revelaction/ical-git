package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)


func TestAlarmsAllowedDesktop(t *testing.T) {
	const testToml = `
notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]
`
	var conf Config
	if _, err := toml.Decode(testToml, &conf); err != nil {
		t.Fatalf("Failed to decode TOML: %v", err)
	}

	expectedAlarms := []Alarm{
		{Type: "desktop", When: "-P1D"},
		{Type: "desktop", When: "-PT15M"},
		{Type: "desktop", When: "-PT1H"},
	}

	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}
}

func TestAlarmsAllowedTelegram(t *testing.T) {
	const testToml = `
notifiers = ["telegram"]
alarms = [
	{type = "telegram", when = "-P1D"},  
	{type = "telegram", when = "-PT15M"},  
	{type = "telegram", when = "-PT1H"},  
]
`
	var conf Config
	if _, err := toml.Decode(testToml, &conf); err != nil {
		t.Fatalf("Failed to decode TOML: %v", err)
	}

	expectedAlarms := []Alarm{
		{Type: "telegram", When: "-P1D"},
		{Type: "telegram", When: "-PT15M"},
		{Type: "telegram", When: "-PT1H"},
	}

	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}

	for i, alarm := range actualAlarms {
		if alarm != expectedAlarms[i] {
			t.Fatalf("Expected alarm %v, got %v", expectedAlarms[i], alarm)
		}
	}
}
