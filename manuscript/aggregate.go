package manuscript

import "github.com/quii/go-piggy"

type AggregateProjection interface {
	GetVersionedManuscript(entityID string) (VersionedManuscript, error)
}

type Aggregate struct {
	projection AggregateProjection
	emitter    go_piggy.Emitter
}

func NewAggregate(projection AggregateProjection, emitter go_piggy.Emitter) *Aggregate {
	return &Aggregate{
		projection: projection,
		emitter:    emitter,
	}
}

func (a *Aggregate) ProcessCommand(event go_piggy.Event) (accepted bool) {
	m, err := a.projection.GetVersionedManuscript(event.EntityID)

	if err != nil {
		return false
	}

	manuscript := m.CurrentRevision()

	if manuscript.Published {
		return false
	}

	a.emitter.Send(event)

	return true
}
