package flashscore

import (
	json2 "encoding/json"
	"github.com/p10r/pedro/serve/domain"
	"io"
	"strconv"
	"strings"
)

type Response struct {
	Leagues Leagues `json:"DATA"`
}

type Leagues []League

type League struct {
	Name   string `json:"NAME"`
	Events Events `json:"EVENTS"`
}

type Events []Event

type Event struct {
	HomeName         string `json:"HOME_NAME"`
	AwayName         string `json:"AWAY_NAME"`
	StartTime        int64  `json:"START_TIME"`
	HomeScoreCurrent string `json:"HOME_SCORE_CURRENT"`
	AwayScoreCurrent string `json:"AWAY_SCORE_CURRENT"`
	Stage            string `json:"STAGE"`
	HomeSetPoints1   string `json:"HOME_SCORE_PART_1"`
	HomeSetPoints2   string `json:"HOME_SCORE_PART_2"`
	HomeSetPoints3   string `json:"HOME_SCORE_PART_3"`
	HomeSetPoints4   string `json:"HOME_SCORE_PART_4"`
	HomeSetPoints5   string `json:"HOME_SCORE_PART_5"`
	AwaySetPoints1   string `json:"AWAY_SCORE_PART_1"`
	AwaySetPoints2   string `json:"AWAY_SCORE_PART_2"`
	AwaySetPoints3   string `json:"AWAY_SCORE_PART_3"`
	AwaySetPoints4   string `json:"AWAY_SCORE_PART_4"`
	AwaySetPoints5   string `json:"AWAY_SCORE_PART_5"`
}

func NewResponse(input io.ReadCloser) (Response, error) {
	var res Response

	err := json2.NewDecoder(input).Decode(&res)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}
func (r Response) ToUntrackedMatches() domain.UntrackedMatches {
	matches := domain.UntrackedMatches{}

	for _, league := range r.Leagues {
		leagueInfo := strings.Split(league.Name, ": ")
		country := leagueInfo[0]
		leagueName := leagueInfo[1]

		for _, event := range league.Events {
			home, _ := strconv.Atoi(event.HomeScoreCurrent)
			away, _ := strconv.Atoi(event.AwayScoreCurrent)

			match := domain.UntrackedMatch{
				HomeName:         event.HomeName,
				AwayName:         event.AwayName,
				StartTime:        event.StartTime,
				FlashscoreName:   league.Name,
				Country:          country,
				League:           leagueName,
				Stage:            event.Stage,
				HomeScoreCurrent: home,
				AwayScoreCurrent: away,
			}

			matches = append(matches, match)
		}
	}

	return matches
}
