package manuscript

import (
	"github.com/quii/go-piggy"
	"regexp"
	"fmt"
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

func (p *Projection) GetManuscript(entityID string) Manuscript {
	return p.versionedManuscripts.CurrentRevision(entityID)
}

//todo: testme
func (p *Projection) GetVersionManuscriptVersion(entityID string, version int) (Manuscript, error) {
	versions, exists := p.versionedManuscripts[entityID]

	if !exists{
		return Manuscript{}, fmt.Errorf("manuscript %s does not exist", entityID)
	}

	if len(versions) < version{
		return Manuscript{}, fmt.Errorf("manuscript version %d of %s does not exist", version, entityID)
	}

	return versions[version], nil
}

func (p *Projection) listenForUpdates() {
	events := p.receiver.Listen(0)

	for event := range events {
		manuscript := p.versionedManuscripts.CurrentRevision(event.ID)
		p.versionedManuscripts[event.ID] = append(p.versionedManuscripts[event.ID], manuscript.Update(event.Facts))
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number
