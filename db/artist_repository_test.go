package db_test

import (
	"pedro-go/db"
	d "pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestInMemoryArtistRepository(t *testing.T) {
	t.Run("verify contract for in-memory repo", func(t *testing.T) {
		d.ArtistRepositoryContract{NewRepository: func() d.ArtistRepository {
			return db.NewInMemoryArtistRepository()
		}}.Test(t)
	})
}

func TestGormArtistRepository(t *testing.T) {
	t.Run("verify contract for sqlite db", func(t *testing.T) {
		d.ArtistRepositoryContract{NewRepository: func() d.ArtistRepository {
			repo, _ := db.NewGormArtistRepository("file::memory:")
			return repo
		}}.Test(t)
	})

	t.Run("map user IDs to string list", func(t *testing.T) {
		for _, tc := range []struct {
			Input d.UserIds
			Want  db.UserIdsString
		}{
			{
				Input: d.UserIds{d.UserId(1), d.UserId(2), d.UserId(3), d.UserId(4)},
				Want:  db.UserIdsString("1,2,3,4"),
			},
			{
				Input: d.UserIds{},
				Want:  db.UserIdsString(""),
			},
		} {
			got := db.NewUserIdString(tc.Input)
			expect.Equal(t, got, tc.Want)
		}
	})

	t.Run("map id string to user ids", func(t *testing.T) {
		for _, tc := range []struct {
			Input db.UserIdsString
			Want  d.UserIds
		}{
			{
				Input: db.UserIdsString("1,2,3,4"),
				Want:  d.UserIds{d.UserId(1), d.UserId(2), d.UserId(3), d.UserId(4)},
			},
			{
				Input: db.UserIdsString(""),
				Want:  d.UserIds{},
			},
			{
				Input: db.UserIdsString("1 , 2 , 3 ,     4     "),
				Want:  d.UserIds{d.UserId(1), d.UserId(2), d.UserId(3), d.UserId(4)},
			},
		} {
			expect.DeepEqual(t, tc.Input.ToUserIds(), tc.Want)
		}
	})
}
