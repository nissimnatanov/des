package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type identifyTriplets struct {
}

func (a identifyTriplets) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	// force local alloc for cleaner profile
	var peersCache [6]int
	peers := peersCache[:0]
	var eliminationCount int

	for index, allowed := range b.AllowedSets {
		if allowed.Size() != 3 {
			// this includes non-empty cells too (allowed set is empty)
			continue
		}
		peers = a.findPeers(b, allowed, indexes.RelatedSequence(index), peers[:0])
		if len(peers) < 2 {
			continue
		}

		// found enough peers
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
		state.AddStep(Step(a.String()), StepComplexityHarder, eliminationCount)
		return StatusSucceeded
	}

	return StatusUnknown
}

func (a identifyTriplets) tryEliminate(
	board *boards.Game, index int, peers []int,
	allowed values.Set,
	seqNumFromIndex func(int) int,
	indexesFromSeq func(int) indexes.Sequence,
) (status Status, stop bool) {
	seqNum := seqNumFromIndex(index)
	var seqPeer1 = -1
	var seqPeer2 = -1

	for _, peer := range peers {
		if seqNum != seqNumFromIndex(peer) {
			continue
		}
		if seqPeer1 == -1 {
			seqPeer1 = peer
			continue
		}
		if seqPeer2 == -1 {
			seqPeer2 = peer
			continue
		}
		// we already found two peers in the same sequence, this is third which means there are
		// 4 cells with same triplet of values
		return StatusNoSolution, true
	}
	if seqPeer1 == -1 || seqPeer2 == -1 {
		// we didn't find two peers in the same sequence
		return StatusUnknown, false
	}
	return a.tryEliminateSeq(board, [3]int{index, seqPeer1, seqPeer2}, allowed, indexesFromSeq(seqNum))
}

func (a identifyTriplets) tryEliminateSeq(
	board *boards.Game, peers [3]int,
	toEliminate values.Set, seq indexes.Sequence,
) (status Status, stop bool) {
	status = StatusUnknown
	stop = false
	for index := range seq.Indexes {
		if index == peers[0] || index == peers[1] || index == peers[2] || !board.IsEmpty(index) {
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
			stop = true
		}
	}

	return status, stop
}

func (a identifyTriplets) findPeers(
	board *boards.Game, allowed values.Set, seq indexes.Sequence, peers []int,
) []int {
	for peerIndex := range seq.Indexes {
		if !board.IsEmpty(peerIndex) {
			continue
		}
		peerAllowed := board.AllowedSet(peerIndex)
		if peerAllowed == allowed {
			peers = append(peers, peerIndex)
		}
	}
	return peers
}

func (a identifyTriplets) Complexity() StepComplexity {
	return StepComplexityHard
}

func (a identifyTriplets) String() string {
	return "Identify Triplets"
}
