package manuscript

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/quii/go-piggy"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

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

func NewServer(repo Repo, emitter go_piggy.Emitter, options ...func(*Server)) *Server {
	s := new(Server)
	s.Repo = repo
	s.Emitter = emitter
	s.EntityIdGenerator = go_piggy.RandomID

	for _, op := range options {
		op(s)
	}

	staticEditor := http.FileServer(http.Dir("manuscript/editor"))

	r := mux.NewRouter()
	r.HandleFunc("/manuscripts/{entityID}", s.getManuscriptJSON).Headers("accept", "application/json")
	r.HandleFunc("/manuscripts/{entityID}", s.manuscriptForm)

	r.HandleFunc("/manuscripts", s.createManuscript).Methods("POST")

	r.HandleFunc("/manuscripts/{entityID}/events", s.addEventsToManuscript).Methods("POST")

	r.Handle("/", staticEditor)

	s.handler = handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Location"}),
	)(r)

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

	if r.Header.Get("content-type") == "application/json" {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(body, &facts)
		s.Emitter.Send(NewManuscriptChangesEvent(Manuscript{EntityID: entityID}, facts...))
		w.WriteHeader(http.StatusAccepted)

	} else {
		r.ParseForm()
		log.Println(r.Form)
		facts = append(facts, go_piggy.Fact{Op: "SET", Key: "Title", Value: r.Form.Get("title")})
		facts = append(facts, go_piggy.Fact{Op: "SET", Key: "Abstract", Value: r.Form.Get("abstract")})
		s.Emitter.Send(NewManuscriptChangesEvent(Manuscript{EntityID: entityID}, facts...))
		w.Header().Add("location", "/manuscripts/"+entityID)
		w.WriteHeader(http.StatusSeeOther)
	}

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

func (s *Server) manuscriptForm(w http.ResponseWriter, r *http.Request) {
	entityID := entityIDFromRequest(r)
	manuscript := s.Repo.GetManuscript(entityID)

	t, err := template.ParseFiles("manuscript/editor/editor.html")
	if err != nil {
		log.Fatal("Problem parsing template", err)
	}

	t.Execute(w, manuscript)
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
