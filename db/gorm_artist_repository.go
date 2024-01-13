package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"pedro-go/domain"
	"pedro-go/ra"
	"strconv"
	"strings"
)

type GormArtistRepository struct {
	db *gorm.DB
}

func NewGormArtistRepository(dsn string) (*GormArtistRepository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return &GormArtistRepository{}, fmt.Errorf("can not connect to database %v", err)
	}

	// AutoMigrate will create the table if it doesn't exist
	err = db.AutoMigrate(&artistEntity{})
	if err != nil {
		return &GormArtistRepository{}, fmt.Errorf("can not run db migration %v", err)
	}

	return &GormArtistRepository{db: db}, nil
}

func (r GormArtistRepository) Save(artist domain.Artist) (domain.Artist, error) {
	entity := &artistEntity{
		ID:            artist.ID,
		RAId:          artist.RAID,
		RASlug:        string(artist.RASlug),
		Name:          artist.Name,
		FollowedBy:    newUserIdString(artist.FollowedBy),
		TrackedEvents: newEventIDsString(artist.TrackedEvents),
	}

	r.db.Save(entity)

	artist.ID = entity.ID
	return artist, nil
}

func (r GormArtistRepository) All() (domain.Artists, error) {
	var entities []artistEntity
	r.db.Find(&entities)

	var artists []domain.Artist
	for _, e := range entities {
		a := domain.Artist{
			ID:            e.ID,
			RAID:          e.RAId,
			RASlug:        ra.Slug(e.RASlug),
			Name:          e.Name,
			FollowedBy:    e.FollowedBy.toUserIds(),
			TrackedEvents: e.TrackedEvents.toEventIds(),
		}
		artists = append(artists, a)
	}

	return artists, nil
}

// TODO move
type artistEntity struct {
	gorm.Model
	ID            int64
	RAId          string
	RASlug        string
	Name          string
	FollowedBy    commaSeparatedStr
	TrackedEvents commaSeparatedStr
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
