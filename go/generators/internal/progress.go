package internal

import "fmt"

type Progress int

const (
	ProgressUnknown Progress = iota
	TooEarly
	BelowMinLevel
	InRangeKeepGoing
	InRangeStop
	AboveMaxLevel
)

func (p Progress) String() string {
	switch p {
	case ProgressUnknown:
		return "Unknown"
	case TooEarly:
		return "TooEarly"
	case BelowMinLevel:
		return "BelowMinLevel"
	case InRangeKeepGoing:
		return "InRangeKeepGoing"
	case InRangeStop:
		return "InRangeStop"
	case AboveMaxLevel:
		return "AboveMaxLevel"
	default:
		return fmt.Sprintf("UnknownProgress(%d)", p)
	}
}
