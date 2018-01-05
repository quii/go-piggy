package manuscript

import (
	"github.com/quii/go-piggy"
	"testing"
	"time"
	"reflect"
)

func TestItAddsManuscriptsAsTheyAreAdded(t *testing.T) {

	manuscript1 := Manuscript{
		entityID: go_piggy.RandomID(),
	}
	manuscript2 := Manuscript{
		entityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(NewManuscriptEvent(manuscript1))
	eventSource.Send(NewManuscriptEvent(manuscript2))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	if len(repo.manuscripts) != 2 {
		t.Errorf("Repo has not processed 2 manuscripts, it has done %d", len(repo.manuscripts))
	}

	if !repo.manuscriptExists(manuscript1.entityID) {
		t.Errorf("Could not find manuscript 1 in repo %s", repo.manuscripts)
	}

	if !repo.manuscriptExists(manuscript2.entityID) {
		t.Errorf("Could not find manuscript 2 in repo %s", repo.manuscripts)
	}

	if repo.manuscriptExists("unknown") {
		t.Error("Should not be able to find unknown manuscript in repo!")
	}
}

func TestItReadsNewManuscriptEvent(t *testing.T) {

	manuscript := Manuscript{
		entityID: go_piggy.RandomID(),
		Title:    "Hello, world",
		Abstract: "the catcher in the rye",
		Authors:  []string{"CJ", "TS"},
	}

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(NewManuscriptEvent(manuscript))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, exists := repo.manuscripts[manuscript.entityID]

	if !exists {
		t.Error("The manuscript was not saved after event was sent")
	}

	if !reflect.DeepEqual(parsedManuscript, manuscript) {
		t.Fatalf("The manuscript sent in is not the same as the one that was parsed from facts expected: %+v actual: %+v", manuscript, parsedManuscript)
	}
}

func TestAuthorEventsAreProjectedAcrossMultipleEvents(t *testing.T) {
	manuscript := Manuscript{
		entityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}

	eventSource.Send(NewManuscriptEvent(manuscript))
	eventSource.Send(NewManuscriptChanges(manuscript, AuthorsSet(0, "CJ")))
	eventSource.Send(NewManuscriptChanges(manuscript, AuthorsSet(1, "TV")))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, _ := repo.manuscripts[manuscript.entityID]

	if parsedManuscript.Authors[0] != "CJ" {
		t.Errorf("authors not set correctly, expect CJ at [0] %s", parsedManuscript.Authors)
	}

	if parsedManuscript.Authors[1] != "TV" {
		t.Errorf("authors not set correctly, expect TV at [1] %s", parsedManuscript.Authors)
	}
}

func (r *Repo) manuscriptExists(id string) bool {
	_, exists := r.manuscripts[id]
	return exists
}
