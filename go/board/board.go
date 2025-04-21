package board

const (
	BoardSize         = 81
	SequenceSize      = 9
	MinValidBoardSize = 17
)

type BoardMode int

const (
	Immutable BoardMode = iota
	Edit
	Play
)

// BoardBase is a board interface that only allows reads
type BoardBase interface {
	Get(index int) Value
	IsEmpty(index int) bool
	IsReadOnly(index int) bool

	String() string
}

type boardBaseInternal interface {
	BoardBase
	setInternal(index int, v Value, readOnly bool) Value
}

type Board interface {
	BoardBase

	Mode() BoardMode

	RowSet(row int) ValueSet
	ColumnSet(col int) ValueSet
	SquareSet(square int) ValueSet

	AllowedSet(index int) ValueSet
	Count(v Value) int
	FreeCellCount() int

	IsValidCell(index int) bool
	IsValid() bool
	IsSolved() bool

	// available in Play and Edit modes only
	Set(index int, v Value)
	// available in Edit modes only
	SetReadOnly(index int, v Value)

	// available in Play and Edit modes only
	Disallow(index int, v Value)
	DisallowSet(index int, vs ValueSet)
	DisallowReset(index int)

	// available in Play and Edit modes only
	Restart()

	// available in any mode
	Clone(mode BoardMode) Board
}
