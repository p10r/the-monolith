package specifications

import (
	"github.com/p10r/pedro/pedro/telegram"
	"github.com/p10r/pedro/pkg/sqlite"
	"gopkg.in/telebot.v3"
	"log"
	"testing"
)

func TestE2E(t *testing.T) {
	t.Skip() //TODO
	poller := newTestPoller()

	msg := &telebot.Message{
		ID:      1,
		Payload: "/follow https://ra.co/dj/crilletamalt",
	}

	conn := sqlite.NewDB(":memory:")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewPedro(conn, "botToken", []int64{1})
	bot.Poller = poller
	bot.Start()

	poller.updates <- telebot.Update{ID: 1, Message: msg}
}
