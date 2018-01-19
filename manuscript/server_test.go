package manuscript

import (
	"encoding/json"
	"fmt"
	"github.com/quii/go-piggy"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeManuscriptRepo struct {
	manuscripts map[string]Manuscript
}

func (f *fakeManuscriptRepo) GetManuscript(id string) Manuscript {
	man, _ := f.manuscripts[id]
	return man
}

func (f *fakeManuscriptRepo) GetVersionedManuscript(entityID string, version int) (Manuscript, error) {
	return Manuscript{}, fmt.Errorf("haven't made this yet...")
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
		manuscripts: map[string]Manuscript{
			"random-id": manuscript,
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
