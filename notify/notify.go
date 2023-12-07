package notify

import (
	"time"
)

type Notifier interface {
	Notify(n Notification) error
}

type Notification struct {
	EventTime   time.Time
	TimeZone    *time.Location
	Summary     string
	Description string
	Type        string
	Time        time.Time
}

