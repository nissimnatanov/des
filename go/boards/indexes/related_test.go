package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"gotest.tools/v3/assert"
)

func assertRelated(t *testing.T, index, related, expected int) {
	rs := indexes.RelatedSequence(index)
	assert.Equal(t, rs.At(related), expected, "Got wrong related of %v at %v", index, related)
	assert.Equal(t, rs.Size(), indexes.RelatedSize, "Got wrong size for related sequence %v", index)
}

func TestRelatedIterator(t *testing.T) {
	assertRelated(t, 0, 0, 1)
	assertRelated(t, 0, 7, 8)
	assertRelated(t, 0, 8, 9)
	assertRelated(t, 0, 15, 72)
	assertRelated(t, 0, 19, 20)

	assertRelated(t, 40, 0, 36)
	assertRelated(t, 40, 7, 44)
	assertRelated(t, 40, 8, 4)
	assertRelated(t, 40, 15, 76)
	assertRelated(t, 40, 19, 50)

	assertRelated(t, 80, 0, 72)
	assertRelated(t, 80, 7, 79)
	assertRelated(t, 80, 8, 8)
	assertRelated(t, 80, 15, 71)
	assertRelated(t, 80, 19, 70)
}
