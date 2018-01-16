package manuscript

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeManuscriptRepo struct {
	createManuscriptRaised string
	manuscripts            map[string]Manuscript
}

func (f *fakeManuscriptRepo) CreateManuscript(id string) {
	f.createManuscriptRaised = id
}

func (f *fakeManuscriptRepo) GetManuscript(id string) Manuscript {
	man, _ := f.manuscripts[id]
	return man
}

func TestItRaisesNewManuscriptEventOnPost(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	repo := &fakeManuscriptRepo{}

	server := Server{
		Repo: repo,
		EntityIdGenerator: func() string {
			return "random-id"
		},
	}

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "/random-id", response.Header().Get("location"))
	assert.Equal(t, "random-id", repo.createManuscriptRaised)
}

func TestItGetsManuscripts(t *testing.T) {
	manuscript := Manuscript{
		EntityID: "random-id",
		Title:    "Pepper pot",
		Abstract: "Is a cat from egypt",
		Authors:  nil,
	}

	repo := &fakeManuscriptRepo{
		manuscripts: map[string]Manuscript{
			"random-id": manuscript,
		},
	}

	server := Server{
		Repo: repo,
	}

	request, _ := http.NewRequest(http.MethodGet, "/random-id", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	receivedManuscript := Manuscript{}
	err := json.Unmarshal(response.Body.Bytes(), &receivedManuscript)

	assert.NoError(t, err)
	assert.Equal(t, manuscript, receivedManuscript)
}
