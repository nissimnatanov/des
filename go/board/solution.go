package board

import "fmt"

// Solution is a full and valid board, with no markers for read-only. It is used as a source
// to generate other boards, interface is limited to make sure redundant methods not called.
type Solution interface {
	BoardBase
	Edit() Board
	Play() Board
}

type solutionImpl struct {
	boardBase
}

func (sol *solutionImpl) validateAndLock() error {
	var err error
	for seq := range SequenceSize {
		err = sol.validateSequence(RowSequence(seq))
		if err != nil {
			return err
		}
		err = sol.validateSequence(ColumnSequence(seq))
		if err != nil {
			return err
		}
		err = sol.validateSequence(SquareSequence(seq))
		if err != nil {
			return err
		}
	}
	sol.mode = Immutable
	return nil
}

func (sol *solutionImpl) clone(mode BoardMode) Board {
	var newBoard boardImpl
	newBoard.init(mode)
	newBoard.copyValues(&sol.boardBase)
	newBoard.recalculateAllStats()
	return &newBoard
}

func (sol *solutionImpl) Edit() Board {
	return sol.clone(Edit)
}

func (sol *solutionImpl) Play() Board {
	b := sol.clone(Play)
	b.Restart()
	return b
}

func (sol *solutionImpl) validateSequence(s Sequence) error {
	vs, dupes := sol.calcSequence(s)
	if !dupes.IsEmpty() {
		return fmt.Errorf("duplicate values")
	}
	if vs != FullValueSet() {
		return fmt.Errorf("incomplete board")
	}
	return nil
}
