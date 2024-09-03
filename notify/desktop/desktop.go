package desktop

import (
	"bytes"
	"github.com/gen2brain/beeep"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"html/template"
	"time"
)

// Desktop implements the notify.Notifier interface
type Desktop struct {
	config config.Config
}

func New(conf config.Config) *Desktop {
	return &Desktop{
		config: conf,
	}
}

// icon https://specifications.freedesktop.org/icon-theme-spec/icon-theme-spec-latest.html#directory_layout
func (d *Desktop) Notify(n notify.Notification) error {

	body, err := d.renderNotification(n)
	if err != nil {
		return err
	}

	beeep.Notify(n.Summary, body, d.config.Desktop.Icon)
	return nil
}

func (d *Desktop) renderNotification(n notify.Notification) (string, error) {

	// The Summary is already in the header for desktop notifications
	n.Summary = ""

	wrapper := struct {
		notify.Notification
		EventTimeZone     string
		EventTimeConf     time.Time
		EventTimeZoneConf string
	}{

		Notification:      n,
		EventTimeZone:     n.EventTimeTz(),
		EventTimeConf:     n.EventTimeConf(d.config.Location.Location),
		EventTimeZoneConf: d.config.Location.Location.String(),
	}

	// Confirmed: ✅, Postponed: 🔄Cancelled: ❌Pending: ⌛Tentative: 🤔Not Attending: 🚫
	t, err := template.New("notification").Parse(notify.Tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, wrapper); err != nil {
		return "", err
	}

	return buf.String(), nil
}
