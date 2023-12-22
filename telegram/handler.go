package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	Pedro(os.Getenv("TELEGRAM_TOKEN"))
}

func Pedro(botToken string) {
	pref := tele.Settings{
		Token:   botToken,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: true,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/hello", func(c tele.Context) error {
		tags := c.Args() // list of arguments split by a space
		for _, tag := range tags {
			println(tag)
		}

		return c.Send("Hello!")
	})

	bot.Start()
}
