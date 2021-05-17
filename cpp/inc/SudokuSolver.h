#pragma once

#include "SudokuBoard.h"
#include "SudokuResult.h"
#include "SudokuSolverOptions.h"

class SudokuSolver
{
  public:
    SudokuSolver() = default;
    SudokuSolver(const SudokuSolver &) = delete;
    SudokuSolver &operator=(const SudokuSolver &) = delete;

    SudokuResultConstShared run(const SudokuSolverOptions &options, SudokuBoardConstShared board)
    {
        return run(options, board, SudokuSolutionOptional::empty());
    }

    virtual SudokuResultConstShared run(
        const SudokuSolverOptions &options,
        const SudokuBoardConstShared inputBoard,
        SudokuSolutionOptional solution) = 0;
};

using SudokuSolverShared = std::shared_ptr<SudokuSolver>;

SudokuSolverShared createSolver();
