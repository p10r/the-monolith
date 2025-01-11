package domain

import (
	"slices"
	"strings"
)

// New leagues inside this must be lowercase
var favouriteLeagues = []string{
	"italy: superlega",
	"italy: superlega - play offs",
	"italy: coppa italia a1",
	"italy: coppa italia a1 women",
	"italy: serie a1 women",
	"italy: serie a1 women - playoffs",
	"poland: plusliga",
	"poland: plusliga - play offs",
	"france: ligue a - play offs",
	"france: ligue a",
	"russia: super league - play offs",
	"russia: super league",
	"russia: russia cup",
	"world: nations league",
	"world: nations league - play offs",
	"world: nations league women",
	"world: nations league women - play offs",
	"world: pan-american cup",
	"world: world championship - first round",
	"world: world championship - second round",
	"world: world championship - play offs",
	"world: world championship women - first round",
	"germany: vbl supercup",
	"germany: 1. bundesliga",
	"germany: 1. bundesliga - play offs",
	"germany: dvv cup",
	"turkey: sultanlar ligi women",
	"turkey: sultanlar ligi women - play offs",
	"turkey: efeler ligi",
	"turkey: efeler ligi - play offs",
	"turkey: efeler ligi - 5th-8th places",
	"europe: champions league",
	"europe: champions league women",
	"europe: champions league women - play offs",
	"europe: champions league - play offs",
	"europe: cev cup",
	"europe: european championships women",
	"europe: european championships",
	"japan: sv.league",
	"usa: pvf women",
}

func (matches Matches) Favourites() Matches {
	filtered := Matches{}
	for _, match := range matches {
		if slices.Contains(favouriteLeagues, strings.ToLower(match.FlashscoreName)) {
			filtered = append(filtered, match)
		}
	}

	if len(filtered) == 0 {
		return Matches{}
	}

	return filtered
}
