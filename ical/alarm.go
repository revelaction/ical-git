package ical

import (
	"github.com/revelaction/ical-git/alarm"
	"github.com/arran4/golang-ical"
)

func getEventAlarm(event *ics.VEvent) alarm.Alarm {
	// Retrieve the alarms from the event
	alarms := event.Alarms()
	if len(alarms) == 0 {
		return alarm.Alarm{}
	}

	a := alarms[0]

	triggerProp := a.GetProperty(ics.ComponentPropertyTrigger)
	actionProp := a.GetProperty(ics.ComponentPropertyAction)

	return alarm.Alarm{
		Action:     actionProp.Value,
		DurIso8601: triggerProp.Value,
	}
}
