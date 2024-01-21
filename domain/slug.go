package domain

import (
	"errors"
	"regexp"
)

type RASlug string

func NewSlug(url string) (RASlug, error) {
	re := regexp.MustCompile(`(?:https?://ra\.co/dj/|ra\.co/dj/)([a-zA-Z]+)`)

	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return RASlug(match[1]), nil
	}
	return "", errors.New("could not parse artist")
}
