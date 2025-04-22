package board

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/board/indexes"
	"github.com/nissimnatanov/des/go/board/values"
)

// Solution is a full and valid board. It may or may not include read-only markers, it can be used
// as a source to generate other boards.
//
// Its mode is Immutable to indicate that it cannot be modified.
type Solution interface {
	Board

	// isSolution prevents casting regular board to solution without validating it first
	isSolution()
}

// NewSolution creates a new solution from a valid and solved board.
func NewSolution(b Board) Solution {
	if sol, ok := b.(Solution); ok {
		return sol
	}
	var sol solutionImpl
	sol.init(Immutable)

	if bb, ok := b.(*boardImpl); ok {
		sol.copyValues(&bb.base)
	} else {
		for i := range Size {
			sol.values[i] = b.Get(i)
			sol.readOnlyFlags.Set(i, b.IsReadOnly(i))
		}
	}

	err := sol.validateAndLock()
	if err != nil {
		panic(err)
	}
	return &sol
}

type solutionImpl struct {
	base
}

func (sol *solutionImpl) validateAndLock() error {
	var err error
	for seq := range indexes.SequenceSize {
		err = sol.validateSequence(indexes.RowSequence(seq))
		if err != nil {
			return err
		}
		err = sol.validateSequence(indexes.ColumnSequence(seq))
		if err != nil {
			return err
		}
		err = sol.validateSequence(indexes.SquareSequence(seq))
		if err != nil {
			return err
		}
	}
	sol.mode = Immutable
	return nil
}

func (sol *solutionImpl) isSolution() {}

func (sol *solutionImpl) RowSet(row int) values.Set {
	if row < 0 || row >= indexes.SequenceSize {
		panic(fmt.Sprintf("Row %d is out of bounds", row))
	}
	return values.FullSet()
}

func (sol *solutionImpl) ColumnSet(col int) values.Set {
	if col < 0 || col >= indexes.SequenceSize {
		panic(fmt.Sprintf("Column %d is out of bounds", col))
	}
	return values.FullSet()
}

func (sol *solutionImpl) SquareSet(sq int) values.Set {
	if sq < 0 || sq >= indexes.SequenceSize {
		panic(fmt.Sprintf("Square %d is out of bounds", sq))
	}
	return values.FullSet()
}

func (sol *solutionImpl) AllowedSet(index int) values.Set {
	if index < 0 || index >= Size {
		panic(fmt.Sprintf("Index %d is out of bounds", index))
	}
	return values.EmptySet()
}

func (sol *solutionImpl) Count(v values.Value) int {
	v.Validate()
	if v == 0 {
		return 0
	}
	return indexes.SequenceSize
}

func (sol *solutionImpl) FreeCellCount() int {
	return 0
}

func (sol *solutionImpl) IsValidCell(index int) bool {
	if index < 0 || index >= Size {
		panic(fmt.Sprintf("Index %d is out of bounds", index))
	}
	return true
}

func (sol *solutionImpl) IsValid() bool {
	return true
}

func (sol *solutionImpl) IsSolved() bool {
	return true
}

func (b *solutionImpl) String() string {
	var sb strings.Builder
	WriteValues(b, bufio.NewWriter(&sb))
	return sb.String()
}

func (sol *solutionImpl) Clone(mode Mode) Board {
	if mode == Immutable {
		return sol
	}

	var newBoard boardImpl
	newBoard.init(mode)
	newBoard.copyValues(&sol.base)
	newBoard.recalculateAllStats()
	return &newBoard
}

func (sol *solutionImpl) CloneInto(mode Mode, dst Board) {
	panic("CloneInto is not supported for Solution, use Clone")
}

func (sol *solutionImpl) validateSequence(s indexes.Sequence) error {
	vs, dupes := sol.calcSequence(s)
	if !dupes.IsEmpty() {
		return fmt.Errorf("duplicate values")
	}
	if vs != values.FullSet() {
		return fmt.Errorf("incomplete board")
	}
	return nil
}

func (sol *solutionImpl) Set(index int, v values.Value) {
	panic("Cannot set value on Solution board")
}

func (sol *solutionImpl) Disallow(index int, v values.Value) {
	panic("Cannot disallow value on Solution board")
}
func (sol *solutionImpl) DisallowSet(index int, vs values.Set) {
	panic("Cannot disallow value set on Solution board")
}
func (sol *solutionImpl) DisallowReset(index int) {
	panic("Cannot reset disallowed values on Solution board")
}

func (sol *solutionImpl) Restart() {
	panic("Cannot restart Solution board")
}

func (sol *solutionImpl) SetReadOnly(index int, v values.Value) {
	panic("Cannot set read-only value on Solution board")
}
