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
	})

}
