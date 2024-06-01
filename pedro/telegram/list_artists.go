package telegram

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pkg/l"
	"gopkg.in/telebot.v3"
	"log/slog"
	"strings"
)

func genericErrMsg(endpoint string, err error) string {
	return fmt.Sprintf("There was an error when calling %s! err: %v", endpoint, err)
}

func listArtistsMsg(artists domain.Artists) (string, error) {
	var res []string
	for _, artist := range artists {
		res = append(res, "- "+artist.Name)
	}

	if len(res) == 0 {
		//goland:noinspection ALL
		return "", fmt.Errorf("You're not following anyone yet.")
	}

	return fmt.Sprintf("You're following:\n%v", strings.Join(res, "\n")), nil
}

func listArtists(
	r *domain.ArtistRegistry,
	s Sender,
	log *slog.Logger,
) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		id := c.Sender().ID
		artists, err := r.ArtistsFor(ctx, domain.UserID(id))
		if err != nil {
			log.Error(l.Error(fmt.Sprintf("%v has no artists", id), err))
			return c.Send(genericErrMsg("/artists", err))
		}

		msg, err := listArtistsMsg(artists)
		if err != nil {
			log.Error(l.Error(fmt.Sprintf("%v has no artists", id), err))
			return c.Send(err.Error())
		}

		return s.Send(c, msg)
	}
}
