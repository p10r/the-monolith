package discord

import (
	"fmt"
	"github.com/p10r/pedro/serve/domain"
	"sort"
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
			Name:   flag(league) + fullName,
			Value:  upcomingText(matchesForCountry),
			Inline: false,
		})
	}

	return Message{fmt.Sprintf("Games for %s", date), []Embeds{{fields}}}
}

func NewFinishedMatchesMsg(matches domain.FinishedMatches, currentTime time.Time) Message {
	date := currentTime.Format("Monday, 2 January 2006")

	var fields []Fields
	for league, matchesForCountry := range sortFinishedByLeague(matches) {
		fullName := matchesForCountry[0].Country + ": " + matchesForCountry[0].League
		fields = append(fields, Fields{
			Name:   flag(league) + fullName,
			Value:  finishedText(matchesForCountry),
			Inline: false,
		})
	}

	return Message{fmt.Sprintf("Results for %s", date), []Embeds{{fields}}}
}

func sortUpcomingByLeague(matches domain.Matches) map[string]domain.Matches {
	countries := make(map[string]domain.Matches)
	for _, m := range matches {
		//TODO store leagues also in DB, use this identifier here
		fullName := m.Country + ": " + m.League
		countries[fullName] = append(countries[fullName], m)
	}
	return countries
}

func sortFinishedByLeague(matches domain.FinishedMatches) map[string]domain.FinishedMatches {
	countries := make(map[string]domain.FinishedMatches)
	for _, m := range matches {
		//TODO store leagues also in DB, use this identifier here
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

func finishedText(matches domain.FinishedMatches) string {
	// Sort, to always have the same order in the message
	sort.Slice(matches, func(i, j int) bool {
		one := matches[i].HomeName
		other := matches[j].HomeName

		return len(one) > len(other)
	})

	var texts []string
	for _, m := range matches {
		//nolint:lll
		formatted := fmt.Sprintf("**%v - %v**\t\t\tScore:\t||%v\t:\t%v||", m.HomeName, m.AwayName, m.HomeSetScore, m.AwaySetScore)
		if m.StatsUrl != "" {
			formatted = formatted + "\t\t\t[📊 Statistics](" + m.StatsUrl + ")"
		}
		texts = append(texts, formatted)
	}
	return strings.Join(texts, "\n")
}

func flag(leagueName string) string {
	if strings.Contains(leagueName, "Poland") {
		return "🇵🇱"
	}
	if strings.Contains(leagueName, "Italy") {
		return "🇮🇹"
	}
	if strings.Contains(leagueName, "France") {
		return "🇫🇷"
	}
	if strings.Contains(leagueName, "Germany") {
		return "🇩🇪"
	}
	if strings.Contains(leagueName, "Russia") {
		return "🇷🇺"
	}
	if strings.Contains(leagueName, "Turkey") {
		return "🇹🇷"
	}
	if strings.Contains(leagueName, "Europe") {
		return "🇪🇺"
	}
	if strings.Contains(leagueName, "USA") {
		return "🇺🇸"
	}
	if strings.Contains(leagueName, "Japan") {
		return "🇯🇵"
	}
	return ""
}

// See https://hammertime.cyou/ for more info
func hour(unixTs int64) string {
	return fmt.Sprintf("<t:%s:t>", strconv.FormatInt(unixTs, 10))
}
