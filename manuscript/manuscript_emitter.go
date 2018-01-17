package manuscript

import "github.com/quii/go-piggy"

type Emitter struct {
	emitter go_piggy.Emitter
}

func NewEmitter(emitter go_piggy.Emitter) *Emitter {
	return &Emitter{
		emitter: emitter,
	}
}

func (p *Emitter) CreateManuscript(id string) {
	p.emitter.Send(NewManuscriptEvent(Manuscript{
		EntityID: id,
		Title:    "Title not set yet",
	}))
}
