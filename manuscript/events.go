package manuscript

import (
	"fmt"
	"github.com/quii/go-piggy"
)

const manuscriptEventType = "manuscript"

func NewManuscriptEvent(m Manuscript) go_piggy.Event {
	event := go_piggy.Event{
		Name:     "NewManuscript",
		EntityID: m.EntityID,
		Type:     manuscriptEventType,
		Facts: []go_piggy.Fact{
			TitleChanged(m.Title),
			AbstractChanged(m.Abstract),
		},
	}

	event.Facts = append(event.Facts, AllAuthorsSet(m.Authors)...)

	return event
}

func NewManuscriptVersionEvent(m Manuscript, facts ...go_piggy.Fact) go_piggy.Event {
	return go_piggy.Event{
		Name:     "NewManuscriptVersion",
		EntityID: m.EntityID,
		Type:     manuscriptEventType,
		Facts:    facts,
	}
}

func TitleChanged(value string) go_piggy.Fact {
	return go_piggy.Fact{Op: "SET", Key: "Title", Value: value}
}

func AbstractChanged(value string) go_piggy.Fact {
	return go_piggy.Fact{Op: "SET", Key: "Abstract", Value: value}
}

func AuthorsSet(position int, name string) go_piggy.Fact {
	return go_piggy.Fact{Op: "SET", Key: fmt.Sprintf("Authors[%d]", position), Value: name}
}

func AllAuthorsSet(authors []string) (facts []go_piggy.Fact) {
	for i, a := range authors {
		facts = append(facts, AuthorsSet(i, a))
	}
	return
}
