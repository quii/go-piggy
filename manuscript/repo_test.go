package manuscript

import (
	"github.com/quii/go-piggy"
	"reflect"
	"testing"
	"time"
)

func TestItAddsManuscriptsAsTheyAreAdded(t *testing.T) {

	manuscript1 := Manuscript{
		EntityID: go_piggy.RandomID(),
	}
	manuscript2 := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(NewManuscriptEvent(manuscript1))
	eventSource.Send(NewManuscriptEvent(manuscript2))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	if len(repo.manuscripts) != 2 {
		t.Errorf("Repo has not processed 2 manuscripts, it has done %d", len(repo.manuscripts))
	}

	if !repo.manuscriptExists(manuscript1.EntityID) {
		t.Errorf("Could not find manuscript 1 in repo %s", repo.manuscripts)
	}

	if !repo.manuscriptExists(manuscript2.EntityID) {
		t.Errorf("Could not find manuscript 2 in repo %s", repo.manuscripts)
	}

	if repo.manuscriptExists("unknown") {
		t.Error("Should not be able to find unknown manuscript in repo!")
	}
}

func TestItReadsNewManuscriptEvent(t *testing.T) {

	manuscript := Manuscript{
		EntityID: go_piggy.RandomID(),
		Title:    "Hello, world",
		Abstract: "the catcher in the rye",
		Authors:  []string{"CJ", "TS"},
	}

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(NewManuscriptEvent(manuscript))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, exists := repo.manuscripts[manuscript.EntityID]

	if !exists {
		t.Error("The manuscript was not saved after event was sent")
	}

	if !reflect.DeepEqual(parsedManuscript, manuscript) {
		t.Fatalf("The manuscript sent in is not the same as the one that was parsed from facts expected: %+v actual: %+v", manuscript, parsedManuscript)
	}
}

func TestItReadsFactsToTheCorrectManuscripts(t *testing.T) {
	man1 := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	man2 := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}
	eventSource.Send(NewManuscriptEvent(man1))
	eventSource.Send(NewManuscriptEvent(man2))
	eventSource.Send(NewManuscriptChangesEvent(man1, TitleChanged("Showered and blue blazered")))
	eventSource.Send(NewManuscriptChangesEvent(man2, TitleChanged("Fill yourself with quarters")))

	repo := NewRepo(eventSource)
	time.Sleep(5 * time.Millisecond) //todo: bleh

	expectedMan1State := Manuscript{
		EntityID: man1.EntityID,
		Title:    "Showered and blue blazered",
	}

	expectedMan2State := Manuscript{
		EntityID: man2.EntityID,
		Title:    "Fill yourself with quarters",
	}

	if actualMan1State, _ := repo.manuscripts[man1.EntityID]; !reflect.DeepEqual(actualMan1State, expectedMan1State) {
		t.Errorf("Man1 end state is not correct, expected %+v got %+v", expectedMan1State, actualMan1State)
	}

	if actualMan2State, _ := repo.manuscripts[man2.EntityID]; !reflect.DeepEqual(actualMan2State, expectedMan2State) {
		t.Errorf("Man1 end state is not correct, expected %+v got %+v", expectedMan1State, actualMan2State)
	}
}

func TestAuthorEventsAreProjectedAcrossMultipleEvents(t *testing.T) {
	manuscript := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}

	eventSource.Send(NewManuscriptEvent(manuscript))
	eventSource.Send(NewManuscriptChangesEvent(manuscript, AuthorsSet(0, "CJ")))
	eventSource.Send(NewManuscriptChangesEvent(manuscript, AuthorsSet(1, "TV")))

	repo := NewRepo(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, _ := repo.manuscripts[manuscript.EntityID]

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
