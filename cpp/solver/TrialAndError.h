#pragma once

#include "SudokuSolverAlgorithm.h"

class TrialAndError : public SudokuSolverAlgorithm
{
private:
  SudokuSolverStatus tryAndCheckAlgorithm(
      SudokuSolver *solver,
      SudokuResultShared result,
      const std::vector<int> &sortedIndexes,
      const SudokuSolverOptions &childOptions);
  SudokuSolverStatus useLayeredRecursion(
      SudokuSolver *solver,
      SudokuResultShared result,
      const std::vector<int> &sortedIndexes);

public:
  SudokuSolverStatus run(SudokuSolver *solver, SudokuResultShared result) override;
};
