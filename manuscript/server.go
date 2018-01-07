package manuscript

import (
	"encoding/json"
	"net/http"
	"strings"
)

// deffo WIP, just experimenting with how this should all hang together

type eventSource interface {
	CreateManuscript(id string)
}

type manuscriptRepo interface {
	GetManuscript(id string) Manuscript
}

type Server struct {
	eventSource       eventSource
	manuscriptRepo    manuscriptRepo
	entityIdGenerator func() string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		newEntityID := s.entityIdGenerator()
		s.eventSource.CreateManuscript(newEntityID)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("location", "/"+newEntityID)
	}

	if r.Method == http.MethodGet {
		entityId := strings.TrimPrefix(r.URL.Path, "/")
		manuscript := s.manuscriptRepo.GetManuscript(entityId)
		manuscriptJSON, _ := json.Marshal(manuscript)

		w.Write(manuscriptJSON)
		w.WriteHeader(http.StatusOK)
	}
}
