package ical

import (
	"github.com/arran4/golang-ical"
	"github.com/revelaction/ical-git/alarm"
)

func getEventAlarm(event *ics.VEvent, allowed []string) []alarm.Alarm {
	// Retrieve the alarms from the event
	var result = []alarm.Alarm{}
	alarms := event.Alarms()
	if len(alarms) == 0 {
		return result
	}

	for _, a := range alarms {
		triggerProp := a.GetProperty(ics.ComponentPropertyTrigger)
		//actionProp := a.GetProperty(ics.ComponentPropertyAction)

		parsedDur, err := alarm.ParseIso8601(triggerProp.Value)
		if err != nil {
			// TODO Handle the error, for now, log it and continue
			continue
		}

		for _, allow := range allowed {

			result = append(result, alarm.Alarm{
				Action:     allow,
				DurIso8601: triggerProp.Value,
				Dur:        parsedDur,
				Source:     "event",
			})
		}
	}

	return result
}
