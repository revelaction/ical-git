package ical

import (
	"github.com/revelaction/ical-git/alarm"
	"github.com/arran4/golang-ical"
)

func getEventAlarm(event *ics.VEvent) []alarm.Alarm {
	// Retrieve the alarms from the event
	alarms := event.Alarms()
	result := make([]alarm.Alarm, 0, len(alarms))

	for _, a := range alarms {
		triggerProp := a.GetProperty(ics.ComponentPropertyTrigger)
		actionProp := a.GetProperty(ics.ComponentPropertyAction)

		result = append(result, alarm.Alarm{
			Action:     actionProp.Value,
			DurIso8601: triggerProp.Value,
		})
	}

	return result
}
