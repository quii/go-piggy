package manuscript

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/quii/go-piggy"
	"io/ioutil"
	"net/http"
	"strconv"
)

// deffo WIP, just experimenting with how this should all hang together

type Repo interface {
	GetManuscript(id string) Manuscript
	GetVersionedManuscript(entityID string, version int) (Manuscript, error)
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

func WithEntityIdGenerator(f func() string) func(*Server) {
	return func(server *Server) {
		server.EntityIdGenerator = f
	}
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
	r.HandleFunc("/{entityID}", s.getManuscript)
	r.HandleFunc("/", s.createManuscript).Methods("POST")
	r.HandleFunc("/{entityID}/events", s.addEventsToManuscript).Methods("POST")

	s.handler = r

	return s
}

func (s *Server) addEventsToManuscript(w http.ResponseWriter, r *http.Request) {
	entityID := mux.Vars(r)["entityID"]

	var facts []go_piggy.Fact

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	json.Unmarshal(body, &facts)

	s.Emitter.Send(NewManuscriptChangesEvent(Manuscript{EntityID: entityID}, facts...))

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) createManuscript(w http.ResponseWriter, r *http.Request) {
	newEntityID := s.EntityIdGenerator()

	s.Emitter.Send(NewManuscriptEvent(Manuscript{
		EntityID: newEntityID,
	}))

	w.Header().Add("location", "/"+newEntityID)
	w.WriteHeader(http.StatusCreated)

}

func (s *Server) getManuscript(w http.ResponseWriter, r *http.Request) {
	entityID := mux.Vars(r)["entityID"]
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
