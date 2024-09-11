package ical

import (
	"github.com/arran4/golang-ical"
)

// getEventAttachment retrieves the first attachment from a VEvent.
// Note: The current implementation of golang-ical does not support multiple ATTACH lines,
// so this function only returns the first attachment found.
func getEventAttachment(event *ics.VEvent) string {
	attachment := event.GetProperty(ics.ComponentPropertyAttach)
	if nil == attachment {
		return ""
	}

	return attachment.Value
}
