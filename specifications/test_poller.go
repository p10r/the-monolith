package specifications

import "gopkg.in/telebot.v3"

type testPoller struct {
	updates chan telebot.Update
	done    chan struct{}
}

func newTestPoller() *testPoller {
	return &testPoller{
		updates: make(chan telebot.Update, 1),
		done:    make(chan struct{}, 1),
	}
}

func (p *testPoller) Poll(b *telebot.Bot, updates chan telebot.Update, stop chan struct{}) {
	for {
		select {
		case upd := <-p.updates:
			updates <- upd
		case <-stop:
			return
		default: //nolint
		}
	}
}
