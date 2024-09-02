package ical

import (
	"alarm"
	"ics"
)

func getEventAlarm(event *ics.VEvent) alarm.Alarm {
	// Retrieve the alarms from the event
	alarms := event.Alarms()
	if len(alarms) == 0 {
		return alarm.Alarm{}
	}

	// Use the first alarm for simplicity
	alarm := alarms[0]

	// Get the Trigger and Action properties
	triggerProp := alarm.GetProperty(ics.ComponentPropertyTrigger)
	actionProp := alarm.GetProperty(ics.ComponentPropertyAction)

	// Create an alarm.Alarm literal with these properties
	return alarm.Alarm{
		Action:     actionProp.Value,
		DurIso8601: triggerProp.Value,
	}
}
