package ical

import (
	"bytes"
	"github.com/arran4/golang-ical"
	"testing"
)

func TestGetEventAttachments(t *testing.T) {
	// Create an iCal literal with a VEVENT containing attachments
	icalContent := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20240902T100000Z
DTEND:20240902T110000Z
SUMMARY:Test Event
ATTACH:http://example.com/attachment1.pdf
ATTACH:http://example.com/attachment2.pdf
END:VEVENT
END:VCALENDAR`

	// Parse the iCal literal
	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	// Get the first event
	event := cal.Events()[0]

	// Call getEventAttachments with the parsed event
	result := getEventAttachments(event)

	// Verify the result struct
	expected := []string{
		"http://example.com/attachment1.pdf",
		"http://example.com/attachment2.pdf",
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d attachments, but got %d", len(expected), len(result))
	}

	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Expected attachment %s, but got %s", exp, result[i])
		}
	}
}
