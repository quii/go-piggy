package manuscript

import (
	"testing"
	"github.com/quii/go-piggy"
	"time"
)

func TestItAddsManuscriptsAsTheyAreAdded(t *testing.T) {

	manuscript1 := go_piggy.NewEvent("manuscript", nil)
	manuscript2 := go_piggy.NewEvent("manuscript", nil)

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(manuscript1)
	eventSource.Send(manuscript2)

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	if len(repo.manuscripts) != 2 {
		t.Errorf("Repo has not processed 2 manuscripts, it has done %d", len(repo.manuscripts))
	}

	if !repo.manuscriptExists(manuscript1.ID) {
		t.Errorf("Could not find manuscript 1 in repo %s", repo.manuscripts)
	}

	if !repo.manuscriptExists(manuscript2.ID) {
		t.Errorf("Could not find manuscript 2 in repo %s", repo.manuscripts)
	}

	if repo.manuscriptExists("unknown") {
		t.Error("Should not be able to find unknown manuscript in repo!")
	}
}

func TestItReadsFactsIntoManuscripts(t *testing.T) {
	manuscript := go_piggy.NewEvent("manuscript", []go_piggy.Fact{
		{"SET", "Title", "Hello, world"},
		{"SET", "Abstract", "the catcher in the rye"},
	})

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(manuscript)

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, exists := repo.manuscripts[manuscript.ID]

	if !exists{
		t.Error("The manuscript was not saved after event was sent")
	}

	if parsedManuscript.Title != "Hello, world" {
		t.Errorf("manuscript title is incorrect, expect Hello, world but got %s", parsedManuscript.Title)
	}

	if parsedManuscript.Abstract != "the catcher in the rye" {
		t.Errorf("manuscript abstract is incorrect, expect catcher in the rye byt got %s", parsedManuscript.Abstract)
	}
}

func (r *Repo) manuscriptExists(id string) bool {
	_, exists := r.manuscripts[id]
	return exists
}