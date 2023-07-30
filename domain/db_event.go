package domain

type DbEvent struct {
	Msg string `json:"msg"`
}

func (e DbEvent) EventName() string {
	return "DbEvent"
}
