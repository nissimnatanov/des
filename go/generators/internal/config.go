package internal

import "github.com/nissimnatanov/des/go/solver"

// FastStageCount includes the last stage that is always a failure stage when fast
// generation fails
const FastStageCount = 4
const SlowStageCount = 10 // TODO: hardcode for now

const FastGenerationCap = solver.LevelVeryHard
