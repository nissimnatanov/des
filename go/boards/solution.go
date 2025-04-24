package boards

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

// NewSolution creates a new solution from a valid and solved board.
func NewSolution(b *Play) *Solution {
	var sol Solution
	sol.copyValues(&b.base)
	err := sol.validateAndLock()
	if err != nil {
		panic(err)
	}
	return &sol
}

// Solution is a full and valid board. It may or may not include read-only markers, it can be used
// as a source to generate other boards.
//
// Its mode is Immutable to indicate that it cannot be modified.
type Solution struct {
	base
}

func (sol *Solution) validateAndLock() error {
	var err error
	for seq := range SequenceSize {
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
	return nil
}
func (b *Solution) IsValidCell(index int) bool {
	// solutions are pre-validated as long as they are not nil
	if b == nil {
		panic("solution is nil")
	}
	return true
}

func (b *Solution) String() string {
	var sb strings.Builder
	WriteValues(b, bufio.NewWriter(&sb))
	return sb.String()
}

func (sol *Solution) Clone(mode Mode) *Play {
	newBoard := &Play{}
	newBoard.init(mode)
	newBoard.copyValues(&sol.base)
	newBoard.recalculateAllStats()
	return newBoard
}

func (sol *Solution) validateSequence(s indexes.Sequence) error {
	vs, dupes := sol.calcSequence(s)
	if !dupes.IsEmpty() {
		return fmt.Errorf("duplicate values")
	}
	if vs != values.FullSet() {
		return fmt.Errorf("incomplete board")
	}
	return nil
}
