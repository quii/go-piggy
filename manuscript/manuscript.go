package manuscript

import (
	"github.com/quii/go-piggy"
	"strconv"
)

type Manuscript struct {
	EntityID        string
	Title, Abstract string
	Authors         []string
	Version         int
}

func (m *Manuscript) InsertAuthorIn(index int, name string) {

	if m.Authors == nil {
		m.Authors = make([]string, 0)
	}
	authorsArraySize := len(m.Authors)

	if (index + 1) > authorsArraySize {
		authorsArraySize = index + 1
	}

	newArray := make([]string, authorsArraySize)

	copy(newArray, m.Authors)

	newArray[index] = name

	m.Authors = newArray
}

func (m *Manuscript) Update(facts []go_piggy.Fact) Manuscript {
	newVersion := *m

	newVersion.Version++

	for _, f := range facts {
		switch f.Op {
		case "SET":
			switch f.Key {
			case "Title":
				newVersion.Title = f.Value
			case "Abstract":
				newVersion.Abstract = f.Value
			}

			if authorsRegex.MatchString(f.Key) {
				extractedIndex := authorIndexRegex.FindString(f.Key)
				i, _ := strconv.Atoi(extractedIndex)
				newVersion.InsertAuthorIn(i, f.Value)
			}
		}
	}

	return newVersion
}
