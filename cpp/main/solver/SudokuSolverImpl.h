#pragma once

#include "SudokuSolver.h"
#include "SudokuResult.h"

class SudokuSolverImpl : public SudokuSolver
{
private:
  SudokuSolutionOptional inferSolution(SudokuSolutionOptional solution, SudokuBoardConstShared board);
  SudokuSolverStatus runStep(SudokuResultShared result);

protected:
  virtual SudokuResultConstShared runInternal(
      const SudokuSolverOptions &options,
      SudokuBoardConstShared inputBoard,
      SudokuSolutionOptional solution);

public:
  SudokuSolverImpl() {}

  SudokuResultConstShared run(
      const SudokuSolverOptions &options,
      SudokuBoardConstShared inputBoard,
      SudokuSolutionOptional solution) override;
};
