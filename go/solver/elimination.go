package solver

import (
	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

const crossSequenceConstraintComplexity = StepComplexityHarder

func eliminateInSequence(
	b *boards.Game, v values.Value, eliminateInSeq indexes.Sequence,
) Status {
	succeeded := false
	for _, index := range eliminateInSeq {
		if !b.IsEmpty(index) {
			// skip the main sequence we are working on, or if non-empty
			continue
		}
		allowed := b.AllowedValues(index)
		if !allowed.Contains(v) {
			continue
		}
		switch allowed.Size() {
		case 0:
			panic("eliminateInSequence: value was allowed but not in the set")
		case 1:
			// about to eliminate the last allowed value, no solution
			return StatusNoSolution
		case 2:
			// there are two values left, we are about to eliminate one,
			// we can just set the other value instead of elimination
			b.Set(index, allowed.Without(v.AsSet()).First())
		default:
			b.DisallowValue(index, v)
		}
		succeeded = true
	}
	if succeeded {
		return StatusSucceeded
	}
	return StatusUnknown
}
