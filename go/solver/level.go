package solver

import "fmt"

type Level int8

const (
	LevelUnknown Level = iota
	LevelEasy
	LevelMedium
	LevelHard
	LevelVeryHard
	LevelEvil
	LevelDarkEvil
	LevelNightmare
	LevelBlackHole
)

func (l Level) String() string {
	switch l {
	case LevelUnknown:
		return "Unknown"
	case LevelEasy:
		return "Easy"
	case LevelMedium:
		return "Medium"
	case LevelHard:
		return "Hard"
	case LevelVeryHard:
		return "VeryHard"
	case LevelEvil:
		return "Evil"
	case LevelDarkEvil:
		return "DarkEvil"
	case LevelNightmare:
		return "Nightmare"
	case LevelBlackHole:
		return "BlackHole"
	default:
		return fmt.Sprintf("WRONG SudokuLevel: %d", l)
	}
}

type LevelComplexityBar int64

const (
	// TODO: need to adjust level bars
	LevelEasyBar      LevelComplexityBar = 100
	LevelMediumBar    LevelComplexityBar = 250
	LevelHardBar      LevelComplexityBar = 500
	LevelVeryHardBar  LevelComplexityBar = 1000
	LevelEvilBar      LevelComplexityBar = 2000
	LevelDarkEvilBar  LevelComplexityBar = 5000
	LevelNightmareBar LevelComplexityBar = 15000
)

// LevelFromComplexity returns the level of the Sudoku puzzle based on its complexity
func LevelFromComplexity(complexity int64) Level {
	if complexity <= int64(LevelEasyBar) {
		return LevelEasy
	} else if complexity <= int64(LevelMediumBar) {
		return LevelMedium
	} else if complexity <= int64(LevelHardBar) {
		return LevelHard
	} else if complexity <= int64(LevelVeryHardBar) {
		return LevelVeryHard
	} else if complexity <= int64(LevelEvilBar) {
		return LevelEvil
	} else if complexity <= int64(LevelDarkEvilBar) {
		return LevelDarkEvil
	} else if complexity <= int64(LevelNightmareBar) {
		return LevelNightmare
	} else {
		return LevelBlackHole
	}
}
