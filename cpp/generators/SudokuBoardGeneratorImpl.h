#include <memory>
#include "SudokuBoardGeneratorState.h"
#include "SudokuBoardGenerator.h"

using namespace std;

class SudokuBoardGeneratorImpl : public SudokuBoardGenerator
{
public:
  struct Phase
  {
    const int freeCount;
    const int minComplexity;
    const bool enforceMinComplexity;
    const int generateCount;
    const int selectCount;
    const bool restore;
    const bool merge;
  };

private:
  // Fast generation (single board per solution).
  SudokuResultConstShared generate(Random &r, SudokuBoardGeneratorState &state);

  // Slow generation (multiple fast generators per solution, keeping the best N each iteration).
  SudokuResultConstShared generateSlow(Random &r, SudokuBoardGeneratorStateConstShared initialState);

  SudokuBoardGeneratorStateConstShared removePhase0(
      Random &r, SudokuBoardGeneratorStateConstShared initialState, int freeCount);
  vector<SudokuBoardGeneratorStateConstShared> generatePhase0Batch(
      Random &r,
      SudokuBoardGeneratorStateConstShared initialState,
      const Phase &phase);
  vector<SudokuBoardGeneratorStateConstShared> generatePhaseNBatch(
      int phi,
      Random &r,
      const vector<SudokuBoardGeneratorStateConstShared> &sources,
      const Phase &phase);

public:
  SudokuResultConstShared generate(SudokuLevel level, Random &r, SudokuSolutionOptional solutionOrNull) override;
};
