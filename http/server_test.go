package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestApiRoutes(t *testing.T) {
	t.Run("get all artists", func(t *testing.T) {
		artists := domain.Artists{domain.Artist{Name: "Boys Noize"}}
		api, _ := newTestApiFixture(nil, artists)

		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		var got domain.Artists
		err := json.Unmarshal(res.Body.Bytes(), &got)
		want := artists

		expect.NoErr(t, err)
		expect.SliceEqual(t, got, want)
	})

	t.Run("tracks incoming calls", func(t *testing.T) {
		initialEvents := domain.Events{}
		artists := domain.Artists{domain.Artist{Name: "Boys Noize"}}
		api, eventRecorder := newTestApiFixture(initialEvents, artists)

		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		got := eventRecorder.Events
		want := domain.Events{HttpEvent{"/artists"}}

		expect.SliceEqual(t, got, want)
	})
}

func newTestApiFixture(events domain.Events, storedArtists domain.Artists) (func(http.ResponseWriter, *http.Request), *domain.JsonEventRecorder) {
	if events == nil {
		events = domain.Events{}
	}
	eventRecorder := domain.NewEventRecorder(events)

	var (
		registry = ArtistRegistry{id: 0, Artists: storedArtists}
		server   = NewServer(0, &eventRecorder, registry)
		api      = server.routes.ServeHTTP
	)
	return api, &eventRecorder
}

type ArtistRegistry struct {
	id int
	domain.Artists
}

func (r ArtistRegistry) FindAll(ctx context.Context) (domain.Artists, error) {
	return r.Artists, nil
}

func (r ArtistRegistry) Add(_ context.Context, artist domain.NewArtist) (domain.Artist, error) {
	id := r.id + 1
	storedArtist := domain.Artist{Id: domain.Id(id), Name: artist.Name}

	r.Artists = append(r.Artists, storedArtist)

	return storedArtist, nil
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
