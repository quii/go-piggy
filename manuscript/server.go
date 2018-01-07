package manuscript

import (
	"encoding/json"
	"net/http"
	"strings"
)

// deffo WIP, just experimenting with how this should all hang together

type ManuscriptRepo interface {
	CreateManuscript(id string)
	GetManuscript(id string) Manuscript
}

type Server struct {
	Repo              ManuscriptRepo
	EntityIdGenerator func() string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		newEntityID := s.EntityIdGenerator()
		s.Repo.CreateManuscript(newEntityID)

		w.Header().Add("location", "/"+newEntityID)
		w.WriteHeader(http.StatusCreated)
	}

	if r.Method == http.MethodGet {
		entityId := strings.TrimPrefix(r.URL.Path, "/")
		manuscript := s.Repo.GetManuscript(entityId)
		manuscriptJSON, _ := json.Marshal(manuscript)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(manuscriptJSON)
		w.WriteHeader(http.StatusOK)
	}
}
