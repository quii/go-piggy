package manuscript

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	if response.Code != http.StatusCreated {
		t.Errorf("expected a status created but got %d", response.Code)
	}

	if response.Header().Get("Location") != "/random-id" {
		t.Errorf("did not get a location header pointing to new document, headers were %+v", response.Header())
	}

	if repo.createManuscriptRaised != "random-id" {
		t.Errorf("event source did not have a new document raised with random-id")
	}
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

	if response.Code != http.StatusOK {
		t.Errorf("expected a 200 when fetching document but got %d", response.Code)
	}

	receivedManuscript := Manuscript{}
	err := json.Unmarshal(response.Body.Bytes(), &receivedManuscript)

	if err != nil {
		t.Fatalf("problem parsing manuscript from response %s %+v", response.Body, err)
	}

	if !reflect.DeepEqual(receivedManuscript, manuscript) {
		t.Errorf("manuscript returned is wrong expected %+v got %+v", manuscript, receivedManuscript)
	}
}
