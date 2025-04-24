package board

import (
	"github.com/nissimnatanov/des/go/board/indexes"
	"github.com/nissimnatanov/des/go/board/values"
)

const (
	SequenceSize = indexes.BoardSequenceSize
	Size         = indexes.BoardSize

	MinValidBoardSize = 17
)

type Mode int

const (
	// Immutable mode means input board shall not be modified.
	Immutable Mode = iota
	Edit
	Play
)

type Board interface {
	Mode() Mode

	Get(index int) values.Value
	IsEmpty(index int) bool
	IsReadOnly(index int) bool

	RowSet(row int) values.Set
	ColumnSet(col int) values.Set
	SquareSet(square int) values.Set

	AllowedSet(index int) values.Set
	FreeCellCount() int

	IsValidCell(index int) bool
	IsValid() bool
	IsSolved() bool

	Clone(mode Mode) Board

	// CloneInto copies the current board into the dst board, both must be boards
	// created with New (not NewSolution) or cloned from such.
	CloneInto(mode Mode, dst Board)

	String() string

	// available in Play and Edit modes only
	Set(index int, v values.Value)

	// available in Play and Edit modes only
	Disallow(index int, v values.Value)
	DisallowSet(index int, vs values.Set)
	DisallowReset(index int)

	// available in Play and Edit modes only
	Restart()

	// available in Edit mode only
	SetReadOnly(index int, v values.Value)
}
