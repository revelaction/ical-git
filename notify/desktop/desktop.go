package desktop

import (
	"github.com/gen2brain/beeep"
	"github.com/revelaction/ical-git/notify"
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

// icon https://specifications.freedesktop.org/icon-theme-spec/icon-theme-spec-latest.html#directory_layout
func (d *Desktop) Notify(n notify.Notification) error {

	beeep.Notify(n.Summary+n.EventTime.Format(time.RFC822), n.Description, "/usr/share/icons/hicolor/48x48/apps/filezilla.png")
	return nil
}
