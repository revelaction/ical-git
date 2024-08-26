package config

import (
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

const testToml = `
timezone = "Europe/Berlin"
tick = "24h"

alarms = [
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
]

notifiers = ["desktop"]

[fetcher_filesystem]
directory = "testdata"

[notifier_telegram]
token = "yuu3b3k"
chat_id = 588488

[notifier_desktop]
icon = "/usr/share/icons/hicolor/48x48/apps/filezilla.png"
`

func TestAlarmsAllowedDesktop(t *testing.T) {
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
