package telegram

import (
	"bytes"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"html/template"
	"time"
)

// Telegram
type Telegram struct {
	bot    *tg.BotAPI
	config config.Config
}

func New(conf config.Config) *Telegram {
	bot, err := tg.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		// TODO
		return nil
	}

	return &Telegram{
		bot:    bot,
		config: conf,
	}
}

func (t *Telegram) Notify(n notify.Notification) error {

	message, err := t.renderNotification(n)

	if err != nil {
		return err
	}

	var msg tg.Chattable

	// If we have a image, the text is the caption.
	//"https://image.tmdb.org/t/p/original/1nS8AxnoYE2Y1ANMpVKZnm8iLxP.jpg"
	if n.ImageUrl != "" {
		photo := tg.NewPhoto(t.config.Telegram.ChatId, tg.FileURL(n.ImageUrl))
		photo.Caption = message
		photo.ParseMode = "html"
		msg = photo
	} else {
		text := tg.NewMessage(t.config.Telegram.ChatId, message)
		text.ParseMode = "html"
		msg = text
	}

	_, err = t.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// https://core.telegram.org/bots/api#html-style
func (t *Telegram) renderNotification(n notify.Notification) (string, error) {

	wrapper := struct {
		notify.Notification
		EventTimeZone     string
		EventTimeConf     time.Time
		EventTimeZoneConf string
	}{

		Notification:      n,
		EventTimeZone:     n.EventTimeTz(),
		EventTimeConf:     n.EventTimeConf(t.config.Location.Location),
		EventTimeZoneConf: t.config.Location.Location.String(),
	}

	// Confirmed: ✅, Postponed: 🔄Cancelled: ❌Pending: ⌛Tentative: 🤔Not Attending: 🚫
	tmpl, err := template.New("notification").Parse(notify.Tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, wrapper); err != nil {
		return "", err
	}

	return buf.String(), nil
}
