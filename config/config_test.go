package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)


func TestAlarmsAllowedTelegram(t *testing.T) {
	const testToml = `
notifiers = ["telegram"]
alarms = [
	{type = "telegram", when = "-P1D"},  
	{type = "telegram", when = "-PT15M"},  
]
`
	var conf Config
	if _, err := toml.Decode(testToml, &conf); err != nil {
		t.Fatalf("Failed to decode TOML: %v", err)
	}

	expectedAlarms := []Alarm{
		{Type: "telegram", When: "-P1D"},
		{Type: "telegram", When: "-PT15M"},
	}

	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}
}
