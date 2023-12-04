package desktop

import (
	"github.com/revelaction/ical-git/notify"
	"github.com/gen2brain/beeep"
    "time"
)

// Desktop implements the notify.Notifier interface
type Desktop struct {
	Icon string
}

func New(icon string) *Desktop {
	return &Desktop{
		Icon: icon,
	}
}

func (d *Desktop) Notify(n notify.Notification) error {

    beeep.Notify(n.Summary + n.EventTime.Format(time.RFC822), n.Description, d.Icon)
	return nil
}
