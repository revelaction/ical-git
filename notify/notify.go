package notify

import (
	"time"
)

const Tpl = `
{{- if .Summary}}
<b>{{.Summary}}</b>
{{- end}}
📅 <b>{{.EventTime.Format "Monday, 2006-01-02"}}</b> <b>{{.EventTime.Format "🕒 15:04"}}</b> 🌍 {{.EventTimeZone}}
📅 <i>{{.EventTimeConf.Format "Monday, 2006-01-02"}}</i> <i>{{.EventTimeConf.Format "🕒 15:04"}}</i> 🌍 <i>{{.EventTimeZoneConf}}</i>

{{- if .Duration}}
⏳ Duration: <b>{{.Duration}}</b><br>
{{- end}}
{{- if .Location}}
📌 Location: <b>{{.Location}}</b><br>
{{- end}}
{{- if .Description}}
📝 Description: {{.Description}}<br>
{{- end}}
{{- if .Status}}
🚦 Status: <b>{{.Status}}</b>
{{- end}}
{{- if .Attendees}}
Attendees:
{{- range .Attendees}}
🔸{{.}}
{{- end}}
{{- end}}

Set by {{.Source}} 🔔 with duration {{.DurIso8601}}
`

type Notifier interface {
	Notify(n Notification) error
}

// Notification represents a notification that should be delivered in the
// current daemon Tick
type Notification struct {
	// Time is the time when the notification should be delivered
	// Its timezone (location) is derived from the ical VEVENT itself
	Time time.Time
	// EventTime is the time of the event.
	// Its timezone (location) is derived from the ical VEVENT itself
	EventTime time.Time
	// EventPath is the path of the ics file containing the event
	EventPath string
	// Type specifies the type of the notification. Currently a desktop or
	// telegram notification
	Type string
	// DurIso8601 is the duration between the event time and the notification
	// time in ISO 8601 format
	DurIso8601 string
	// Duration is the length of the event
	Dur time.Duration
	// Summary is a brief summary of the event
	Summary string
	// Description provides more details about the event
	Description string
	// ImageUrl is the Url of a image to add to the notification. (only telegram)
	ImageUrl string
	// Location is the place where the event takes place
	Location string
	// Status indicates the status of the event
	Status string
	// Duration is the length of the event
	Duration time.Duration
	// Attendees lists the people attending the event
	Attendees []string
	// Source indicates the source of the notification, which can be either "event" or "config"
	Source string
}

// EventTimeConf returns the EventTime in the configured location
func (n Notification) EventTimeConf(loc *time.Location) time.Time {
	return n.EventTime.In(loc)
}

// EventTimeTz extracts the location of EventTime as string
func (n Notification) EventTimeTz() string {
	return n.EventTime.Location().String()
}
