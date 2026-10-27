package discord

import (
	"cmp"
	"fmt"
	"github.com/p10r/monolith/serve/domain"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Content string   `json:"content"`
	Embeds  []Embeds `json:"embeds"`
}

type Embeds struct {
	Fields []Fields `json:"fields"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func NewUpcomingMatchesMsg(matches domain.Matches, currentTime time.Time) Message {
	date := currentTime.Format("Monday, 2 January 2006")

	var fields []Fields
	for league, matchesForCountry := range sortUpcomingByLeague(matches) {
		fullName := matchesForCountry[0].Country + ": " + matchesForCountry[0].League
		fields = append(fields, Fields{
			Name:   flag(domain.NewLeagueKey(league)) + strings.ToTitle(fullName),
			Value:  upcomingText(matchesForCountry),
			Inline: false,
		})
	}

	msg := Message{
		Content: fmt.Sprintf("Games for %s", date),
		Embeds:  []Embeds{{fields}},
	}

	// Sort, to always have the same order in the message to help approval tests
	slices.SortFunc(msg.Embeds[0].Fields, func(a, b Fields) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return msg
}

func NewFinishedMatchesMsg(
	matches domain.MatchesByLeague,
	currentTime time.Time,
) Message {
	var fields []Fields
	for league, matchesForCountry := range matches {
		f := Fields{
			Name:   flag(league) + string(league),
			Value:  finishedText(matchesForCountry),
			Inline: false,
		}

		fields = append(fields, f)
	}

	return Message{
		Content: fmt.Sprintf("Results for %s", currentTime.Format("Monday, 2 January 2006")),
		Embeds:  []Embeds{{fields}},
	}
}

func sortUpcomingByLeague(matches domain.Matches) map[string]domain.Matches {
	countries := make(map[string]domain.Matches)
	for _, m := range matches {
		fullName := m.Country + ": " + m.League
		countries[fullName] = append(countries[fullName], m)
	}
	return countries
}

func upcomingText(matches domain.Matches) string {
	var texts []string
	for _, e := range matches {
		formatted := fmt.Sprintf("**%v - %v**\t %v", e.HomeName, e.AwayName, hour(e.StartTime))
		texts = append(texts, formatted)
	}
	return strings.Join(texts, "\n")
}

func finishedText(matches domain.Matches) string {
	var texts []string
	for _, m := range matches {
		texts = append(texts, fmt.Sprintf("**%s - %s**\t\t\tScore:\t||%d\t:\t%d||",
			m.HomeName, m.AwayName, m.HomeScoreCurrent, m.AwayScoreCurrent))
	}

	return strings.Join(texts, "\n")
}

func flag(key domain.LeagueKey) string {
	if key.CountryEquals("Poland") {
		return "🇵🇱"
	}
	if key.CountryEquals("Italy") {
		return "🇮🇹"
	}
	if key.CountryEquals("France") {
		return "🇫🇷"
	}
	if key.CountryEquals("Germany") {
		return "🇩🇪"
	}
	if key.CountryEquals("Russia") {
		return "🇷🇺"
	}
	if key.CountryEquals("Turkey") {
		return "🇹🇷"
	}
	if key.CountryEquals("Europe") {
		return "🇪🇺"
	}
	if key.CountryEquals("USA") {
		return "🇺🇸"
	}
	if key.CountryEquals("Japan") {
		return "🇯🇵"
	}
	return ""
}

// See https://hammertime.cyou/ for more info
func hour(unixTs int64) string {
	return fmt.Sprintf("<t:%s:t>", strconv.FormatInt(unixTs, 10))
}
