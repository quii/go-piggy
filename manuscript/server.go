package manuscript

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/quii/go-piggy"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//todo: i have totally ruined this file by hacking ;/

type Repo interface {
	GetManuscript(id string) Manuscript
	GetVersionedManuscript(entityID string, version int) (Manuscript, error)
	Versions(entityID string) int
	Events(entityID string) []go_piggy.Event
}

type Server struct {
	Repo              Repo
	Emitter           go_piggy.Emitter
	EntityIdGenerator func() string
	handler           http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func NewServer(repo Repo, emitter go_piggy.Emitter, options ...func(*Server)) *Server {
	s := new(Server)
	s.Repo = repo
	s.Emitter = emitter
	s.EntityIdGenerator = go_piggy.RandomID

	for _, op := range options {
		op(s)
	}

	r := mux.NewRouter()

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	r.HandleFunc("/manuscripts/{entityID}", s.getManuscriptJSON)
	r.HandleFunc("/manuscripts", s.createManuscript).Methods("POST")
	r.HandleFunc("/manuscripts/{entityID}/events", s.addEventsToManuscript).Methods("POST")
	r.HandleFunc("/manuscripts/{entityID}/events", s.showEvents).Methods("GET")

	s.handler = cors(r)

	return s
}

func WithEntityIdGenerator(f func() string) func(*Server) {
	return func(server *Server) {
		server.EntityIdGenerator = f
	}
}

func (s *Server) addEventsToManuscript(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)

	var facts []go_piggy.Fact

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(body, &facts)
	s.Emitter.Send(NewManuscriptVersionEvent(Manuscript{EntityID: entityID}, facts...))
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) createManuscript(w http.ResponseWriter, r *http.Request) {
	newEntityID := s.EntityIdGenerator()

	s.Emitter.Send(NewManuscriptEvent(Manuscript{
		EntityID: newEntityID,
	}))

	location := "/manuscripts/" + newEntityID
	log.Println("manuscript created at", location)
	w.Header().Add("location", location)
	w.Header().Add("Access-Control-Expose-Headers", "Location")
	w.WriteHeader(http.StatusCreated)

}

//todo: testme
func (s *Server) showEvents(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)
	events := s.Repo.Events(entityID)

	eventsAsJSON, _ := json.Marshal(events)

	w.Header().Set("content-type", "application/json")
	w.Write(eventsAsJSON)
}

func (s *Server) getManuscriptJSON(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)
	version := r.URL.Query().Get("version")

	var manuscript Manuscript

	if version != "" {
		v, _ := strconv.Atoi(version)
		m, _ := s.Repo.GetVersionedManuscript(entityID, v)
		manuscript = m
	} else {
		manuscript = s.Repo.GetManuscript(entityID)
	}

	if manuscript.EntityID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	manuscriptJSON, _ := json.Marshal(manuscript)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(manuscriptJSON)
}

func entityIDFromRequest(r *http.Request) string {
	return mux.Vars(r)["entityID"]
}
