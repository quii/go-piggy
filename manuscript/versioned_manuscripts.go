package manuscript

import "fmt"

type VersionedManuscripts map[string][]Manuscript

//todo: testme
func (v VersionedManuscripts) CurrentRevision(entityID string) Manuscript {
	versions, exists := v[entityID]

	if !exists {
		return Manuscript{
			EntityID: entityID,
		}
	}

	return versions[len(versions)-1]
}

//todo: testme
func (v VersionedManuscripts) GetVersionedManuscript(entityID string, version int) (Manuscript, error) {
	versions, exists := v[entityID]

	if !exists {
		return Manuscript{}, fmt.Errorf("manuscript %s does not exist", entityID)
	}

	if len(versions) < version {
		return Manuscript{}, fmt.Errorf("manuscript version %d of %s does not exist", version, entityID)
	}

	return versions[version-1], nil
}

func (v VersionedManuscripts) Versions(entityID string) int {
	versions, _ := v[entityID]
	return len(versions)
}
