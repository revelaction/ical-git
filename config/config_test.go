package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)

notifiers = ["desktop"]
alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]
`

func TestAlarmsAllowedDesktop(t *testing.T) {
	const testToml = `
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
`

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}
}
