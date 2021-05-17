#pragma once

#include <iostream>

enum class SudokuSolverAction
{
    /* Prove fast, level will be inaccurate. */
    PROVE,
    /* Solve fast, level will be inaccurate. */
    SOLVE_FAST,
    /* Solve and determine the level. */
    SOLVE,
    // HINT
};

inline const char *nameOf(const SudokuSolverAction &status)
{
    switch (status)
    {
    case SudokuSolverAction::PROVE:
        return "Prove";
    case SudokuSolverAction::SOLVE_FAST:
        return "FastSolve";
    case SudokuSolverAction::SOLVE:
        return "Solve";
    default:
        return "WRONG SudokuSolverAction";
    }
}

inline std::ostream &operator<<(std::ostream &os, const SudokuSolverAction &action)
{
    os << nameOf(action);
    return os;
}
