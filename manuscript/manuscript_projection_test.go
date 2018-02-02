package manuscript

import (
	"github.com/quii/go-piggy"
	"github.com/stretchr/testify/assert"
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

	changes := make(chan int, 1)
	projection := NewProjection(eventSource, &ProjectionOptions{VersionChanges: changes})

	waitForManuscriptVersion(t, 2, changes)

	assert.Len(t, projection.versionedManuscripts, 2)
	assert.Contains(t, projection.versionedManuscripts, manuscript1.EntityID)
	assert.Contains(t, projection.versionedManuscripts, manuscript2.EntityID)
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

	changes := make(chan int, 1)
	projection := NewProjection(eventSource, &ProjectionOptions{VersionChanges: changes})

	waitForManuscriptVersion(t, 1, changes)

	parsedManuscript := projection.currentRevisionOrDefault(manuscript.EntityID)

	expectedManuscript := manuscript
	expectedManuscript.Version = 1

	assert.Equal(t, expectedManuscript, parsedManuscript)
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
	eventSource.Send(NewManuscriptVersionEvent(man1, TitleChanged("Showered and blue blazered")))
	eventSource.Send(NewManuscriptVersionEvent(man2, TitleChanged("Fill yourself with quarters")))
	eventSource.Send(NewManuscriptVersionEvent(man2, TitleChanged("lol")))

	changes := make(chan int, 10)
	projection := NewProjection(eventSource, &ProjectionOptions{VersionChanges: changes})

	waitForManuscriptVersion(t, 5, changes)

	expectedMan1State := Manuscript{
		EntityID: man1.EntityID,
		Title:    "Showered and blue blazered",
		Version:  2,
	}

	expectedMan2State := Manuscript{
		EntityID: man2.EntityID,
		Title:    "lol",
		Version:  3,
	}

	assert.Equal(t, expectedMan1State, projection.currentRevisionOrDefault(man1.EntityID))
	assert.Equal(t, expectedMan2State, projection.currentRevisionOrDefault(man2.EntityID))
}

func TestAuthorEventsAreProjectedAcrossMultipleEvents(t *testing.T) {
	manuscript := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}

	eventSource.Send(NewManuscriptEvent(manuscript))
	eventSource.Send(NewManuscriptVersionEvent(manuscript, AuthorsSet(0, "CJ")))
	eventSource.Send(NewManuscriptVersionEvent(manuscript, AuthorsSet(1, "TV")))

	changes := make(chan int, 1)
	projection := NewProjection(eventSource, &ProjectionOptions{VersionChanges: changes})

	waitForManuscriptVersion(t, 3, changes)

	parsedManuscript := projection.currentRevisionOrDefault(manuscript.EntityID)

	assert.Equal(t, parsedManuscript.Authors[0], "CJ")
	assert.Equal(t, parsedManuscript.Authors[1], "TV")
}

func (p *Projection) manuscriptExists(id string) bool {
	_, exists := p.versionedManuscripts[id]
	return exists
}

func waitForManuscriptVersion(t *testing.T, version int, changes chan int) {
	for i := 0; i < version; i++ {
		select {
		case <-changes:
		case <-time.After(1 * time.Second):
			t.Fatal("timed out waiting for changes")
		}
	}
}
