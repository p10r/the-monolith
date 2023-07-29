package domain

type Event interface {
	EventName() string
}

type Events []Event

type EventRecorder interface {
	Record(event Event)
}

type TestEventRecorder struct {
	Events Events
}

func (r *TestEventRecorder) Record(event Event) {
	r.Events = append(r.Events, event)
}
