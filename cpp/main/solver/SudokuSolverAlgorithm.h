#pragma once

#include "SudokuSolverStatus.h"
#include "SudokuResult.h"

class SudokuSolver;

class SudokuSolverAlgorithm
{
protected:
  virtual SudokuSolverStatus run(SudokuResultShared result)
  {
    throw std::logic_error("algorithm must implement either of the run methods");
  }

public:
  SudokuSolverAlgorithm() = default;
  SudokuSolverAlgorithm(const SudokuSolverAlgorithm &) = delete;
  SudokuSolverAlgorithm &operator=(const SudokuSolverAlgorithm &) = delete;

  virtual SudokuSolverStatus run(SudokuSolver *solver, SudokuResultShared result)
  {
    return run(result);
  }
};
