package stats

import (
	"fmt"
	"strings"
	"time"
)

// we currently only have ~8 stages in slow generators, so we can use a fixed size array
// with an extra capacity for future stages
type GameStages [12]GameStage

func (stages *GameStages) ReportCandidateCount(stage int, count int) {
	if stage < 0 || stage >= len(*stages) {
		panic("negative stage")
	}
	if count < 0 {
		panic("negative count")
	}
	(*stages)[stage].CandidateCount += 1
	(*stages)[stage].Candidate += int64(count)
}

func (stages *GameStages) ReportBestComplexity(stage int, complexity int64) {
	if stage < 0 || stage >= len(*stages) {
		panic("negative stage")
	}
	if complexity < 0 {
		panic("negative complexity")
	}
	if complexity == 0 {
		return
	}
	(*stages)[stage].TotalComplexity += complexity
	(*stages)[stage].BestComplexity = max((*stages)[stage].BestComplexity, complexity)
	(*stages)[stage].ComplexityCount++
}

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
		(*stages)[i].TotalComplexity += other[i].TotalComplexity
		(*stages)[i].ComplexityCount += other[i].ComplexityCount
		(*stages)[i].BestComplexity = max((*stages)[i].BestComplexity, other[i].BestComplexity)
		(*stages)[i].Candidate += other[i].Candidate
		(*stages)[i].CandidateCount += other[i].CandidateCount
	}
}

func (stages GameStages) String() string {
	// trim down zero stages
	last := len(stages) - 1
	for last >= 0 && stages[last].Total == 0 {
		last--
	}
	var sb strings.Builder
	sb.WriteString("[")
	complexities := make([]string, 0, last+1)
	for si, s := range stages[:last+1] {
		if s.ComplexityCount > 0 {
			complexityAve := s.TotalComplexity / int64(s.ComplexityCount)
			complexities = append(complexities, fmt.Sprintf("[%d] %d/%d", si, s.BestComplexity, complexityAve))
		} else {
			complexities = append(complexities, fmt.Sprintf("[%d] 0", si))
		}
		if s.Succeeded == 0 && s.Failed == 0 {
			// skip non-productive stages
			continue
		}
		var successAve, failAve, candidateAve float64
		if s.Total > 0 {
			successAve = float64(s.Succeeded) / float64(s.Total) * 100
			failAve = float64(s.Failed) / float64(s.Total) * 100
		}
		if s.CandidateCount > 0 {
			candidateAve = float64(s.Candidate) / float64(s.CandidateCount)
		}

		sb.WriteString(fmt.Sprintf("\n  [%d] Total: %d, Success: %d(%.02f%%), Failed: %d(%.02f%%), ~Candidates: %.01f",
			si, s.Total, s.Succeeded, successAve, s.Failed, failAve, candidateAve))
	}
	sb.WriteString("\n], Complexities: ")
	sb.WriteString(strings.Join(complexities, " "))
	return sb.String()
}

type GameStage struct {
	Total           int
	Succeeded       int
	Failed          int
	TotalComplexity int64
	BestComplexity  int64 // best complexity for this stage
	ComplexityCount int
	Candidate       int64
	CandidateCount  int64
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
	return fmt.Sprintf("Game Generations: %d, ~Elapsed: %s, ~Retries: %.3f. Stages: %s",
		gs.Count, gs.AverageElapsed(), gs.AverageRetries(), gs.StageStats)
}

func (gs *GameStats) report(count int, elapsed time.Duration, retries int64, stageStats GameStages) {
	gs.Count += int64(count)
	gs.Elapsed += elapsed
	gs.Retries += retries
	gs.StageStats.merge(stageStats)
}
