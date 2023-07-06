package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errNoReq    = errors.New("no req")
	errNoAccess = errors.New("no access")
	errNoCount  = errors.New("no count")
	errUnlucky  = errors.New("unlucky")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Default)
	switch err {
	case errNoReq:
		msg.Text = b.cfg.Messages.NoReq
		b.bot.Send(msg)
	case errNoAccess:
		msg.Text = b.cfg.Messages.NoAccess
		b.bot.Send(msg)
	case errNoCount:
		msg.Text = b.cfg.Messages.NoCounter
		b.bot.Send(msg)
	case errUnlucky:
		msg.Text = b.cfg.Messages.Unlucky
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
