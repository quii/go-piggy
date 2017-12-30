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

func (r *Repo) manuscriptExists(id string) bool {
	_, exists := r.manuscripts[id]
	return exists
}