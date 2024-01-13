package db

import (
	d "pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestInMemoryArtistRepository(t *testing.T) {
	t.Run("verify contract for in-memory repo", func(t *testing.T) {
		d.ArtistRepositoryContract{NewRepository: func() d.ArtistRepository {
			return NewInMemoryArtistRepository()
		}}.Test(t)
	})
}

func TestGormArtistRepository(t *testing.T) {
	t.Run("verify contract for sqlite db", func(t *testing.T) {
		d.ArtistRepositoryContract{NewRepository: func() d.ArtistRepository {
			repo, _ := NewGormArtistRepository("file::memory:")
			return repo
		}}.Test(t)
	})

	//same is being mapped for domain.EventIDs
	t.Run("map domain IDs to string list", func(t *testing.T) {
		for _, tc := range []struct {
			Input d.UserIDs
			Want  commaSeparatedStr
		}{
			{
				Input: d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
				Want:  commaSeparatedStr("1,2,3,4"),
			},
			{
				Input: d.UserIDs{},
				Want:  commaSeparatedStr(""),
			},
		} {
			got := newUserIdString(tc.Input)
			expect.Equal(t, got, tc.Want)
		}
	})

	//same is being mapped for domain.EventIDs
	t.Run("map id string to domain ids", func(t *testing.T) {
		for _, tc := range []struct {
			Input commaSeparatedStr
			Want  d.UserIDs
		}{
			{
				Input: commaSeparatedStr("1,2,3,4"),
				Want:  d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
			},
			{
				Input: commaSeparatedStr(""),
				Want:  d.UserIDs{},
			},
			{
				Input: commaSeparatedStr("1 , 2 , 3 ,     4     "),
				Want:  d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
			},
		} {
			expect.DeepEqual(t, tc.Input.toUserIds(), tc.Want)
		}
	})
}
