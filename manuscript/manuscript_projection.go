package manuscript

import (
	"fmt"
	"github.com/quii/go-piggy"
	"regexp"
	"time"
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
	Logger         go_piggy.Logger
}

func readyOptions(userOptions *ProjectionOptions) *ProjectionOptions {
	if userOptions == nil {
		userOptions = &ProjectionOptions{}
	}

	if userOptions.Logger == nil {
		userOptions.Logger = go_piggy.NewStdoutLogger()
	}

	return userOptions
}

func NewProjection(receiver go_piggy.Receiver, options *ProjectionOptions) (m *Projection) {

	m = new(Projection)
	m.receiver = receiver
	m.changes = make(chan int, 1)
	m.versionedManuscripts = make(VersionedManuscripts)
	m.events = make(map[string][]go_piggy.Event)
	m.options = readyOptions(options)

	go m.listenForUpdates()

	return m
}

func (p *Projection) Events(entityID string) []go_piggy.Event {
	events, _ := p.events[entityID]
	return events
}

func (p *Projection) GetVersionedManuscript(entityID string) (VersionedManuscript, error) {
	manuscripts, exists := p.versionedManuscripts[entityID]

	if !exists {
		return VersionedManuscript{}, fmt.Errorf("manuscript %s does not exist", entityID)
	}

	return manuscripts, nil
}

func (p *Projection) currentRevisionOrDefault(entityID string) Manuscript {
	versions, err := p.GetVersionedManuscript(entityID)

	if err != nil {
		return Manuscript{EntityID: entityID}
	}

	return versions.CurrentRevision()
}

func (p *Projection) listenForUpdates() {
	events := p.receiver.Listen(0)

	for event := range events {
		p.options.Logger.Info(fmt.Sprintf("got event %+v", event))

		manuscript := p.currentRevisionOrDefault(event.EntityID)

		pastEvents, _ := p.events[event.EntityID]
		p.events[event.EntityID] = append(pastEvents, event)
		p.versionedManuscripts[event.EntityID] = append(p.versionedManuscripts[event.EntityID], manuscript.Update(event.Facts))

		go p.incrementVersion()
	}
}

func (p *Projection) incrementVersion() {
	p.version++

	if p.options.VersionChanges != nil {
		select {
		case p.options.VersionChanges <- p.version:
		case <-time.After(1 * time.Second):
		}
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
