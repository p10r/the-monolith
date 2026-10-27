package serve

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/p10r/monolith/pkg/l"
)

type discordMessage struct {
	Content string         `json:"content"`
	Embeds  []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Fields []discordField `json:"fields"`
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func (s *Serve) sendScheduledMatches(matches Matches, now time.Time) error {
	msg := newScheduledMatchesMessage(matches, now)
	return s.sendDiscordMessage(msg)
}

func (s *Serve) sendFinishedMatches(matches MatchesByLeague, now time.Time) error {
	msg := newFinishedMatchesMessage(matches, now)
	return s.sendDiscordMessage(msg)
}

func (s *Serve) sendDiscordMessage(msg discordMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal discord message: %w", err)
	}

	const maxRetries = 5
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err := s.httpClient.Post(s.discordURL, "application/json", bytes.NewReader(payload))
		if err != nil {
			lastErr = err
			s.log.Error(l.Error("send discord message", err))
			time.Sleep(time.Second)
			continue
		}
		res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf("discord request failed: %s", res.Status)
		}
		return nil
	}
	return fmt.Errorf("discord request failed after %d attempts: %w", maxRetries, lastErr)
}

func newScheduledMatchesMessage(matches Matches, now time.Time) discordMessage {
	var fields []discordField
	for league, matchesForLeague := range groupByLeagueName(matches) {
		fullName := matchesForLeague[0].Country + ": " + matchesForLeague[0].League
		fields = append(fields, discordField{
			Name:  flag(NewLeagueKey(league)) + strings.ToTitle(fullName),
			Value: scheduledText(matchesForLeague),
		})
	}

	msg := discordMessage{
		Content: fmt.Sprintf("Games for %s", now.Format("Monday, 2 January 2006")),
		Embeds:  []discordEmbed{{Fields: fields}},
	}

	// Keep a stable order so the message is predictable.
	slices.SortFunc(msg.Embeds[0].Fields, func(a, b discordField) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return msg
}

func newFinishedMatchesMessage(matches MatchesByLeague, now time.Time) discordMessage {
	var fields []discordField
	for league, matchesForLeague := range matches {
		fields = append(fields, discordField{
			Name:  flag(league) + string(league),
			Value: finishedText(matchesForLeague),
		})
	}

	return discordMessage{
		Content: fmt.Sprintf("Results for %s", now.Format("Monday, 2 January 2006")),
		Embeds:  []discordEmbed{{Fields: fields}},
	}
}

func groupByLeagueName(matches Matches) map[string]Matches {
	grouped := make(map[string]Matches)
	for _, m := range matches {
		fullName := m.Country + ": " + m.League
		grouped[fullName] = append(grouped[fullName], m)
	}
	return grouped
}

func scheduledText(matches Matches) string {
	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("**%s - %s**\t %s", m.HomeName, m.AwayName, hour(m.StartTime)))
	}
	return strings.Join(lines, "\n")
}

func finishedText(matches Matches) string {
	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("**%s - %s**\t\t\tScore:\t||%d\t:\t%d||",
			m.HomeName, m.AwayName, m.HomeScoreCurrent, m.AwayScoreCurrent))
	}
	return strings.Join(lines, "\n")
}

func flag(key LeagueKey) string {
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

// hour formats a unix timestamp as a Discord relative timestamp.
// See https://hammertime.cyou/ for more info.
func hour(unixTs int64) string {
	return fmt.Sprintf("<t:%s:t>", strconv.FormatInt(unixTs, 10))
}
