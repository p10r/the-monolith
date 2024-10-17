package specs

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/giftbox"
	"github.com/p10r/pedro/giftbox/specs/env"
	"github.com/p10r/pedro/pkg/sqlite"
	"net/http"
	"strings"
	"testing"
)

func TestAddGift(t *testing.T) {
	t.Parallel()

	t.Run("adds a sweet", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		defer sqlite.MustCloseDB(t, server.DB)

		res := server.AddSweet()
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, `{"id":"1"}`, strings.TrimSpace(res.Body.String()))

		gift, ok := server.FindInDB(t, "1")
		assert.True(t, ok)
		assert.Equal(t, giftbox.TypeSweet, gift.Type)
	})

	t.Run("adds a wish", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		defer sqlite.MustCloseDB(t, server.DB)

		res := server.AddWish()
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, `{"id":"1"}`, strings.TrimSpace(res.Body.String()))

		gift, ok := server.FindInDB(t, "1")
		assert.True(t, ok)
		assert.Equal(t, giftbox.TypeWish, gift.Type)
	})

	t.Run("redeems gifts", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		defer sqlite.MustCloseDB(t, server.DB)

		server.AddSweet()
		server.RedeemGift("1")
		value, _ := server.FindInDB(t, "1")
		assert.Equal(t, giftbox.Gift{ID: "1", Type: giftbox.TypeSweet, Redeemed: true}, value)

		server.AddWish()
		server.RedeemGift("2")
		value, _ = server.FindInDB(t, "2")
		assert.Equal(t, giftbox.Gift{ID: "2", Type: giftbox.TypeWish, Redeemed: true}, value)
	})

	t.Run("blocks redeeming a gift twice", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		defer sqlite.MustCloseDB(t, server.DB)

		server.AddSweet()
		res := server.RedeemGift("1")
		assert.Equal(t, http.StatusOK, res.Code)

		res = server.RedeemGift("1")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("returns 400 if no id is given", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		defer sqlite.MustCloseDB(t, server.DB)
		emptyID := ""

		server.AddSweet()
		res := server.RedeemGift(emptyID)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("shows status of all gifts", func(t *testing.T) {

	})
}
