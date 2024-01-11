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

func (r GormArtistRepository) All() (domain.Artists, error) {
	var entities []artistEntity
	r.db.Find(&entities)

	var artists []domain.Artist
	for _, e := range entities {
		a := domain.Artist{
			Id:         e.ID,
			RAId:       e.RAId,
			RASlug:     ra.Slug(e.RASlug),
			Name:       e.Name,
			FollowedBy: e.FollowedBy.ToUserIds(),
		}
		artists = append(artists, a)
	}

	return artists, nil
}

func (r GormArtistRepository) Save(artist domain.Artist) (domain.Artist, error) {
	entity := &artistEntity{
		ID:         artist.Id,
		RAId:       artist.RAId,
		RASlug:     string(artist.RASlug),
		Name:       artist.Name,
		FollowedBy: NewUserIdString(artist.FollowedBy),
	}

	r.db.Save(entity)

	artist.Id = entity.ID
	return artist, nil
}

type artistEntity struct {
	gorm.Model
	ID         int64
	RAId       string
	RASlug     string
	Name       string
	FollowedBy UserIdsString
}

type UserIdsString string

func NewUserIdString(ids domain.UserIds) UserIdsString {
	var strIds []string
	for _, id := range ids {
		strIds = append(strIds, strconv.FormatInt(int64(id), 10))
	}
	return UserIdsString(strings.Join(strIds, ","))
}

func (r UserIdsString) ToUserIds() domain.UserIds {
	input := string(r)
	if len(strings.TrimSpace(input)) == 0 {
		return domain.UserIds{}
	}

	var ids domain.UserIds
	for _, s := range strings.Split(input, ",") {
		i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			log.Printf("SKIPPING: Could not convert '%v' to UserIds - this should never happen\n", r)
			continue
		}

		ids = append(ids, domain.UserId(i))
	}

	return ids
}
