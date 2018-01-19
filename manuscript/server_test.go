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

func TestItRaisesNewManuscriptEventOnPost(t *testing.T) {
	emitter := &fakeEmitter{}

	server := NewServer(
		&fakeManuscriptRepo{},
		emitter,
		withFixedEntityId("random-id"),
	)

	request, _ := http.NewRequest(http.MethodPost, "/manuscripts", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "/manuscripts/random-id", response.Header().Get("location"))
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

	server := NewServer(
		repo,
		&fakeEmitter{},
		withFixedEntityId("random-id"),
	)

	request, _ := http.NewRequest(http.MethodGet, "/manuscripts/random-id", nil)
	request.Header.Add("accept", "application/json")
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

	server := NewServer(
		repo,
		&fakeEmitter{},
		WithEntityIdGenerator(func() string {
			return "random-id"
		}),
	)

	request, _ := http.NewRequest(http.MethodGet, "/manuscripts/random-id?version=2", nil)
	request.Header.Add("accept", "application/json")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	receivedManuscript := Manuscript{}
	err := json.Unmarshal(response.Body.Bytes(), &receivedManuscript)

	assert.NoError(t, err)
	assert.Equal(t, manuscriptV2, receivedManuscript)
}

func TestIt404sForManuscriptsThatDontExist(t *testing.T) {
	server := NewServer(
		&fakeManuscriptRepo{},
		&fakeEmitter{},
	)

	request, _ := http.NewRequest(http.MethodGet, "/manuscripts/random-id?version=2", nil)
	request.Header.Add("accept", "application/json")

	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestItAddsEventsToExistingManuscripts(t *testing.T) {
	emitter := &fakeEmitter{}

	server := NewServer(
		&fakeManuscriptRepo{},
		emitter,
	)

	eventJSON := `[
		{"OP":"SET", "Key":"Title", "Value": "Bob"},
		{"OP":"SET", "Key":"Abstract", "Value": "Smith"}
	]`

	request, _ := http.NewRequest(http.MethodPost, "/manuscripts/random-id/events", strings.NewReader(eventJSON))
	request.Header.Set("content-type", "application/json")
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

type fakeManuscriptRepo struct {
	manuscripts VersionedManuscripts
}

func (f *fakeManuscriptRepo) Versions(entityID string) int {
	panic("implement me")
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

func withFixedEntityId(id string) func(*Server) {
	return WithEntityIdGenerator(func() string {
		return id
	})
}
