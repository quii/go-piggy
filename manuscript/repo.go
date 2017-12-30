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
	m := Manuscript{}
	
	for _, f := range facts{
		switch f.Op {
		case "SET":
			switch f.Key {
			case "Title":
				m.Title = f.Value
			case "Abstract":
				m.Abstract = f.Value
			}
		}
	}
	
	return m
}

