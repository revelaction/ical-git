package ical

import (
	"github.com/arran4/golang-ical"
)

func getEventAttachments(event *ics.VEvent) []string {
	var result = []string{}
	attachments := event.GetProperties(ics.ComponentPropertyAttach)
	if len(attachments) == 0 {
		return result
	}

	for _, attachment := range attachments {
		result = append(result, attachment.Value)
	}

	return result
}
