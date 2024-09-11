package ical

import (
	"bytes"
	"github.com/arran4/golang-ical"
	"testing"
)

func TestGetEventAttachment(t *testing.T) {
	// Create an iCal literal with a VEVENT containing attachments
	icalContent := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20240902T100000Z
DTEND:20240902T110000Z
SUMMARY:Test Event
ATTACH;FMTTYPE=image/jpeg:http://example.com/attachment.jpg
END:VEVENT
END:VCALENDAR`

	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	event := cal.Events()[0]

	result := GetEventAttachment(event)
	t.Logf(": %#v", result)

	expected := "http://example.com/attachment.jpg"

	if result != expected {
		t.Fatalf("Expected %s attachment, but got %s", expected, result)
	}
}

func TestGetEventAttachmentWithMultipleAttachments(t *testing.T) {
	icalContent := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20240902T100000Z
DTEND:20240902T110000Z
SUMMARY:Test Event
ATTACH;FMTTYPE=image/jpeg:http://example.com/attachment1.jpg
ATTACH;FMTTYPE=image/png:http://example.com/attachment2.png
END:VEVENT
END:VCALENDAR`

	reader := bytes.NewReader([]byte(icalContent))
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	event := cal.Events()[0]

	result := GetEventAttachment(event)
	t.Logf(": %#v", result)

    // only first selected
	expected := "http://example.com/attachment1.jpg"

	if result != expected {
		t.Fatalf("Expected %s attachment, but got %s", expected, result)
	}
}
