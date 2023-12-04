package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5" 
	"github.com/revelaction/ical-git/notify"
	"time"
)

const (
	telegran_token = "sadfsag"
	chatId         = int64(373747346)
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
	bot *tg.BotAPI
}

func New() *Telegram {
	bot, err := tg.NewBotAPI(telegran_token)
	if err != nil {
		// TODO
		return nil
	}

	return &Telegram{
		bot: bot,
	}
}

func (t *Telegram) Notify(n notify.Notification) error {
    message := n.Summary + n.EventTime.Format(time.RFC822)
	msg := tg.NewMessage(chatId, message)
	msg.ParseMode = "markdown"
	_, err := t.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
