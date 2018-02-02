package manuscript

import (
	"github.com/quii/go-piggy"
	"regexp"
	"log"
)

type Projection struct {
	receiver             go_piggy.Receiver
	versionedManuscripts VersionedManuscripts
	events               map[string][]go_piggy.Event
	changes              chan int
	version              int
	options              *ProjectionOptions
}

type ProjectionOptions struct {
	VersionChanges chan int
}

func NewProjection(receiver go_piggy.Receiver, options *ProjectionOptions) (m *Projection) {

	m = new(Projection)
	m.receiver = receiver
	m.changes = make(chan int, 1)
	m.versionedManuscripts = make(VersionedManuscripts)
	m.events = make(map[string][]go_piggy.Event)
	m.options = options

	go m.listenForUpdates()

	return m
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
		log.Printf("Event received %s\n", event)
		manuscript := p.versionedManuscripts.CurrentRevision(event.EntityID)
		pastEvents, _ := p.events[event.EntityID]
		p.events[event.EntityID] = append(pastEvents, event)
		p.versionedManuscripts[event.EntityID] = append(p.versionedManuscripts[event.EntityID], manuscript.Update(event.Facts))

		p.version++
		if p.options != nil {
			p.options.VersionChanges <- p.version
		}
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
