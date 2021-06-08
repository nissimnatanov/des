package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertRelated(t *testing.T, index, related, expected int) {
	ri := RelatedSequence(index).Iterator()
	actual := -1
	for i := 0; i <= related; i++ {
		assert.Truef(t, ri.Next(), "Expected Next on related index of %v at %v", index, i)
		actual = ri.Value()
	}
	assert.Equal(t, expected, actual, "Got wrong related of %v at %v", index, related)
	if related == RelatedSize-1 {
		assert.Falsef(t, ri.Next(), "Expected Next to return false on related end index")
	}
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
