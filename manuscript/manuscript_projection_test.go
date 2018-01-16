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

	projection := NewProjection(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	assert.Len(t, projection.manuscripts, 2)
	assert.Contains(t, projection.manuscripts, manuscript1.EntityID)
	assert.Contains(t, projection.manuscripts, manuscript2.EntityID)
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

	projection := NewProjection(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, exists := projection.manuscripts[manuscript.EntityID]

	assert.True(t, exists)
	assert.Equal(t, parsedManuscript, manuscript)
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

	projection := NewProjection(eventSource)
	time.Sleep(5 * time.Millisecond) //todo: bleh

	expectedMan1State := Manuscript{
		EntityID: man1.EntityID,
		Title:    "Showered and blue blazered",
	}

	expectedMan2State := Manuscript{
		EntityID: man2.EntityID,
		Title:    "Fill yourself with quarters",
	}

	assert.Equal(t, projection.manuscripts[man1.EntityID], expectedMan1State)
	assert.Equal(t, projection.manuscripts[man2.EntityID], expectedMan2State)
}

func TestAuthorEventsAreProjectedAcrossMultipleEvents(t *testing.T) {
	manuscript := Manuscript{
		EntityID: go_piggy.RandomID(),
	}

	eventSource := &go_piggy.InMemorySource{}

	eventSource.Send(NewManuscriptEvent(manuscript))
	eventSource.Send(NewManuscriptChangesEvent(manuscript, AuthorsSet(0, "CJ")))
	eventSource.Send(NewManuscriptChangesEvent(manuscript, AuthorsSet(1, "TV")))

	projection := NewProjection(eventSource)

	time.Sleep(5 * time.Millisecond) //todo: bleh

	parsedManuscript, _ := projection.manuscripts[manuscript.EntityID]

	assert.Equal(t, parsedManuscript.Authors[0], "CJ")
	assert.Equal(t, parsedManuscript.Authors[1], "TV")
}

func (p *Projection) manuscriptExists(id string) bool {
	_, exists := p.manuscripts[id]
	return exists
}
