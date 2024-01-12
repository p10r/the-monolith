package ra

import (
	"errors"
	"regexp"
)

type Slug string

func NewSlug(url string) (Slug, error) {
	re := regexp.MustCompile(`(?:https?://ra\.co/dj/|ra\.co/dj/)([a-zA-Z]+)`)

	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return Slug(match[1]), nil
	}
	return "", errors.New("could not parse artist")
}
