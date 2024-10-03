package notify

import (
	"time"
)

const Tpl = `
{{- if .Summary}}
<b>{{.Summary}}</b>
<b> </b>
{{- end}}
{{- if and .ShowDate .EventTime}}
📅 <b>{{.EventTime.Format "Monday, 2006-01-02"}}</b> <b>{{.EventTime.Format "🕒 15:04"}}</b> 🌍 {{.EventTimeZone}}
{{- end}}
{{- if and .ShowDate .EventTimeConf}}
📅 <i>{{.EventTimeConf.Format "Monday, 2006-01-02"}}</i> <i>{{.EventTimeConf.Format "🕒 15:04"}}</i> 🌍 <i>{{.EventTimeZoneConf}}</i>
{{- end}}

{{- if .Duration}}
⏳ Duration: <b>{{.Duration}}</b>
{{- end}}
{{- if .Location}}
📌 Location: <b>{{.Location}}</b>
{{- end}}
{{- if .Description}}
{{.Description}}
{{- end}}
{{- if .Comment}}
{{.Comment}}
{{- end}}
{{- if .Categories}}
Categories:
{{- range .Categories}}
🔖{{.}}
{{- end}}
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
{{- if .ShowAlarm}}
{{.Source}} 🔔 {{.DurIso8601}}
{{- end}}
{{- if .IsUrgent}}
📢 🔥 Urgent! Time Difference: {{.TimeDifference}} 🔥 📢
{{- end}}
`

const UrgencyThreshold = 1 * time.Hour

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
	// ImageName is the (file) Name of a image to add to the notification. (only telegram)
	ImageName string
	// ImageUrl is the Url of a image to add to the notification. (only telegram)
	ImageUrl string
	// ImageData is the image binary data to add to the notification. (only telegram)
	ImageData []byte
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
	// Comment provides additional text about the notification.
	// The ical standard allows many COMMENT lines. We use this feature to provide "one of many" messages:
	// a ics file for recurrent motivational reminder can have many comment lines and every instance of
	// the message will have a different related comment.
	Comment string
	// Categories lists the categories associated with the event
	Categories []string
	// ShowDate indicates whether to display the event dates in the notification
	ShowDate bool
	// ShowAlarm indicates whether to display the alarm details in the notification
	ShowAlarm bool
}

// EventTimeConf returns the EventTime in the configured location
func (n Notification) EventTimeConf(loc *time.Location) time.Time {
	return n.EventTime.In(loc)
}

// EventTimeTz extracts the location of EventTime as string
func (n Notification) EventTimeTz() string {
	return n.EventTime.Location().String()
}

// IsUrgent checks if the notification is urgent based on the difference between EventTime and Time
func (n Notification) IsUrgent() bool {
	return n.EventTime.Sub(n.Time) > UrgencyThreshold
}

// TimeDifference returns the difference between EventTime and Time
func (n Notification) TimeDifference() time.Duration {
	return n.EventTime.Sub(n.Time)
}
