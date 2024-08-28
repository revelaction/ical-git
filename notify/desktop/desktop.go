package desktop

import (
	"github.com/gen2brain/beeep"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"html/template"
	//"time"
	"bytes"
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

	body, err := renderNotification(n)
	if err != nil {
		return err
	}

	beeep.Notify(n.Summary, body, d.config.Desktop.Icon)
	return nil
}

func renderNotification(n notify.Notification) (string, error) {

	type NotificationWrapper struct {
		notify.Notification
        EventTimeZone *time.Location
	}

	const tpl = `
📅 <b>{{.EventTime.Format "Monday, 2006-01-02"}}</b> <b>{{.EventTime.Format "🕒 15:04"}}</b> 🌍 {{.EventTimeZone}}

{{if .Duration}}
⏳ Duration: <b>{{.Duration}}</b>
{{end}}
{{if .Location}}
📌 Location: <b>{{.Location}}</b>
{{end}}
{{if .Description}}
📝 Description: {{.Description}}
{{end}}
{{if .Status}}
🚦 Status: <b>{{.Status}}</b>
{{end}}
{{if .Attendees}}
Attendees:
{{- range .Attendees}}
🔸{{.}}
{{- end}}
{{end}}
`
	// Confirmed: ✅, Postponed: 🔄Cancelled: ❌Pending: ⌛Tentative: 🤔Not Attending: 🚫
	t, err := template.New("notification").Parse(tpl)
	if err != nil {
		return "", err
	}


	wrapper := NotificationWrapper{Notification: n}

	var buf bytes.Buffer
	if err := t.Execute(&buf, wrapper); err != nil {
		return "", err
	}

	return buf.String(), nil
}
