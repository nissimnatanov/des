package generators

import "github.com/nissimnatanov/des/go/solver"

type stage struct {
	GeneratePerCandidate int `json:"generate_per_candidate,omitempty"`
	SelectBest           int `json:"select_best,omitempty"`
	FreeCells            int `json:"free_cells,omitempty"`
	MinToRemove          int `json:"min_to_remove,omitempty"`
	MaxToRemove          int `json:"max_to_remove,omitempty"`
	// TopN means we are interested in the top N sub-candidates for each candidate
	TopN              int          `json:"top_n,omitempty"`
	ProveOnlyLevelCap solver.Level `json:"prove_only_level_cap,omitempty"`
}

var slowStages = []stage{
	// at first we only have one candidate (e.g. solution), we can create many best forks
	{FreeCells: 49, MinToRemove: 10, MaxToRemove: 15, GeneratePerCandidate: 20, SelectBest: 10},
	{FreeCells: 51, TopN: 3, SelectBest: 50},
	{FreeCells: 55, MinToRemove: 2, MaxToRemove: 3, GeneratePerCandidate: 10, SelectBest: 50},
	{FreeCells: solver.MaxFreeCellsForValidBoard, TopN: 3, SelectBest: 100},
}
