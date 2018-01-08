package manuscript

import (
	"github.com/quii/go-piggy"
	"log"
	"regexp"
	"strconv"
)

type Projection struct {
	receiver    go_piggy.Receiver
	emitter     go_piggy.Emitter
	manuscripts map[string]Manuscript
}

//todo: should a thing ever use both? is it really just one interface? or too many concerns?
func NewProjection(eventSource go_piggy.EventSource) (m *Projection) {

	m = new(Projection)
	m.receiver = eventSource
	m.emitter = eventSource
	m.manuscripts = make(map[string]Manuscript)

	go m.listenForUpdates()

	return
}

func (p *Projection) CreateManuscript(id string) {
	p.emitter.Send(NewManuscriptEvent(Manuscript{
		EntityID: id,
		Title:    "Title not set yet",
	}))
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
		p.manuscripts[event.ID] = updateManuscript(manuscript, event.Facts)
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number

func updateManuscript(m Manuscript, facts []go_piggy.Fact) Manuscript {

	for _, f := range facts {
		switch f.Op {
		case "SET":
			switch f.Key {
			case "Title":
				m.Title = f.Value
			case "Abstract":
				m.Abstract = f.Value
			}

			if authorsRegex.MatchString(f.Key) {
				extractedIndex := authorIndexRegex.FindString(f.Key)
				i, _ := strconv.Atoi(extractedIndex)
				m.InsertAuthorIn(i, f.Value)
			}
		}
	}

	return m
}
