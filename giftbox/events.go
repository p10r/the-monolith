package giftbox

import (
	"fmt"
	slogtelegram "github.com/samber/slog-telegram/v2"
	"log/slog"
)

type Event interface {
	Content() string
}

type RedeemedEvent struct {
	ID   GiftID
	Type GiftType
}

func (e RedeemedEvent) Content() string {
	return fmt.Sprintf("RedeemedEvent Type: %s, ID: %s", string(e.Type), string(e.ID))
}

type IllegalAccessEvent struct {
	URL  string
	Body string
}

func (e IllegalAccessEvent) Content() string {
	return fmt.Sprintf("IllegalAccessEvent %v", e)
}

type AlreadyRedeemedEvent struct {
	ID   GiftID
	Type GiftType
}

func (e AlreadyRedeemedEvent) Content() string {
	return fmt.Sprintf("AlreadyRedeemedEvent Type: %s, ID: %s", string(e.Type), string(e.ID))
}

type NotFoundEvent struct {
	ID string
}

func (e NotFoundEvent) Content() string {
	return fmt.Sprintf("NotFoundEvent ID: %s", e.ID)
}

type EventMonitor interface {
	Track(e Event)
}

type TelegramMonitor struct {
	logger *slog.Logger
}

func NewTelegramMonitor(token, username string) *TelegramMonitor {
	logger := slog.New(slogtelegram.Option{
		Level:    slog.LevelInfo,
		Token:    token,
		Username: username,
	}.NewTelegramHandler())

	return &TelegramMonitor{logger: logger}
}

func (t TelegramMonitor) Track(e Event) {
	t.logger.Info("giftbox", "e", e.Content())
}
