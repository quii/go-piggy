package manuscript

import (
	"github.com/quii/go-piggy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventRejectAndAccept(t *testing.T) {

	emitter := go_piggy.NewEmitterSpy()

	aggregateProjection := &fakeManuscriptRepo{
		manuscripts: VersionedManuscripts{
			"A": []Manuscript{{EntityID: "A", Title: "Lovely Time", Published: true}},
			"B": []Manuscript{{EntityID: "B", Title: "Lovely Time", Published: false}},
		},
	}

	aggregate := NewAggregate(aggregateProjection, emitter)

	t.Run("it rejects commands updating a manuscript which is already published", func(t *testing.T) {
		event := NewManuscriptVersionEvent("A", manuscriptChange)
		accepted := aggregate.ProcessCommand(event)
		assert.False(t, accepted)
		assert.False(t, emitter.EventReceived(event))
	})

	t.Run("it accepts commands when a manuscript is not yet published", func(t *testing.T) {
		event := NewManuscriptVersionEvent("B", manuscriptChange)
		accepted := aggregate.ProcessCommand(event)
		assert.True(t, accepted)
		assert.True(t, emitter.EventReceived(event))
	})

	t.Run("it rejects a new version command if the projection cant get the manuscript", func(t *testing.T) {
		event := NewManuscriptVersionEvent("NotFound", manuscriptChange)
		accepted := aggregate.ProcessCommand(event)
		assert.False(t, accepted)
		assert.False(t, emitter.EventReceived(event))
	})
}

var manuscriptChange = TitleChanged("new title yeah")
