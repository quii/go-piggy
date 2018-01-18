package manuscript

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
