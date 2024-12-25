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

func (matches UntrackedMatches) FilterScheduled(favs []string) UntrackedMatches {
	scheduled := filter(SCHEDULED, matches)
	if len(scheduled) == 0 {
		return UntrackedMatches{}
	}

	filtered := filterFavourites(scheduled, favs) //TODO
	if len(filtered) == 0 {
		return UntrackedMatches{}
	}

	return filtered
}

func (matches UntrackedMatches) FilterFinished(favourites []string) FinishedMatches {
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
			match.HomeName,
			match.HomeScoreCurrent,
			match.AwayName,
			match.AwayScoreCurrent,
			// Empty by default - set through statistics package
			"",
		})
	}

	return finished
}

func filter(stage string, flashscoreMatches UntrackedMatches) UntrackedMatches {
	filtered := UntrackedMatches{}

	for _, match := range flashscoreMatches {
		if lowerCase(match.Stage) == lowerCase(stage) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

// TODO move favourites to struct that has Country and League separate
func filterFavourites(matches UntrackedMatches, favourites []string) UntrackedMatches {
	var favs []string
	for _, favourite := range favourites {
		favs = append(favs, lowerCase(favourite))
	}

	filtered := UntrackedMatches{}
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
