package notify

import (
	"time"
)

type Type int

const (
    Telegram Type = iota
    Email
    LinuxDesktop
)

type Notifier interface {
    Notify(n Notification) error
}

type Notification struct {
    EventTime time.Time
    TimeZone      *time.Location
    Summary       string
    Description   string
    Type Type
    TriggerTime    time.Time
}

