package notify

import (
	"time"
)

type Notifier interface {
	Notify(n Notification) error
}

// Notfication represents a notification that should be delivered in the
// current daemon Tick
type Notification struct {
	EventTime   time.Time
	TimeZone    *time.Location
	Summary     string
	Description string
	Type        string
	Diff        string
	Time        time.Time
}
