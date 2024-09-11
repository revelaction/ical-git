package ical

import (
	"github.com/arran4/golang-ical"
)

// as of now golang_ical does not support many ATTACH lines
// we support only one
func getEventAttachment(event *ics.VEvent) string {
	attachment := event.GetProperty(ics.ComponentPropertyAttach)
	if nil == attachment {
		return ""
	}

	return attachment.Value
}
