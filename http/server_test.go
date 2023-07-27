package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

type TestEventRecorder struct {
	Events Events
}

func (r *TestEventRecorder) Record(event Event) {
	r.Events = append(r.Events, event)
}

func TestApiRoutes(t *testing.T) {
	var (
		eventRecorder = TestEventRecorder{Events{}}
		server        = NewServer(0, &eventRecorder)
		api           = server.routes.ServeHTTP
	)

	t.Run("get all artists", func(t *testing.T) {
		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		var got domain.Artists
		err := json.Unmarshal(res.Body.Bytes(), &got)
		want := domain.Artists{domain.Artist{Name: "Boys Noize"}}

		expect.NoErr(t, err)
		expect.SliceEqual(t, got, want)
	})

	t.Run("tracks incoming calls", func(t *testing.T) {
		eventRecorder.Events = nil

		req, res := testSetupReqCtx(t, http.MethodGet, "/artists")
		api(res, req)

		got := eventRecorder.Events
		want := Events{Event{"/artists"}}

		expect.SliceEqual(t, got, want)
	})
}

func testSetupReqCtx(t *testing.T, method, url string) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("do something") //TODO
	}
	res := httptest.NewRecorder()

	return req, res
}
