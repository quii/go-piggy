package manuscript

import (
	"github.com/quii/go-piggy"
	"log"
	"regexp"
)

type Projection struct {
	receiver    go_piggy.Receiver
	manuscripts map[string]Manuscript
}

func NewProjection(receiver go_piggy.Receiver) (m *Projection) {

	m = new(Projection)
	m.receiver = receiver
	m.manuscripts = make(map[string]Manuscript)

	go m.listenForUpdates()

	return
}

func (p *Projection) GetManuscript(id string) Manuscript {
	man, _ := p.manuscripts[id]
	return man
}

func (p *Projection) listenForUpdates() {
	events := p.receiver.Listen(0)

	for event := range events {
		log.Printf("got a new event %+v\n", event)
		manuscript, _ := p.manuscripts[event.ID]
		manuscript.EntityID = event.ID
		p.manuscripts[event.ID] = manuscript.Update(event.Facts)
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
