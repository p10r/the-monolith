package domain

import (
	"golang.org/x/exp/slog"
	"os"
)

type Event interface {
	EventName() string
}

type Events []Event

type ErrEvent struct {
	Err string `json:"err"`
}

func (e ErrEvent) EventName() string {
	return "ErrEvent"
}

type EventRecorder interface {
	Record(event Event)
}

type JsonEventRecorder struct {
	logger *slog.Logger
	Events Events
}

func NewEventRecorder(events Events) JsonEventRecorder {
	return JsonEventRecorder{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Events: events,
	}
}

func (r *JsonEventRecorder) Record(event Event) {
	r.logger.Info(event.EventName(), "data", event)
	r.Events = append(r.Events, event)
}
