package generators

import (
	"fmt"
	"slices"
	"time"

	"github.com/nissimnatanov/des/go/solver"
)

type GameStageStats struct {
	Total     int
	Succeeded int
	Failed    int
}

type GameStats struct {
	Count      int64
	Retries    int64
	Elapsed    time.Duration
	Complexity int64
	StageStats []GameStageStats
}

func (gs GameStats) clone() GameStats {
	c := gs
	c.StageStats = slices.Clone(gs.StageStats)
	return c
}

func (gs GameStats) AverageElapsed() time.Duration {
	if gs.Count == 0 {
		return 0
	}
	return gs.Elapsed / time.Duration(gs.Count)
}
func (gs GameStats) AverageRetries() float64 {
	if gs.Count == 0 {
		return 0
	}
	return float64(gs.Retries) / float64(gs.Count)
}

func (gs GameStats) AverageComplexity() float64 {
	if gs.Count == 0 {
		return 0
	}
	return float64(gs.Complexity) / float64(gs.Count)
}

func (gs GameStats) String() string {
	return fmt.Sprintf("Generations: %d, ~Elapsed: %s, ~Retries: %.3f, ~Complexity: %.3f, Stages: %v",
		gs.Count, gs.AverageElapsed(), gs.AverageRetries(), gs.AverageComplexity(), gs.StageStats)
}

func (gs *GameStats) reportOne(elapsed time.Duration, retries int64, complexity solver.StepComplexity, stageStags []GameStageStats) {
	gs.Count++
	gs.Elapsed += elapsed
	gs.Retries += retries
	gs.Complexity += int64(complexity)
	for i, stage := range stageStags {
		if i == len(gs.StageStats) {
			gs.StageStats = append(gs.StageStats, stage)
		} else {
			gs.StageStats[i].Total += stage.Total
			gs.StageStats[i].Succeeded += stage.Succeeded
			gs.StageStats[i].Failed += stage.Failed
		}
	}
}
