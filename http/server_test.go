package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pedro-go/db/inmemory"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestApiRoutes(t *testing.T) {
	var (
		artists       = domain.Artists{domain.Artist{Name: "Boys Noize"}}
		registry      = inmemory.ArtistRegistry{}
		eventRecorder = domain.TestEventRecorder{domain.Events{}}
		server        = NewServer(0, &eventRecorder, registry)
		api           = server.routes.ServeHTTP
	)

	t.Run("get all artists", func(t *testing.T) {
		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		var got domain.Artists
		err := json.Unmarshal(res.Body.Bytes(), &got)
		want := artists

		expect.NoErr(t, err)
		expect.SliceEqual(t, got, want)
	})

	t.Run("tracks incoming calls", func(t *testing.T) {
		eventRecorder.Events = nil

		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		got := eventRecorder.Events
		want := domain.Events{HttpEvent{"/artists"}}

		expect.SliceEqual(t, got, want)
	})
}

func testSetupReqCtx(t *testing.T, method, url string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("do something") //TODO
	}
	res := httptest.NewRecorder()

	return req, res
}
