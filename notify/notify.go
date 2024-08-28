package notify

import (
	"time"
)

type Notifier interface {
	Notify(n Notification) error
}

// Notification represents a notification that should be delivered in the
// current daemon Tick
type Notification struct {
	// Time is the time when the notification should be delivered
    // Its timezone (location) is derived from the ical VEVENT itself
	Time time.Time
	// EventTime is the time of the event. 
    // Its timezone (location) is derived from the ical VEVENT itself
	EventTime time.Time
    // Type specifies the type of the notification. Currently a desktop or
    // telegram notification
	Type string
    // DurIso8601 is the duration between the event time and the notification
    // time in ISO 8601 format
	DurIso8601 string
	// Duration is the length of the event
	Dur time.Duration
	// Summary is a brief summary of the event
	Summary string
	// Description provides more details about the event
	Description string
	// Location is the place where the event takes place
	Location string
	// Status indicates the status of the event
	Status string
	// Duration is the length of the event
	Duration time.Duration
	// Attendees lists the people attending the event
	Attendees []string
}
