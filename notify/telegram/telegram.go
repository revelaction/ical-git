package telegram

import (
	"bytes"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"html/template"
	"log/slog"
	"time"
)

// Telegram
type Telegram struct {
	bot    *tg.BotAPI
	config config.Config
}

func New(conf config.Config) (*Telegram, error) {
	bot, err := tg.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("could not create telegram Bot; %w", err)
	}

	return &Telegram{
		bot:    bot,
		config: conf,
	}, nil
}

func (t *Telegram) Notify(n notify.Notification) error {

	message, err := t.renderNotification(n)

	if err != nil {
		return err
	}

	var msg tg.Chattable

	// If we have a image, the text is the caption.
	if n.ImageUrl != "" {
		photo := tg.NewPhoto(t.config.Telegram.ChatId, tg.FileURL(n.ImageUrl))
		photo.Caption = message
		photo.ParseMode = "html"
		msg = photo
	} else if n.ImageData != nil {
		photo := tg.NewPhoto(t.config.Telegram.ChatId, tg.FileBytes{Name: "image.png", Bytes: n.ImageData})
		photo.Caption = message
		photo.ParseMode = "html"
		msg = photo
	} else {
		text := tg.NewMessage(t.config.Telegram.ChatId, message)
		text.ParseMode = "html"
		msg = text
	}

	m, err = t.bot.Send(msg)
	if err != nil {
		return err
	}

	if len(m.Photo) > 0 {
		slog.Info("Image File Id", "id", m.Photo[len(m.Photo)-1].FileID)
	}
	//slog.Info("Image File Id", "id", m)

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
