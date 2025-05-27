#include "SudokuSolverAlgorithm.h"
#include "SudokuIndexes.h"

class IdentifyPairs : public SudokuSolverAlgorithm
{
  private:
    int findPeer(SudokuPlayerBoard &board, SudokuValueSet vs, const related &indexes);
    bool tryEliminate(
        SudokuPlayerBoard &board, int ignore1, int ignore2, SudokuValueSet allowedValues, const sequence &indexes);

  public:
    SudokuSolverStatus run(SudokuResultShared result) override;
};
