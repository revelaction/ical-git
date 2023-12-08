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

    body, err:= renderNotification(n)
    if err != nil {
        return err
    }

	beeep.Notify(n.Summary, body, d.config.Desktop.Icon)
	return nil
}


func renderNotification(n notify.Notification) (string, error) {
    const tpl = `
📅 {{.EventTime.Format "Monday, 2006-02-01"}} {{.EventTime.Format "🕒 15:04"}} 
🌍 Time Zone: {{.TimeZone}}
📢 Summary: <b>{{.Summary}}</b>
📝 Description: {{.Description}}
`
    t, err := template.New("notification").Parse(tpl)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := t.Execute(&buf, n); err != nil {
        return "", err
    }

    return buf.String(), nil
}

