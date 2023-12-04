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

type Conf struct {
	Type      string
	Duration  time.Duration
}

type Confs []Conf

type Notification struct {
    EventTime time.Time
    TimeZone      *time.Location
    Summary       string
    Description   string
    Type Type
    TriggerTime    time.Time
}

var configs = Confs{
	{"telegram", 7 * 24 * time.Hour},   // 1 week
	{"telegram", 24 * time.Hour},       // 1 day
	{"desktop", 15 * time.Minute},      // 15 minutes
}
