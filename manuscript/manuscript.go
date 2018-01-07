package manuscript

type Manuscript struct {
	EntityID        string
	Title, Abstract string
	Authors         []string
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
