package db

import (
	"log"
	"pedro-go/pedro/domain"
	"strconv"
	"strings"
)

type artistDBEntity struct {
	ID            int64
	RAID          string
	RASlug        string
	Name          string
	FollowedBy    commaSeparatedStr
	TrackedEvents commaSeparatedStr
}

func newArtistDBEntity(artist domain.Artist) artistDBEntity {
	return artistDBEntity{
		ID:            artist.ID,
		RAID:          artist.RAID,
		RASlug:        string(artist.RASlug),
		Name:          artist.Name,
		FollowedBy:    newUserIdString(artist.FollowedBy),
		TrackedEvents: newEventIDsString(artist.TrackedEvents),
	}
}

type commaSeparatedStr string

func newUserIdString(ids domain.UserIDs) commaSeparatedStr {
	var strIds []string
	for _, id := range ids {
		strIds = append(strIds, strconv.FormatInt(int64(id), 10))
	}
	return commaSeparatedStr(strings.Join(strIds, ","))
}

func newEventIDsString(ids domain.EventIDs) commaSeparatedStr {
	var strIds []string
	for _, id := range ids {
		strIds = append(strIds, strconv.FormatInt(int64(id), 10))
	}
	return commaSeparatedStr(strings.Join(strIds, ","))
}

func (r commaSeparatedStr) toUserIds() domain.UserIDs {
	ids := r.toInt64Slice()
	if len(ids) == 0 {
		return domain.UserIDs{}
	}

	var userIds domain.UserIDs
	for _, i := range ids {
		userIds = append(userIds, domain.UserID(i))
	}
	return userIds
}

func (r commaSeparatedStr) toEventIds() domain.EventIDs {
	ids := r.toInt64Slice()
	if len(ids) == 0 {
		return domain.EventIDs{}
	}

	var eventIds domain.EventIDs
	for _, i := range ids {
		eventIds = append(eventIds, domain.EventID(i))
	}
	return eventIds
}

func (r commaSeparatedStr) toInt64Slice() []int64 {
	input := string(r)
	if len(strings.TrimSpace(input)) == 0 {
		return []int64{}
	}

	var ids []int64
	for _, s := range strings.Split(input, ",") {
		i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			log.Printf("SKIPPING: Could not convert '%v' to int - this should never happen\n", r)
			continue
		}

		ids = append(ids, i)
	}

	return ids
}
