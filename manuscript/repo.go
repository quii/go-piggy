package manuscript

import "github.com/quii/go-piggy"

type Repo struct{
	receiver go_piggy.Receiver
	manuscripts map[string]Manuscript
}

func NewRepo(receiver go_piggy.Receiver) (m *Repo) {

	m = new(Repo)
	m.receiver = receiver
	m.manuscripts = make(map[string]Manuscript)

	go m.listenForUpdates()

	return
}

func (m *Repo) listenForUpdates() {
	events := m.receiver.Listen(0)

	for event := range events {

		if _, exists := m.manuscripts[event.ID]; !exists{
			m.manuscripts[event.ID] = newManuscriptFromEvent(event.Facts)
		} else {
			//todo: update manuscript from facts
		}
	}
}

func newManuscriptFromEvent(facts []go_piggy.Fact) Manuscript {
	//todo: make me work properly!
	return Manuscript{}
}

