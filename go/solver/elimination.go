package solver

import (
	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

const crossSequenceConstraintComplexity = StepComplexityHarder

func eliminateInSequence(
	b *boards.Game, vs values.Set, eliminateInSeq indexes.Sequence,
) (Status, int) {
	var maxDisallowed int
	for _, index := range eliminateInSeq {
		if !b.IsEmpty(index) {
			// skip the main sequence we are working on, or if non-empty
			continue
		}
		allowed := b.AllowedValues(index)
		disallow := values.Intersect(allowed, vs)
		if disallow.IsEmpty() {
			continue
		}
		if maxDisallowed == 0 {
			maxDisallowed = disallow.Size()
		} else {
			// count the stats based on the num of disallowed per segment
			maxDisallowed = max(maxDisallowed, disallow.Size())
		}
		allowedAfter := allowed.Without(disallow)
		switch allowedAfter.Size() {
		case 0:
			// about to eliminate the last allowed value, no solution
			return StatusNoSolution, 0
		case 1:
			// if, after elimination, there is just one allowed value,
			// we can just set the other value instead of elimination
			b.Set(index, allowedAfter.First())

		default:
			b.DisallowValues(index, disallow)
		}
	}
	if maxDisallowed > 0 {
		return StatusSucceeded, maxDisallowed
	}
	return StatusUnknown, 0
}

func runConstraints(
	step Step,
	state AlgorithmState,
	// index of the sequence being scanned
	seqIndex int,
	// the sequence of board indexes in this row/column
	seq indexes.Sequence,
	crossSeqFromIndex func(index int) int,
	crossSeqNotShared func(crossSeqIndex int, seqIndex int) indexes.Sequence,
) Status {
	b := state.Board()

	// Each row/col is broken into 3 segments, each belong to a square, each crosses exactly 3
	// squares. The two arrays below map the square index to the allowed values in that square
	// segment shared with the row/col.
	//
	// Similarly, each square can be broken into 3 segments per row, and 3 segments per column.
	//
	// cs => cross sequence
	csIndexes := [3]int{-1, -1, -1} // -1 means the entry has not been used
	var csAllowed [3]values.Set

	for index, allowed := range b.AllowedValuesIn(seq) {
		csIndex := crossSeqFromIndex(index)
		for i, csi := range csIndexes {
			if csi == -1 {
				csIndexes[i] = csIndex
				csAllowed[i] = allowed
				break
			} else if csi == csIndex {
				csAllowed[i] = csAllowed[i].With(allowed)
				break
			} else if i == 2 {
				// we are on the last one and got > 3 distinct squares in the same row/col
				// (or vice versa)
				panic(">3 shared segments not possible")
			}
		}
	}

	// now, let's see if any value is present in a single square and if yes, it can be
	// eliminated from all other squares which are part of the same row/column
	var csUnique [3]values.Set
	csUnique[0] = csAllowed[0].Without(csAllowed[1].With(csAllowed[2]))
	csUnique[1] = csAllowed[1].Without(csAllowed[0].With(csAllowed[2]))
	csUnique[2] = csAllowed[2].Without(csAllowed[0].With(csAllowed[1]))

	eliminateCount := 0
	for i, csi := range csIndexes {
		if csi == -1 {
			// it is ok for segments to stay unused in case they were filled with values
			break
		}
		if csUnique[i].IsEmpty() {
			// no unique values in this segment, move to the next one
			continue
		}

		// we can safely eliminate these values
		status, count := eliminateInSequence(b, csUnique[i], crossSeqNotShared(csi, seqIndex))
		switch status {
		case StatusUnknown:
			// keep looking
		case StatusSucceeded:
			eliminateCount += count
		default:
			// if we eliminated the only candidate in the other square, bail out
			return status
		}
	}

	if eliminateCount > 0 {
		state.AddStep(step, crossSequenceConstraintComplexity, eliminateCount)
		return StatusSucceeded
	}
	return StatusUnknown
}
