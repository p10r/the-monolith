package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"pedro-go/domain"
	"strconv"
)

type PedroServer struct {
	routes        http.Handler
	port          int
	EventRecorder domain.EventRecorder
	Registry      domain.ArtistRegistry
}

type Event struct {
	Uri string
}

func (e Event) EventName() string {
	return "HttpEvent"
}

func NewServer(port int, recorder domain.EventRecorder, registry domain.ArtistRegistry) PedroServer {
	s := PedroServer{port: port, EventRecorder: recorder, Registry: registry}

	r := chi.NewRouter()
	r.Use(s.logIncoming)
	r.HandleFunc("/artists", s.getAllArtists)
	s.routes = r

	return s
}

func (s PedroServer) Start() {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.port), s.routes))
}

func (s PedroServer) getAllArtists(w http.ResponseWriter, _ *http.Request) {
	all, _ := s.Registry.FindAll(context.TODO()) //TODO check context
	artists, _ := json.Marshal(all)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(artists)
}

func (s PedroServer) logIncoming(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.EventRecorder.Record(Event{Uri: r.URL.RequestURI()})
		next.ServeHTTP(w, r)
	})
}
