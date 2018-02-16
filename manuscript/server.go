package manuscript

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/quii/go-piggy"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Repo interface {
	Events(entityID string) []go_piggy.Event
	GetVersionedManuscript(entityID string) (VersionedManuscript, error)
}

type Server struct {
	Repo              Repo
	Aggregate         go_piggy.Aggregate
	EntityIdGenerator func() string
	handler           http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func NewServer(repo Repo, aggregate go_piggy.Aggregate, options ...func(*Server)) *Server {
	s := new(Server)
	s.Repo = repo
	s.Aggregate = aggregate
	s.EntityIdGenerator = go_piggy.RandomID

	for _, op := range options {
		op(s)
	}

	r := mux.NewRouter()

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	r.HandleFunc("/manuscripts/{entityID}", s.getManuscriptJSON)
	r.HandleFunc("/manuscripts", s.createManuscript).Methods("POST")
	r.HandleFunc("/manuscripts/{entityID}/events", s.sendCommands).Methods("POST")
	r.HandleFunc("/manuscripts/{entityID}/events", s.showEvents).Methods("GET")

	s.handler = cors(r)

	return s
}

func WithEntityIdGenerator(f func() string) func(*Server) {
	return func(server *Server) {
		server.EntityIdGenerator = f
	}
}

func (s *Server) sendCommands(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)

	var facts []go_piggy.Fact

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(body, &facts)

	accepted := s.Aggregate.ProcessCommand(NewManuscriptVersionEvent(entityID, facts...))

	if accepted {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func (s *Server) createManuscript(w http.ResponseWriter, r *http.Request) {
	newEntityID := s.EntityIdGenerator()

	s.Aggregate.ProcessCommand(NewManuscriptEvent(Manuscript{
		EntityID: newEntityID,
	}))

	location := "/manuscripts/" + newEntityID
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

//todo: handle errors
func (s *Server) getManuscriptJSON(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)
	version := r.URL.Query().Get("version")

	var manuscript Manuscript

	m, err := s.Repo.GetVersionedManuscript(entityID)

	if err != nil {
		//todo: 404 for any error probably isnt quite right!
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if version != "" {
		v, _ := strconv.Atoi(version)
		manuscript, _ = m.Version(v)
	} else {
		manuscript = m.CurrentRevision()
	}

	manuscriptJSON, _ := json.Marshal(manuscript)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(manuscriptJSON)
}

func entityIDFromRequest(r *http.Request) string {
	return mux.Vars(r)["entityID"]
}
