package giftbox

import (
	"fmt"
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

func NewTelegramMonitor(log *slog.Logger) *TelegramMonitor {
	return &TelegramMonitor{logger: log}
}

func (t TelegramMonitor) Track(e Event) {
	t.logger.Info("giftbox", "e", e.Content())
}
