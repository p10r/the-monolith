package domain

import (
	"strings"
)

type StatSheets []StatSheet

type StatSheet struct {
	Home, Away, Url string
}

func (matches Matches) ZipWith(sheets StatSheets) Matches {
	sheetsMap := map[string]StatSheet{}
	for _, sheet := range sheets {
		statKey := strings.ToLower(sheet.Home + "-" + sheet.Away)
		sheetsMap[statKey] = sheet
	}
	// TODO: This will break when to teams play each other in multiple tournaments
	found := Matches{}
	for _, match := range matches {
		matchKey := strings.ToLower(match.HomeName + "-" + match.AwayName)
		val, ok := sheetsMap[matchKey]
		if ok {
			match.StatsUrl = val.Url
		}
		found = append(found, match)
	}
	return found
}
