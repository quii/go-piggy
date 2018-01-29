package manuscript

import (
	"github.com/quii/go-piggy"
	"regexp"
)

type Projection struct {
	receiver             go_piggy.Receiver
	versionedManuscripts VersionedManuscripts
	events               map[string][]go_piggy.Event
	changes              chan int
	version              int
}

//todo: changes should be provided by consumer as an optional thing
func NewProjection(receiver go_piggy.Receiver) (m *Projection, changes chan int) {

	m = new(Projection)
	m.receiver = receiver
	m.changes = make(chan int, 1)
	m.versionedManuscripts = make(VersionedManuscripts)
	m.events = make(map[string][]go_piggy.Event)

	go m.listenForUpdates()

	return m, m.changes
}

//todo: should probably be elsewhere
func (p *Projection) Events(entityID string) []go_piggy.Event {
	events, _ := p.events[entityID]
	return events
}

// are these two functions really needed? perhaps a simpler design exists
func (p *Projection) GetManuscript(entityID string) Manuscript {
	return p.versionedManuscripts.CurrentRevision(entityID)
}

func (p *Projection) GetVersionedManuscript(entityID string, version int) (Manuscript, error) {
	return p.versionedManuscripts.GetVersionedManuscript(entityID, version)
}

func (p *Projection) Versions(entityID string) int {
	return p.versionedManuscripts.Versions(entityID)
}

func (p *Projection) listenForUpdates() {
	events := p.receiver.Listen(0)

	for event := range events {
		manuscript := p.versionedManuscripts.CurrentRevision(event.ID)
		pastEvents, _ := p.events[event.ID]
		p.events[event.ID] = append(pastEvents, event)
		p.versionedManuscripts[event.ID] = append(p.versionedManuscripts[event.ID], manuscript.Update(event.Facts))

		p.version++
		p.changes <- p.version
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
