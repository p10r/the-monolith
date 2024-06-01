package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log/slog"
)

// Sender is a passthrough interface to be able to mock outgoing calls
// TODO maybe this can be somehow caught by telebot
type Sender interface {
	Send(c telebot.Context, msg string) error
}

//goland:noinspection GoNameStartsWithPackageName
type TelegramSender struct {
	log *slog.Logger
}

func NewTelegramSender(log *slog.Logger) *TelegramSender {
	l := log.With(slog.String("adapter", "telegram_out"))
	return &TelegramSender{l}
}

func (s TelegramSender) Send(c telebot.Context, msg string) error {
	s.log.Info(fmt.Sprintf("Sending msg to %v", c.Sender().ID))
	return c.Send(msg)
}
