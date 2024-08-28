package notify

import (
	"time"
)

type Notifier interface {
	Notify(n Notification) error
}

// Notfication represents a notification that should be delivered in the
// current daemon Tick
// Notification represents a notification that should be delivered in the
// current daemon Tick
type Notification struct {
	// EventTime is the time of the event
	EventTime time.Time
	// TimeZone is the location of the event
	TimeZone *time.Location
	// Summary is a brief summary of the event
	Summary string
	// Description provides more details about the event
	Description string
	// Location is the place where the event takes place
	Location string
	// Type specifies the type of the event
	Type string
	// DurIso8601 is the duration of the event in ISO 8601 format
	DurIso8601 string
	// Time is the time when the notification should be delivered
	Time time.Time
	// Status indicates the status of the event
	Status string
	// Duration is the length of the event
	Duration time.Duration
	// Attendees lists the people attending the event
	Attendees []string
	EventTime   time.Time
	TimeZone    *time.Location
	Summary     string
	Description string
	Location    string
	Type        string
	DurIso8601  string
	Time        time.Time
	Status      string
	Duration    time.Duration
	Attendees   []string
}
