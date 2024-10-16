package specs

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/giftbox"
	"github.com/p10r/pedro/giftbox/specs/env"
	"net/http"
	"strings"
	"testing"
)

func TestAddGift(t *testing.T) {
	t.Run("adds a sweet", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))

		res := server.AddSweet()
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, `{"id":"1"}`, strings.TrimSpace(res.Body.String()))

		_, ok := server.CheckStoreFor("1")
		assert.True(t, ok)
	})

	t.Run("redeems sweet", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))

		server.AddSweet()
		res := server.RedeemGift("1")
		assert.Equal(t, http.StatusOK, res.Code)

		expected := giftbox.Gift{ID: "1", Type: "SWEET", Redeemed: true}
		value, ok := server.CheckStoreFor("1")
		assert.True(t, ok)
		assert.Equal(t, expected, value)
	})

	t.Run("blocks redeeming a gift twice", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))

		server.AddSweet()
		res := server.RedeemGift("1")
		assert.Equal(t, http.StatusOK, res.Code)

		res = server.RedeemGift("1")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("returns 400 if no id is given", func(t *testing.T) {
		server := env.NewInMemoryEnv(t, int32(0))
		emptyID := ""

		server.AddSweet()
		res := server.RedeemGift(emptyID)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("shows status of all gifts", func(t *testing.T) {

	})
}
