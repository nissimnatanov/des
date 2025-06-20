package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type identifyPairs struct {
}

func (a identifyPairs) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	// force local alloc for cleaner profile
	var peersCache [5]int
	var eliminationCount int
	freeCellOnStart := b.FreeCellCount()

	for index, allowed := range b.AllAllowedValues {
		if allowed.Size() != 2 {
			// this includes non-empty cells too (allowed set is empty)
			continue
		}
		peers := a.findPeers(b, index, allowed, indexes.RelatedSequence(index), peersCache[:])
		if len(peers) == 0 {
			// no peer found
			continue
		}

		status := a.tryEliminate(state, index, peers, allowed, indexes.RowFromIndex, indexes.RowSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
		} else if status != StatusUnknown {
			return status
		}

		status = a.tryEliminate(state, index, peers, allowed, indexes.ColumnFromIndex, indexes.ColumnSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
		} else if status != StatusUnknown {
			return status
		}

		status = a.tryEliminate(state, index, peers, allowed, indexes.SquareFromIndex, indexes.SquareSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
		} else if status != StatusUnknown {
			return status
		}
		if b.FreeCellCount() < freeCellOnStart {
			// if we set at least one value, we can stop now and try faster algos
			return StatusSucceeded
		}
	}
	if eliminationCount > 0 {
		// if we found at least one index that lead to elimination, let's stop
		// and go back to cheaper algorithm such as theOnlyChoice
		return StatusSucceeded
	}
	return StatusUnknown
}

func (a identifyPairs) tryEliminate(
	state AlgorithmState, index int, peers []int,
	allowed values.Set,
	seqNumFromIndex func(int) int,
	indexesFromSeq func(int) indexes.Sequence,
) Status {
	seqNum := seqNumFromIndex(index)
	var seqPeer = -1

	for _, peer := range peers {
		if seqNum != seqNumFromIndex(peer) {
			continue
		}
		if seqPeer != -1 {
			// we already found one peer in the same sequence, this is second
			// in the same sequence which means there are
			// 3 cells with same pair of values
			state.AddStep(Step(a.String()), a.Complexity(), 1)
			return StatusNoSolution
		}
		seqPeer = peer
	}
	if seqPeer == -1 {
		// we didn't find at least one peer in the same sequence
		return StatusUnknown
	}

	return a.tryEliminateSeq(state, index, seqPeer, allowed, indexesFromSeq(seqNum))
}

func (a identifyPairs) tryEliminateSeq(
	state AlgorithmState, p1, p2 int,
	toEliminate values.Set, seq indexes.Sequence,
) Status {
	status := StatusUnknown
	board := state.Board()
	for _, index := range seq {
		if index == p1 || index == p2 || !board.IsEmpty(index) {
			continue
		}

		tempAllowed := board.AllowedValues(index)
		if values.Intersect(tempAllowed, toEliminate).Size() == 0 {
			// no intersection, continue
			continue
		}

		// found a cell that we can remove values - turn them off
		state.AddStep(Step(a.String()), a.Complexity(), 1)
		status = StatusSucceeded

		// since we are here, we can now easily check if we have only one allowed value left
		tempAllowed = tempAllowed.Without(toEliminate)
		switch tempAllowed.Size() {
		case 0:
			// no allowed values left, this is a dead end
			return StatusNoSolution
		case 1:
			// only one allowed value left, let's set it
			board.Set(index, tempAllowed.First())
		default:
			board.DisallowValues(index, toEliminate)
		}
	}

	return status
}

func (a identifyPairs) findPeers(
	board *boards.Game,
	index int,
	allowed values.Set,
	related indexes.Sequence,
	peersCache []int,
) []int {
	peers := peersCache[:0]
	for _, peerIndex := range related {
		if peerIndex < index {
			// we only need to search forward since we already searched for peers
			// of the previous indexes
			continue
		}
		if !board.IsEmpty(peerIndex) {
			continue
		}
		peerAllowed := board.AllowedValues(peerIndex)
		if peerAllowed == allowed {
			peers = append(peers, peerIndex)
		}
		// keep going to find more peers, this helps us to invalidate boards if > one peer is found
	}
	return peers
}

func (a identifyPairs) Complexity() StepComplexity {
	return StepComplexityHard
}

func (a identifyPairs) String() string {
	return "Identify Pairs"
}
