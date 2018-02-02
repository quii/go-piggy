package manuscript

import (
	"fmt"
)

type VersionedManuscript []Manuscript
type VersionedManuscripts map[string]VersionedManuscript

func (v VersionedManuscript) CurrentRevision() Manuscript {
	return v[len(v)-1]
}

func (v VersionedManuscript) Version(version int) (Manuscript, error) {
	if version > len(v) {
		return Manuscript{}, fmt.Errorf("the latest version number for %s is %d but you asked for %d", v.CurrentRevision().EntityID, len(v), version)
	}

	return v[version-1], nil
}
