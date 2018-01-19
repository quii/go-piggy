package manuscript

import (
	"encoding/json"
	"github.com/quii/go-piggy"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeManuscriptRepo struct {
	manuscripts VersionedManuscripts
}

func (f *fakeManuscriptRepo) GetManuscript(id string) Manuscript {
	return f.manuscripts.CurrentRevision(id)
}

func (f *fakeManuscriptRepo) GetVersionedManuscript(entityID string, version int) (Manuscript, error) {
	return f.manuscripts.GetVersionedManuscript(entityID, version)
}

type fakeEmitter struct {
	events []go_piggy.Event
}

func (f *fakeEmitter) Send(event go_piggy.Event) {
	f.events = append(f.events, event)
}

func TestItRaisesNewManuscriptEventOnPost(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	repo := &fakeManuscriptRepo{}
	emitter := &fakeEmitter{}

	server := NewServer(
		repo,
		emitter,
		WithEntityIdGenerator(func() string {
			return "random-id"
		}),
	)

	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "/random-id", response.Header().Get("location"))

	assert.Contains(t, emitter.events, NewManuscriptEvent(Manuscript{
		EntityID: "random-id",
	}))
}

func TestItGetsManuscripts(t *testing.T) {
	manuscript := Manuscript{
		EntityID: "random-id",
		Title:    "Pepper pot",
		Abstract: "Is a cat from egypt",
		Authors:  nil,
	}

	repo := &fakeManuscriptRepo{
		manuscripts: map[string][]Manuscript{
			"random-id": {manuscript},
		},
	}

	emitter := &fakeEmitter{}

	server := NewServer(
		repo,
		emitter,
		WithEntityIdGenerator(func() string {
			return "random-id"
		}),
	)

	request, _ := http.NewRequest(http.MethodGet, "/random-id", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	receivedManuscript := Manuscript{}
	err := json.Unmarshal(response.Body.Bytes(), &receivedManuscript)

	assert.NoError(t, err)
	assert.Equal(t, manuscript, receivedManuscript)
}

func TestItGetsVersionedManuscripts(t *testing.T) {
	manuscriptV1 := Manuscript{
		EntityID: "random-id",
		Title:    "Pepper pot",
		Abstract: "Is a cat from egypt",
		Authors:  nil,
	}

	manuscriptV2 := manuscriptV1
	manuscriptV2.Title = "new title"

	repo := &fakeManuscriptRepo{
		manuscripts: map[string][]Manuscript{
			"random-id": {manuscriptV1, manuscriptV2},
		},
	}

	emitter := &fakeEmitter{}

	server := NewServer(
		repo,
		emitter,
		WithEntityIdGenerator(func() string {
			return "random-id"
		}),
	)

	request, _ := http.NewRequest(http.MethodGet, "/random-id?version=2", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	receivedManuscript := Manuscript{}
	err := json.Unmarshal(response.Body.Bytes(), &receivedManuscript)

	assert.NoError(t, err)
	assert.Equal(t, manuscriptV2, receivedManuscript)
}

func TestIt404sForVersionsThatDontExist(t *testing.T) {
	repo := &fakeManuscriptRepo{}
	emitter := &fakeEmitter{}

	server := NewServer(
		repo,
		emitter,
	)

	request, _ := http.NewRequest(http.MethodGet, "/random-id?version=2", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestItAddsEventsToExistingManuscripts(t *testing.T) {
	repo := &fakeManuscriptRepo{}
	emitter := &fakeEmitter{}

	server := NewServer(
		repo,
		emitter,
	)

	eventJSON := `[
		{"OP":"SET", "Key":"Title", "Value": "Bob"},
		{"OP":"SET", "Key":"Abstract", "Value": "Smith"}
	]`

	request, _ := http.NewRequest(http.MethodPost, "/random-id/events", strings.NewReader(eventJSON))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusAccepted, response.Code)
	assert.Contains(t, emitter.events, go_piggy.Event{
		ID:   "random-id",
		Type: "manuscript",
		Facts: []go_piggy.Fact{
			TitleChanged("Bob"),
			AbstractChanged("Smith"),
		},
	})
}
