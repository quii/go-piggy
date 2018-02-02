package manuscript

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionedManuscript_CurrentRevision(t *testing.T) {
	m1 := Manuscript{Title: "one"}
	m2 := Manuscript{Title: "two"}
	m3 := Manuscript{Title: "three"}

	vm := VersionedManuscript{m1, m2, m3}

	assert.Equal(t, m3, vm.CurrentRevision())
}

func TestVersionedManuscript_Version(t *testing.T) {
	m1 := Manuscript{Title: "one"}
	m2 := Manuscript{Title: "two"}
	m3 := Manuscript{Title: "three"}

	vm := VersionedManuscript{m1, m2, m3}

	version1, _ := vm.Version(1)
	assert.Equal(t, m1, version1)

	_, notExist := vm.Version(5)
	assert.Error(t, notExist)
}
