#pragma once

#include "SudokuSolverAlgorithm.h"

class TheOnlyChoiceInCell : public SudokuSolverAlgorithm
{
  public:
    SudokuSolverStatus run(SudokuResultShared result) override;
};
