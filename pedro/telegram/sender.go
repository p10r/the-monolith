package telegram

import "gopkg.in/telebot.v3"

// Sender is a passthrough interface to be able to mock outgoing calls
// TODO maybe this can be somehow caught by telebot
type Sender interface {
	Send(c telebot.Context, msg string) error
}

type TelebotSender struct {
}

func (s TelebotSender) Send(c telebot.Context, msg string) error {
	return c.Send(msg)
}
