#pragma once

#include <iostream>

enum class SudokuSolverStatus
{
    UNKNOWN,
    SUCCEEDED,
    NO_SOLUTION,
    TWO_SOLUTIONS,
    LESS_THAN_17,
    // TWO_OR_MORE_VALUES_MISSING is different from TWO_SOLUTIONS since we do not try to solve at all,
    // and there might be no solution here at all. But if there is at least one: there is one more.
    TWO_OR_MORE_VALUES_MISSING,
    TIMEOUT,
};

inline const char *nameOf(const SudokuSolverStatus &status)
{
    switch (status)
    {
    case SudokuSolverStatus::UNKNOWN:
        return "Unknown";
    case SudokuSolverStatus::SUCCEEDED:
        return "Succeeded";
    case SudokuSolverStatus::NO_SOLUTION:
        return "No solution";
    case SudokuSolverStatus::TWO_SOLUTIONS:
        return "Two or more solutions";
    case SudokuSolverStatus::LESS_THAN_17:
        return "Less than 17";
    case SudokuSolverStatus::TIMEOUT:
        return "Timed out";
    case SudokuSolverStatus::TWO_OR_MORE_VALUES_MISSING:
        return "Two or more values missing";
    default:
        return "WRONG SudokuSolverStatus";
    }
}

inline std::ostream &operator<<(std::ostream &os, const SudokuSolverStatus &status)
{
    os << nameOf(status);
    return os;
}
