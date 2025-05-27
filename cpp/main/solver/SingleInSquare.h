#pragma once

#include "SudokuSolverAlgorithm.h"

class SingleInSquare : public SudokuSolverAlgorithm
{
  public:
    SudokuSolverStatus run(SudokuResultShared result) override;
};
