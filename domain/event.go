package domain

type Event struct {
	Id         string
	Title      string
	Venue      string
	Date       string
	StartTime  string
	ContentUrl string
}

type Events []Event
