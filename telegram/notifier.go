package telegram

import (
	"gopkg.in/telebot.v3"
	"pedro-go/domain"
)

type Notifier struct {
	bot      *telebot.Bot
	registry *domain.ArtistRegistry
	users    []int64
}

func (n Notifier) NotifyUsers() {
	//for _, id := range n.users {
	//	events, err := n.registry.NewEventsForUser(domain.UserID(id))
	//	if err != nil {
	//		log.Println(fmt.Errorf("error when trying to notify users %v", err))
	//	}
	//
	//	log.Printf("Sending %v\n", events)

	//_, err := n.bot.Send(user("530586914"), "Test <a href='https://www.google.com/'>Google</a>")
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}
	//}
}

type user string

func (u user) Recipient() string {
	return string(u)
}

type EventsMsg string

//func NewEventsMessage() {
//
//}
