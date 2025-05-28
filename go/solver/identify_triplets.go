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
	freeCellOnStart := b.FreeCellCount()

	for index, allowed := range b.AllAllowedValues {
		if allowed.Size() != 3 {
			// this includes non-empty cells too (allowed set is empty)
			continue
		}
		peers = a.findPeers(b, index, allowed, indexes.RelatedSequence(index), peers[:0])
		if len(peers) < 2 {
			continue
		}

		// found enough peers
		status := a.tryEliminate(b, index, peers, allowed, indexes.RowFromIndex, indexes.RowSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
		} else if status != StatusUnknown {
			return status
		}

		status = a.tryEliminate(b, index, peers, allowed, indexes.ColumnFromIndex, indexes.ColumnSequence)
		if status == StatusSucceeded {
			eliminationCount++
			// do not stop yet if we found a solution, let's check other peers
		} else if status != StatusUnknown {
			return status
		}

		status = a.tryEliminate(b, index, peers, allowed, indexes.SquareFromIndex, indexes.SquareSequence)
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
		state.AddStep(Step(a.String()), a.Complexity(), eliminationCount)
		return StatusSucceeded
	}

	return StatusUnknown
}

func (a identifyTriplets) tryEliminate(
	board *boards.Game, index int, peers []int,
	allowed values.Set,
	seqNumFromIndex func(int) int,
	indexesFromSeq func(int) indexes.Sequence,
) Status {
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
		return StatusNoSolution
	}
	if seqPeer1 == -1 || seqPeer2 == -1 {
		// we didn't find two peers in the same sequence
		return StatusUnknown
	}
	return a.tryEliminateSeq(board, [3]int{index, seqPeer1, seqPeer2}, allowed, indexesFromSeq(seqNum))
}

func (a identifyTriplets) tryEliminateSeq(
	board *boards.Game, peers [3]int,
	toEliminate values.Set, seq indexes.Sequence,
) Status {
	status := StatusUnknown
	for _, index := range seq {
		if index == peers[0] || index == peers[1] || index == peers[2] || !board.IsEmpty(index) {
			continue
		}

		tempAllowed := board.AllowedValues(index)
		if values.Intersect(tempAllowed, toEliminate).Size() == 0 {
			// no intersection, continue
			continue
		}

		// found a cell that we can remove values - turn them off
		board.DisallowValues(index, toEliminate)

		// since we are here, we can now easily check if we have only one allowed value left
		tempAllowed = tempAllowed.Without(toEliminate)
		switch tempAllowed.Size() {
		case 0:
			// no allowed values left, this is a dead end
			return StatusNoSolution
		case 1:
			// only one allowed value left, let's set it
			for _, v := range tempAllowed.Values() {
				// we can safely assume that this is the only value left
				board.Set(index, v)
			}
		}
		status = StatusSucceeded
	}

	return status
}

func (a identifyTriplets) findPeers(
	board *boards.Game, index int, allowed values.Set, seq indexes.Sequence, peers []int,
) []int {
	for _, peerIndex := range seq {
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
	}
	return peers
}

func (a identifyTriplets) Complexity() StepComplexity {
	return StepComplexityHarder
}

func (a identifyTriplets) String() string {
	return "Identify Triplets"
}
