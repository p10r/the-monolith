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
	t.Run("get all artists", func(t *testing.T) {
		api, _ := newTestApiFixture(nil)
		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		var got domain.Artists
		err := json.Unmarshal(res.Body.Bytes(), &got)
		want := domain.Artists{domain.Artist{Name: "Boys Noize"}}

		expect.NoErr(t, err)
		expect.SliceEqual(t, got, want)
	})

	t.Run("tracks incoming calls", func(t *testing.T) {
		initialEvents := domain.Events{}
		api, eventRecorder := newTestApiFixture(initialEvents)

		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		got := eventRecorder.Events
		want := domain.Events{HttpEvent{"/artists"}}

		expect.SliceEqual(t, got, want)
	})
}

func newTestApiFixture(events domain.Events) (func(http.ResponseWriter, *http.Request), *domain.JsonEventRecorder) {
	if events == nil {
		events = domain.Events{}
	}

	var (
		registry      = inmemory.ArtistRegistry{}
		eventRecorder = domain.NewEventRecorder(events)
		server        = NewServer(0, &eventRecorder, registry)
		api           = server.routes.ServeHTTP
	)
	return api, &eventRecorder
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
