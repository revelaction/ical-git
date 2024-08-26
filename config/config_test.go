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
		{Type: "desktop", When: "-P1D", Whend: -24 * time.Hour},
		{Type: "desktop", When: "-PT15M", Whend: -15 * time.Minute},
		{Type: "desktop", When: "-PT1H", Whend: -1 * time.Hour},
		{Type: "desktop", When: "-P5D", Whend: -5 * 24 * time.Hour},
		{Type: "desktop", When: "-P6D", Whend: -6 * 24 * time.Hour},
		{Type: "desktop", When: "-P2DT22H49M", Whend: -2*24*time.Hour - 22*time.Hour - 49*time.Minute},
		{Type: "desktop", When: "-P3D", Whend: -3 * 24 * time.Hour},
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
