#include "SudokuSolverAlgorithm.h"

class SingleInColumn : public SudokuSolverAlgorithm
{
  public:
    SudokuSolverStatus run(SudokuResultShared result) override;
};
