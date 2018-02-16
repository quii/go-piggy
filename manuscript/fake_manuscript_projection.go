package manuscript

import (
	"fmt"
	"github.com/quii/go-piggy"
)

type fakeManuscriptRepo struct {
	manuscripts VersionedManuscripts
}

func (f *fakeManuscriptRepo) Events(entityID string) []go_piggy.Event {
	panic("implement me")
}

func (f *fakeManuscriptRepo) GetVersionedManuscript(entityID string) (VersionedManuscript, error) {
	manuscripts, exists := f.manuscripts[entityID]
	if !exists {
		return VersionedManuscript{}, fmt.Errorf("entityID %s does not exist", entityID)
	}
	return manuscripts, nil
}
