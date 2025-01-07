package domain

import (
	"slices"
	"strings"
)

type MatchStage = string

const (
	SCHEDULED MatchStage = "SCHEDULED"
	FINISHED  MatchStage = "FINISHED"
)

func (matches Matches) FilterScheduled(favs []string) Matches {
	scheduled := matches.filterByStage(SCHEDULED)
	if len(scheduled) == 0 {
		return Matches{}
	}

	filtered := scheduled.filterFavourites(favs) //TODO
	if len(filtered) == 0 {
		return Matches{}
	}

	return filtered
}

func (matches Matches) FilterFinished(favourites []string) Matches {
	finished := matches.filterByStage(FINISHED)
	if len(finished) == 0 {
		return Matches{}
	}

	filtered := finished.filterFavourites(favourites)
	if len(filtered) == 0 {
		return Matches{}
	}

	return filtered
}

func (matches Matches) filterByStage(stage string) Matches {
	filtered := Matches{}

	for _, match := range matches {
		if lowerCase(match.Stage) == lowerCase(stage) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

// TODO move favourites to struct that has Country and League separate
func (matches Matches) filterFavourites(favourites []string) Matches {
	var favs []string
	for _, favourite := range favourites {
		favs = append(favs, lowerCase(favourite))
	}

	filtered := Matches{}
	for _, match := range matches {
		if slices.Contains(favs, lowerCase(match.FlashscoreName)) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

func lowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
