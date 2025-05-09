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
	peers := peersCache[:0]
	var eliminationCount int

	for index, allowed := range b.AllowedSets {
		if allowed.Size() != 2 {
			// this includes non-empty cells too (allowed set is empty)
			continue
		}
		peers = a.findPeers(b, allowed, indexes.RelatedSequence(index), peers[:0])
		switch {
		case len(peers) < 1:
			// no peer found
			continue
		case len(peers) > 4:
			// if cell has more than 3 peers with same pair of allowed values, than at least two of them
			// share the same row/col/square with the current cell meaning within the same
			// scope there are more than two cells with same pair of values
			return StatusNoSolution
		}

		status, stop := a.tryEliminate(b, index, peers, allowed, indexes.RowFromIndex, indexes.RowSequence)
		stopOnSuccess := false
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
			stopOnSuccess = stopOnSuccess || stop
		} else if stop {
			return status
		}

		status, stop = a.tryEliminate(b, index, peers, allowed, indexes.ColumnFromIndex, indexes.ColumnSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
			stopOnSuccess = stopOnSuccess || stop
		} else if stop {
			return status
		}

		status, stop = a.tryEliminate(b, index, peers, allowed, indexes.SquareFromIndex, indexes.SquareSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
			stopOnSuccess = stopOnSuccess || stop
		} else if stop {
			return status
		}

		if stopOnSuccess {
			break
		}
	}
	if eliminationCount > 0 {
		// if we found at least one index that lead to elimination, let's stop
		// and go back to cheaper algorithm such as theOnlyChoice
		state.AddStep(Step(a.String()), StepComplexityHarder, eliminationCount)
		return StatusSucceeded
	}
	return StatusUnknown
}

func (a identifyPairs) tryEliminate(
	board *boards.Game, index int, peers []int,
	allowed values.Set,
	seqNumFromIndex func(int) int,
	indexesFromSeq func(int) indexes.Sequence,
) (status Status, stop bool) {
	seqNum := seqNumFromIndex(index)
	var seqPeer = -1

	for _, peer := range peers {
		if seqNum != seqNumFromIndex(peer) {
			continue
		}
		if seqPeer == -1 {
			seqPeer = peer
			continue
		}
		// we already found one peer in the same sequence, this is second which means there are
		// 3 cells with same pair of values
		return StatusNoSolution, true
	}
	if seqPeer == -1 {
		// we didn't find at least one peer in the same sequence
		return StatusUnknown, false
	}

	return a.tryEliminateSeq(board, [2]int{index, seqPeer}, allowed, indexesFromSeq(seqNum))
}

func (a identifyPairs) tryEliminateSeq(
	board *boards.Game, peers [2]int,
	toEliminate values.Set, seq indexes.Sequence,
) (status Status, stop bool) {
	status = StatusUnknown
	stop = false
	for index := range seq.Indexes {
		if index == peers[0] || index == peers[1] || !board.IsEmpty(index) {
			continue
		}

		tempAllowed := board.AllowedSet(index)
		if values.Intersect(tempAllowed, toEliminate).Size() == 0 {
			// no intersection, continue
			continue
		}

		// found a cell that we can remove values - turn them off
		board.DisallowSet(index, toEliminate)
		status = StatusSucceeded

		// since we are here, we can now easily check if we have only one allowed value left
		tempAllowed = tempAllowed.Without(toEliminate)
		switch tempAllowed.Size() {
		case 0:
			// no allowed values left, this is a dead end
			return StatusNoSolution, true
		case 1:
			// only one allowed value left, let's set it
			for v := range tempAllowed.Values {
				// we can safely assume that this is the only value left
				board.Set(index, v)
			}
			// once we set a value, no need to continue this algorithm since we might
			// get a lot cheaper ones now
			stop = true
		}
	}

	return status, stop
}

func (a identifyPairs) findPeers(board *boards.Game, allowed values.Set, seq indexes.Sequence, peers []int) []int {
	for i := range seq.Indexes {
		if !board.IsEmpty(i) {
			continue
		}
		peerAllowed := board.AllowedSet(i)
		if peerAllowed == allowed {
			peers = append(peers, i)
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
