package specifications

import (
	"gopkg.in/telebot.v3"
	"pedro-go/telegram"
	"testing"
)

func TestE2E(t *testing.T) {
	t.Skip() //TODO
	poller := newTestPoller()

	msg := &telebot.Message{
		ID:      1,
		Payload: "/follow https://ra.co/dj/crilletamalt",
	}

	bot := telegram.Pedro("", ":memory:", []int64{1})
	bot.Poller = poller
	bot.Start()

	poller.updates <- telebot.Update{ID: 1, Message: msg}
}
