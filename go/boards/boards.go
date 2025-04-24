package boards

import (
	"github.com/nissimnatanov/des/go/boards/indexes"
)

const (
	SequenceSize = indexes.BoardSequenceSize
	Size         = indexes.BoardSize

	MinValidBoardSize = 17
)

type Mode int

const (
	// ImmutableMode mode means input board shall not be modified.
	ImmutableMode Mode = iota
	PlayMode
	EditMode
)
