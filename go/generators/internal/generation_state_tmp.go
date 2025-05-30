package internal

/*
		indexes := newIndexManager()
		// Mark all cells as read-only and update index manager
		for index := range boards.Size {
			v := board.Get(index)
			if v == 0 {
				indexes.RemoveIndex(index)
			} else {
				board.SetReadOnly(index, v)
			}
		}

type GenerationState struct {
	indexManager *indexManager
}

func (gs *GenerationState) RemainedIndexes() []int {
	return gs.indexManager.Remained()
}

func (gs *GenerationState) RemainedIndexesSize() int {
	return gs.indexManager.RemainedSize()
}

func (gs *GenerationState) RemoveIndex(index int) {
	gs.indexManager.RemoveIndex(index)
}

func (bs *BoardState) tryProve(ctx context.Context) bool {
	if bs.proven {
		panic("Do not prove twice on the same State (even after clone)!")
	}

	res := bs.state.prover.Run(ctx, bs.board)
	bs.proven = res.Status == solver.StatusSucceeded
	return bs.proven
}

func (bs *BoardState) TryRemoveFromIndex(ctx context.Context, index int) bool {
	bs.indexManager.PrioritizeIndex(index)
	progress := bs.tryRemoveBatch(ctx, 1)
	return progress.KeepGoing() || progress == AtLevelStop
}


func (bs *BoardState) tryRemoveOne(ctx context.Context) Progress {
	defer bs.checkIntegrity()

	for bs.indexManager.RemainedSize() > 0 {
		progress := bs.tryRemoveBatch(ctx, 1)
		if progress != Failed {
			return progress
		}
		// if batch has only one index and it fails, it is removed from remained list.
	}

	return Failed
}

func (bs *BoardState) Clone() *BoardState {
	// shallow clone first
	newState := *bs
	newState.board = newState.board.Clone(boards.Edit)
	newState.indexManager = newState.indexManager.clone()
	newState.checkIntegrity()
	return &newState
}

func (bs *BoardState) RestoreSimpleValues(ctx context.Context, minComplexity solver.StepComplexity) {
	if !bs.proven {
		panic("do not call RestoreSimpleValues without proving")
	}

	bs.indexManager.shuffleRemoved(bs.state.rand)

	lastResult := bs.solveInternal(ctx)

	for _, index := range bs.indexManager.Removed() {
		if !bs.board.IsEmpty(index) {
			// ignore cells with value, their indexes were marked as 'removed' because removing them leads to
			// unsolvable board.
			continue
		}
		bs.board.SetReadOnly(index, bs.state.solution.Get(index))
		result := bs.solveInternal(ctx)
		// TODO: review the 5 here
		if result.Steps.Complexity < minComplexity ||
			(lastResult.Steps.Complexity-result.Steps.Complexity) > 5 {
			// Not a good choice for restoration - remove it
			bs.board.Set(index, 0)
		} else {
			bs.indexManager.RestoreRemoved(index)
			lastResult = result
		}
	}
	bs.checkIntegrity()
	bs.cachedSolve = lastResult
	bs.cachedProgress = bs.shouldContinue(ctx)
}


func (bs *BoardState) mergeWith(other *BoardState) {
	bs.proven = false
	for index := range boards.Size {
		currentValue := bs.board.Get(index)
		if currentValue == 0 {
			continue
		}

		otherValue := other.board.Get(index)
		if otherValue != 0 && otherValue != currentValue {
			panic("Values of current board do not match the values of merged one.")
		}

		if otherValue == 0 {
			if bs.indexManager.TryRemoveIndex(index) {
				bs.board.Set(index, 0)
			}
		}
	}
}

func TryMergeBoardStates(ctx context.Context, state1, state2 *BoardState) *BoardState {
	merge := state1.Clone()
	merge.mergeWith(state2)
	if merge.tryProve(ctx) {
		// if we can prove the merged board, return it
		return merge
	}
	// do not return invalid boards
	return nil
}
*/
