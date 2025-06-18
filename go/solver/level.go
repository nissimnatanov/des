package solver

import (
	"fmt"
	"strings"
)

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

func LevelFromString(s string) Level {
	for l := LevelEasy; l <= LevelBlackHole; l++ {
		if strings.EqualFold(l.String(), s) {
			return l
		}
	}
	return LevelUnknown
}

type LevelComplexityBar int64

const (
	LevelEasyBar     LevelComplexityBar = 125
	LevelMediumBar   LevelComplexityBar = 350
	LevelHardBar     LevelComplexityBar = 1000
	LevelVeryHardBar LevelComplexityBar = 3000

	// many recursive steps are needed
	LevelEvilBar LevelComplexityBar = 10000
	// aligned with the second recursion step complexity
	LevelDarkEvilBar LevelComplexityBar = 40000
	// seen only few puzzles with complexity above this level
	LevelNightmareBar LevelComplexityBar = 100000
)

// LevelFromComplexity returns the level of the Sudoku puzzle based on its complexity
func LevelFromComplexity(complexity StepComplexity) Level {
	cb := LevelComplexityBar(complexity)
	if cb <= LevelEasyBar {
		return LevelEasy
	} else if cb <= LevelMediumBar {
		return LevelMedium
	} else if cb <= LevelHardBar {
		return LevelHard
	} else if cb <= LevelVeryHardBar {
		return LevelVeryHard
	} else if cb <= LevelEvilBar {
		return LevelEvil
	} else if cb <= LevelDarkEvilBar {
		return LevelDarkEvil
	} else if cb <= LevelNightmareBar {
		return LevelNightmare
	} else {
		return LevelBlackHole
	}
}
