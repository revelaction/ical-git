package telegram

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
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
	fmt.Println("token", conf.Telegram.Token)
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
	message := n.Summary + " " + n.EventTime.Format(time.RFC822)
	fmt.Println("in telegram", message, t.conf.Telegram.ChatId)
	msg := tg.NewMessage(t.conf.Telegram.ChatId, message)
	msg.ParseMode = "markdown"
	_, err := t.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
