package manuscript

import (
	"github.com/quii/go-piggy"
	"log"
	"regexp"
)

type Projection struct {
	receiver             go_piggy.Receiver
	versionedManuscripts VersionedManuscripts
}

func NewProjection(receiver go_piggy.Receiver) (m *Projection) {

	m = new(Projection)
	m.receiver = receiver
	m.versionedManuscripts = make(VersionedManuscripts)

	go m.listenForUpdates()

	return
}

// are these two functions really needed? perhaps a simpler design exists
func (p *Projection) GetManuscript(entityID string) Manuscript {
	return p.versionedManuscripts.CurrentRevision(entityID)
}

func (p *Projection) GetVersionedManuscript(entityID string, version int) (Manuscript, error) {
	return p.versionedManuscripts.GetVersionedManuscript(entityID, version)
}

func (p *Projection) listenForUpdates() {
	events := p.receiver.Listen(0)

	for event := range events {
		manuscript := p.versionedManuscripts.CurrentRevision(event.ID)
		log.Println("Got event", event)
		p.versionedManuscripts[event.ID] = append(p.versionedManuscripts[event.ID], manuscript.Update(event.Facts))
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
