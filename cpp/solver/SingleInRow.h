#pragma once

#include "SudokuSolverAlgorithm.h"

class SingleInRow : public SudokuSolverAlgorithm
{
  public:
    SudokuSolverStatus run(SudokuResultShared result) override;
};
