package manuscript

import (
	"github.com/quii/go-piggy"
	"regexp"
	"strconv"
)

type Repo struct {
	receiver    go_piggy.Receiver
	manuscripts map[string]Manuscript
}

func NewRepo(receiver go_piggy.Receiver) (m *Repo) {

	m = new(Repo)
	m.receiver = receiver
	m.manuscripts = make(map[string]Manuscript)

	go m.listenForUpdates()

	return
}

func (r *Repo) listenForUpdates() {
	events := r.receiver.Listen(0)

	for event := range events {
		manuscript, _ := r.manuscripts[event.ID]
		manuscript.EntityID = event.ID
		r.manuscripts[event.ID] = updateManuscript(manuscript, event.Facts)
	}
}

var authorsRegex = regexp.MustCompile(`Authors\[\d+\]`)
var authorIndexRegex = regexp.MustCompile(`(\d+)`) //todo: i suck at regex and it feels dangerous just to match the first number

func updateManuscript(m Manuscript, facts []go_piggy.Fact) Manuscript {

	for _, f := range facts {
		switch f.Op {
		case "SET":
			switch f.Key {
			case "Title":
				m.Title = f.Value
			case "Abstract":
				m.Abstract = f.Value
			}

			if authorsRegex.MatchString(f.Key) {
				extractedIndex := authorIndexRegex.FindString(f.Key)
				i, _ := strconv.Atoi(extractedIndex)
				m.InsertAuthorIn(i, f.Value)
			}
		}
	}

	return m
}
