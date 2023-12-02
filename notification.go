package main

import (
	"time"

	"github.com/gen2brain/beeep"
)

type NotificationType int

const (
    TelegramNotification NotificationType = iota
    EmailNotification
    LinuxDesktopNotification
)

type Notification struct {
    NextEventTime time.Time
    TimeZone      *time.Location
    Summary       string
    Description   string
    NotificationType NotificationType
    TriggerTime    time.Time
    // Add any additional fields you deem necessary
}


func main() {

	err = beeep.Notify("Notification Title", "Notification Message", "./logo.png")
	if err != nil {
		panic(err)
	}

	// Sleep to allow the notification to be displayed
	time.Sleep(5 * time.Second)
}
