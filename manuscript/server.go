package manuscript

import (
	"encoding/json"
	"github.com/quii/go-piggy"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	entityId := strings.TrimPrefix(r.URL.Path, "/")

	if r.Method == http.MethodPost && entityId == "" {
		newEntityID := s.EntityIdGenerator()

		s.Emitter.Send(NewManuscriptEvent(Manuscript{
			EntityID: newEntityID,
		}))

		w.Header().Add("location", "/"+newEntityID)
		w.WriteHeader(http.StatusCreated)

		//todo: test me!
	} else if r.Method == http.MethodPost {
		var facts []go_piggy.Fact

		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		json.Unmarshal(body, &facts)

		log.Println("got some facts", facts)

		s.Emitter.Send(NewManuscriptChangesEvent(Manuscript{EntityID: entityId}, facts...))
	}

	if r.Method == http.MethodGet {
		version := r.URL.Query().Get("version")

		var manuscript Manuscript

		//todo: test this version stuff
		if version != "" {
			v, _ := strconv.Atoi(version)
			log.Println("getting version", v)
			m, _ := s.Repo.GetVersionedManuscript(entityId, v)
			manuscript = m
		} else {
			manuscript = s.Repo.GetManuscript(entityId)
		}

		manuscriptJSON, _ := json.Marshal(manuscript)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(manuscriptJSON)
	}

}
