package specs

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/giftbox"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestAddGift(t *testing.T) {
	t.Run("adds a sweet", func(t *testing.T) {
		var store sync.Map

		req := httptest.NewRequest("POST", "/gifts/sweets", nil)
		w := httptest.NewRecorder()

		id := func() (string, error) {
			return "1", nil
		}
		s := giftbox.NewServer(&store, id)
		s.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, `{"id":"1"}`, strings.TrimSpace(w.Body.String()))

		_, ok := store.Load("1")
		assert.True(t, ok)
	})

	t.Run("redeems sweet", func(t *testing.T) {
		var store sync.Map

		req := httptest.NewRequest("POST", "/gifts/sweets", nil)
		w := httptest.NewRecorder()

		id := func() (string, error) {
			return "1", nil
		}
		s := giftbox.NewServer(&store, id)
		s.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		req2 := httptest.NewRequest("POST", "/gifts/redeem?id=1", nil)
		w2 := httptest.NewRecorder()
		s.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		expected := giftbox.Gift{ID: "1", Type: "SWEET", Redeemed: true}
		value, ok := store.Load("1")
		assert.True(t, ok)
		assert.Equal(t, expected, value.(giftbox.Gift))
	})

	t.Run("blocks redeeming a gift twice", func(t *testing.T) {
		var store sync.Map

		req := httptest.NewRequest("POST", "/gifts/sweets", nil)
		w := httptest.NewRecorder()

		id := func() (string, error) {
			return "1", nil
		}
		s := giftbox.NewServer(&store, id)
		s.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		req2 := httptest.NewRequest("POST", "/gifts/redeem?id=1", nil)
		w2 := httptest.NewRecorder()
		s.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		expected := giftbox.Gift{ID: "1", Type: "SWEET", Redeemed: true}
		value, ok := store.Load("1")
		assert.True(t, ok)
		assert.Equal(t, expected, value.(giftbox.Gift))

		req3 := httptest.NewRequest("POST", "/gifts/redeem?id=1", nil)
		w3 := httptest.NewRecorder()
		s.ServeHTTP(w3, req3)

		assert.Equal(t, http.StatusBadRequest, w3.Code)

		expected = giftbox.Gift{ID: "1", Type: "SWEET", Redeemed: true}
		value, ok = store.Load("1")
		assert.True(t, ok)
		assert.Equal(t, expected, value.(giftbox.Gift))
	})

	t.Run("returns 400 if no id is given", func(t *testing.T) {
		var store sync.Map
		id := func() (string, error) {
			return "1", nil
		}
		s := giftbox.NewServer(&store, id)
		req := httptest.NewRequest("POST", "/gifts/redeem", nil)
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("shows status of all gifts", func(t *testing.T) {

	})
}
