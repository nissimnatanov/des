#pragma once

#include <iostream>

#include "SudokuStepComplexity.h"

// Step (or algorithm).
enum class SudokuStep
{
    SINGLE_IN_SQUARE,
    SINGLE_IN_ROW,
    SINGLE_IN_COLUMN,
    THE_ONLY_CHOICE,
    IDENTIFY_PAIRS,
    TRIAL_AND_ERROR,
};

inline const char *nameOf(const SudokuStep &step)
{
    switch (step)
    {
    case SudokuStep::SINGLE_IN_SQUARE:
        return "Single In Square";
    case SudokuStep::SINGLE_IN_ROW:
        return "Single In Row";
    case SudokuStep::SINGLE_IN_COLUMN:
        return "Single In Column";
    case SudokuStep::THE_ONLY_CHOICE:
        return "The Only Choice";
    case SudokuStep::IDENTIFY_PAIRS:
        return "Identify Pairs";
    case SudokuStep::TRIAL_AND_ERROR:
        return "Trial & Error";
    default:
        throw std::logic_error("WRONG SudokuStep");
    }
}

inline int indexOf(const SudokuStep &step)
{
    switch (step)
    {
    case SudokuStep::SINGLE_IN_SQUARE:
        return 0;
    case SudokuStep::SINGLE_IN_ROW:
        return 1;
    case SudokuStep::SINGLE_IN_COLUMN:
        return 2;
    case SudokuStep::THE_ONLY_CHOICE:
        return 3;
    case SudokuStep::IDENTIFY_PAIRS:
        return 4;
    case SudokuStep::TRIAL_AND_ERROR:
        return 5;
    default:
        throw std::logic_error("WRONG SudokuStep");
    }
}

inline SudokuStep sudokuStepFromIndex(int index)
{
    switch (index)
    {
    case 0:
        return SudokuStep::SINGLE_IN_SQUARE;
    case 1:
        return SudokuStep::SINGLE_IN_ROW;
    case 2:
        return SudokuStep::SINGLE_IN_COLUMN;
    case 3:
        return SudokuStep::THE_ONLY_CHOICE;
    case 4:
        return SudokuStep::IDENTIFY_PAIRS;
    case 5:
        return SudokuStep::TRIAL_AND_ERROR;
    default:
        throw std::logic_error("WRONG SudokuStep");
    }
}

constexpr int sudokuStepSize() noexcept
{
    return 6;
}

inline std::ostream &operator<<(std::ostream &os, const SudokuStep &step)
{
    os << nameOf(step);
    return os;
}
