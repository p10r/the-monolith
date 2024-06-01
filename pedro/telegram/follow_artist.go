package telegram

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pkg/l"
	"gopkg.in/telebot.v3"
	"log/slog"
)

//nolint:lll
var invalidArtistMsg = "Could not parse input, make sure to send it as follows: https://ra.co/dj/dj123"

func artistFollowedMsg(artist string) string {
	return fmt.Sprintf("You're now following %s", artist)
}

func followArtist(
	r *domain.ArtistRegistry,
	s Sender,
	log *slog.Logger,
) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		tags := c.Args()
		if len(tags) == 0 {
			return s.Send(c, invalidArtistMsg)
		}

		slug, err := domain.NewSlug(tags[0])
		if err != nil {
			log.Error(l.Error("invalid artist", err))
			return s.Send(c, invalidArtistMsg)
		}
		userId := domain.UserID(c.Sender().ID)

		err = r.Follow(ctx, slug, userId)
		if err != nil {
			log.Error(l.Error("cannot follow artist", err))
			return s.Send(c, genericErrMsg("/follow", err))
		}

		return s.Send(c, artistFollowedMsg(string(slug)))
	}
}
