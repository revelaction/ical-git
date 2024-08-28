package telegram

import (
	"bytes"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"html/template"
	"time"
)

type Notifier interface {
	Notify(message string, level string) error
}

type message struct {
	msg       string
	timeStamp time.Time
}

// Telegram
type Telegram struct {
	bot  *tg.BotAPI
	conf config.Config
}

func New(conf config.Config) *Telegram {
	bot, err := tg.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		// TODO
		return nil
	}

	return &Telegram{
		bot:  bot,
		conf: conf,
	}
}

func (t *Telegram) Notify(n notify.Notification) error {

	message, err := renderNotification(n)

	if err != nil {
		return err
	}
	//message = "YOLO"

	msg := tg.NewMessage(t.conf.Telegram.ChatId, message)
	msg.ParseMode = "html"
	_, err = t.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// https://core.telegram.org/bots/api#html-style
func renderNotification(n notify.Notification) (string, error) {
	const tpl = `
<b>{{.Summary}}</b>

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

	var buf bytes.Buffer
	if err := t.Execute(&buf, n); err != nil {
		return "", err
	}

	return buf.String(), nil
}
