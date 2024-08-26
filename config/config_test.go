package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)

const testToml = `
timezone = "Europe/Berlin"
tick = "24h"

alarms = [
	{type = "telegram", when = "-P7D"},  
	{type = "desktop", when = "-P1D"},  
	{type = "desktop", when = "-PT15M"},  
	{type = "desktop", when = "-PT1H"},  
	{type = "desktop", when = "-P5D"}, 
	{type = "desktop", when = "-P6D"}, 
	{type = "telegram", when = "-P4D"}, 
	{type = "desktop", when = "-P2DT22H49M"}, 
	{type = "desktop", when = "-P3D"}, 
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

func TestAlarmsAllowed(t *testing.T) {
	var conf Config
	if _, err := toml.Decode(testToml, &conf); err != nil {
		t.Fatalf("Failed to decode TOML: %v", err)
	}

	expectedAlarms := []Alarm{
		{Type: "desktop", When: "-P1D"},
		{Type: "desktop", When: "-PT15M"},
		{Type: "desktop", When: "-PT1H"},
		{Type: "desktop", When: "-P5D"},
		{Type: "desktop", When: "-P6D"},
		{Type: "desktop", When: "-P2DT22H49M"},
		{Type: "desktop", When: "-P3D"},
	}

	actualAlarms := conf.AlarmsAllowed()

	if len(actualAlarms) != len(expectedAlarms) {
		t.Fatalf("Expected %d alarms, got %d", len(expectedAlarms), len(actualAlarms))
	}

	for i, alarm := range actualAlarms {
		if alarm != expectedAlarms[i] {
			t.Errorf("Expected alarm %v, got %v", expectedAlarms[i], alarm)
		}
	}
}
