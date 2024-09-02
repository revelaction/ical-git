package ical

import (
	"github.com/revelaction/ical-git/alarm"
	"github.com/arran4/golang-ical"
)

func getEventAlarm(event *ics.VEvent) []alarm.Alarm {
	// Retrieve the alarms from the event
    var result =[]alarm.Alarm{} 
	alarms := event.Alarms()
	if len(alarms) == 0 {
		return result
	}

	for _, a := range alarms {
		triggerProp := a.GetProperty(ics.ComponentPropertyTrigger)
		actionProp := a.GetProperty(ics.ComponentPropertyAction)

		parsedDur, err := alarm.ParseIso8601()
		if err != nil {
			// Handle the error, for now, log it and continue
			continue
		}

		result = append(result, alarm.Alarm{
			Action:     actionProp.Value,
			DurIso8601: triggerProp.Value,
			Dur:        parsedDur,
			Source:     "event",
		})
	}

	return result
}
