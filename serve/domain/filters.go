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
	scheduled := filter(SCHEDULED, matches)
	if len(scheduled) == 0 {
		return Matches{}
	}

	filtered := filterFavourites(scheduled, favs) //TODO
	if len(filtered) == 0 {
		return Matches{}
	}

	return filtered
}

func (matches Matches) FilterFinished(favourites []string) FinishedMatches {
	scheduled := filter(FINISHED, matches)
	if len(scheduled) == 0 {
		return FinishedMatches{}
	}

	filtered := filterFavourites(scheduled, favourites) //TODO
	if len(filtered) == 0 {
		return FinishedMatches{}
	}

	var finished FinishedMatches
	for _, match := range filtered {
		finished = append(finished, FinishedMatch{
			match,
			// Empty by default - set through statistics package
			"",
		})
	}

	return finished
}

func filter(stage string, flashscoreMatches Matches) Matches {
	filtered := Matches{}

	for _, match := range flashscoreMatches {
		if lowerCase(match.Stage) == lowerCase(stage) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

// TODO move favourites to struct that has Country and League separate
func filterFavourites(matches Matches, favourites []string) Matches {
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
