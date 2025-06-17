package boards

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

// NewSolution creates a new solution from a valid and solved board.
func NewSolution(b *Game) *Solution {
	var sol Solution
	sol.init(Edit)
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
	sol.base.mode = Immutable
	return nil
}

func (b *Solution) isValidCell(index int) bool {
	return true
}

func (b *Solution) getDisallowedByUser(index int) values.Set {
	return values.EmptySet
}

func (b *Solution) String() string {
	var sb strings.Builder
	writeValues(b, &sb)
	return sb.String()
}

func (b *Solution) CloneInto(mode Mode, board *Game) {
	board.init(mode)
	board.copyValues(&b.base)
	board.recalculateAllStats()
}

func (sol *Solution) Clone(mode Mode) *Game {
	newBoard := &Game{}
	sol.CloneInto(mode, newBoard)
	return newBoard
}

func (sol *Solution) validateSequence(s indexes.Sequence) error {
	vs, dupes := sol.calcSequence(s)
	if !dupes.IsEmpty() {
		return fmt.Errorf("duplicate values")
	}
	if vs != values.FullSet {
		return fmt.Errorf("incomplete board")
	}
	return nil
}

func (sol *Solution) MarshalJSON() ([]byte, error) {
	return json.Marshal(Serialize(sol))
}
