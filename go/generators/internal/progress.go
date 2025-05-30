package internal

import "fmt"

type Progress int

const (
	ProgressUnknown Progress = iota
	TooEarly
	BelowLevel
	AtLevelKeepGoing
	AtLevelStop
	AboveLevel
)

func (p Progress) String() string {
	switch p {
	case ProgressUnknown:
		return "Unknown"
	case TooEarly:
		return "TooEarly"
	case BelowLevel:
		return "BelowLevel"
	case AtLevelKeepGoing:
		return "AtLevelKeepGoing"
	case AtLevelStop:
		return "AtLevelStop"
	case AboveLevel:
		return "AboveLevel"
	default:
		return fmt.Sprintf("UnknownProgress(%d)", p)
	}
}
