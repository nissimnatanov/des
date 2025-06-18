package stats

import (
	"fmt"
	"time"
)

// we currently only have ~8 stages in slow generators, so we can use a fixed size array
// with an extra capacity for future stages
type GameStages [12]GameStage

func (stages *GameStages) Report(count int, stage int) {
	if stage < 0 || stage >= len(*stages) {
		panic("negative stage")
	}
	if count > 0 {
		(*stages)[stage].Succeeded += count
	} else {
		(*stages)[stage].Failed++
		// failed loop also counts towards total
		count = 1
	}
	for s := range stage + 1 {
		(*stages)[s].Total += count
	}
}

func (stages *GameStages) merge(other GameStages) {
	for i := range other {
		(*stages)[i].Total += other[i].Total
		(*stages)[i].Succeeded += other[i].Succeeded
		(*stages)[i].Failed += other[i].Failed
	}
}

func (stages GameStages) String() string {
	// trim down zero stages
	last := len(stages) - 1
	for last >= 0 && stages[last].Total == 0 {
		last--
	}
	return fmt.Sprint(stages[:last+1])
}

type GameStage struct {
	Total     int
	Succeeded int
	Failed    int
}

type GameStats struct {
	Count      int64
	Retries    int64
	Elapsed    time.Duration
	StageStats GameStages
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
	return fmt.Sprintf("Game stats:\n* Generations: %d, ~Elapsed: %s, ~Retries: %.3f\n* Stages: %s",
		gs.Count, gs.AverageElapsed(), gs.AverageRetries(), gs.StageStats)
}

func (gs *GameStats) report(count int, elapsed time.Duration, retries int64, stageStats GameStages) {
	gs.Count += int64(count)
	gs.Elapsed += elapsed
	gs.Retries += retries
	gs.StageStats.merge(stageStats)
}
