package internal

import (
	"fmt"
	"time"
)

type GamePerStageStats [SlowStageCount + 1]GameStageStats

func (stages *GamePerStageStats) Report(stage int, success bool) {
	if stage < 0 || stage >= len(stages) {
		panic("stage out of range")
	}
	for s := range stage + 1 {
		stages[s].Total++
	}
	if success {
		stages[stage].Succeeded++
	} else {
		stages[stage].Failed++
	}
}

func (stages *GamePerStageStats) merge(other GamePerStageStats) {
	for i := range stages {
		stages[i].Total += other[i].Total
		stages[i].Succeeded += other[i].Succeeded
		stages[i].Failed += other[i].Failed
	}
}

type GameStageStats struct {
	Total     int
	Succeeded int
	Failed    int
}

type GameStats struct {
	Count      int64
	Retries    int64
	Elapsed    time.Duration
	StageStats GamePerStageStats
}

func (gs GameStats) clone() GameStats {
	return gs
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

func (gs GameStats) String() string {
	return fmt.Sprintf("Generations: %d, ~Elapsed: %s, ~Retries: %.3f, Stages: %v",
		gs.Count, gs.AverageElapsed(), gs.AverageRetries(), gs.StageStats)
}

func (gs *GameStats) reportCount(count int, elapsed time.Duration, retries int64, stageStats GamePerStageStats) {
	gs.Count += int64(count)
	gs.Elapsed += elapsed
	gs.Retries += retries
	gs.StageStats.merge(stageStats)
}
