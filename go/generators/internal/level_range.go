package internal

import (
	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/solver"
)

type LevelRange struct {
	Min solver.Level
	Max solver.Level
}

func (lr LevelRange) WithDefaults() LevelRange {
	if lr.Min == solver.LevelUnknown {
		// If Min is unknown, then we want to generate all levels.
		lr.Min = solver.LevelEasy
	}
	if lr.Max == solver.LevelUnknown {
		// If Max is unknown, generate any level.
		lr.Max = solver.LevelBlackHole
	}
	switch {
	case lr.Min < solver.LevelUnknown || lr.Min > solver.LevelBlackHole:
		panic("invalid Min Level: " + lr.Min.String())
	case lr.Max < solver.LevelUnknown || lr.Max > solver.LevelBlackHole:
		panic("invalid MaxLevel: " + lr.Max.String())
	case lr.Min > lr.Max:
		// if both are set, then Min must be less than or equal to Max.
		panic("MinLevel cannot be greater than MaxLevel: " + lr.Min.String() + " > " + lr.Max.String())
	}
	return lr
}

func (lr LevelRange) shouldContinue(r *random.Random, board *boards.Game, res *solver.Result) Progress {
	if board.FreeCellCount() < 32 {
		// too early even for easy games.
		return TooEarly
	}

	if res.Level < lr.Min {
		return BelowMinLevel
	}

	if res.Level > lr.Max {
		// Overflow, stop.
		return AboveMaxLevel
	}

	// shoot for the max desired level
	if shouldContinueAtLevel(lr.Max, r) {
		// Keep going, we are at the desired level.
		return InRangeKeepGoing
	}

	// We are at the desired level, but do not want to continue.
	return InRangeStop
}

func shouldContinueAtLevel(desiredLevel solver.Level, r *random.Random) bool {
	switch desiredLevel {
	case solver.LevelEasy:
		// For easy games - keep trying (otherwise, game can be too easy).
		return r.PercentProbability(95)
	case solver.LevelMedium:
		// For medium games - keep trying a bit less.
		return r.PercentProbability(85)
	case solver.LevelHard, solver.LevelVeryHard:
		// For hard games - continue in half of the cases..
		return r.PercentProbability(75)
	case solver.LevelEvil:
		// For evil games - make it even harder
		return r.PercentProbability(85)
	default:
		// For harder games, keep going until overflows...
		return true
	}
}
