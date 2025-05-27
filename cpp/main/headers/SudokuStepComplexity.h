#pragma once

#include <iostream>

enum class SudokuStepComplexity : unsigned long
{
    // Step complexity for simple algorithms, mapped to the weight of the step.
    EASY = 1,   // single in square
    MEDIUM = 5, // single in row/column
    HARD = 20,  // identify pairs

    // Step complexity for layered (BFS) recursion, mapped to the weight of the step.
    RECURSION1 = 100,    // single recursion
    RECURSION2 = 1000,   // second-level
    RECURSION3 = 10000,  // third-level - not reached yet with layered recursion
    RECURSION4 = 100000, // 4th-level and beyond
};

inline constexpr unsigned long operator*(const SudokuStepComplexity &step, unsigned int count) noexcept
{
    return static_cast<unsigned long>(step) * count;
}

inline const char *nameOf(const SudokuStepComplexity &stepComplexity)
{
    switch (stepComplexity)
    {
    case SudokuStepComplexity::EASY:
        return "Easy";
    case SudokuStepComplexity::MEDIUM:
        return "Medium";
    case SudokuStepComplexity::HARD:
        return "Hard";
    case SudokuStepComplexity::RECURSION1:
        return "Recursion1";
    case SudokuStepComplexity::RECURSION2:
        return "Recursion2";
    case SudokuStepComplexity::RECURSION3:
        return "Recursion3";
    case SudokuStepComplexity::RECURSION4:
        return "Recursion4";
    default:
        throw std::logic_error("WRONG SudokuStepComplexity");
    }
}

inline std::ostream &operator<<(std::ostream &os, const SudokuStepComplexity &stepComplexity)
{
    os << nameOf(stepComplexity);
    return os;
}

inline int indexOf(const SudokuStepComplexity &stepComplexity)
{
    switch (stepComplexity)
    {
    case SudokuStepComplexity::EASY:
        return 0;
    case SudokuStepComplexity::MEDIUM:
        return 1;
    case SudokuStepComplexity::HARD:
        return 2;
    case SudokuStepComplexity::RECURSION1:
        return 3;
    case SudokuStepComplexity::RECURSION2:
        return 4;
    case SudokuStepComplexity::RECURSION3:
        return 5;
    case SudokuStepComplexity::RECURSION4:
        return 6;
    default:
        throw std::logic_error("WRONG SudokuStepComplexity");
    }
}
inline SudokuStepComplexity stepComplexityFromIndex(int index)
{
    switch (index)
    {
    case 0:
        return SudokuStepComplexity::EASY;
    case 1:
        return SudokuStepComplexity::MEDIUM;
    case 2:
        return SudokuStepComplexity::HARD;
    case 3:
        return SudokuStepComplexity::RECURSION1;
    case 4:
        return SudokuStepComplexity::RECURSION2;
    case 5:
        return SudokuStepComplexity::RECURSION3;
    case 6:
        return SudokuStepComplexity::RECURSION4;
    default:
        throw std::logic_error("WRONG SudokuStepComplexity index");
    }
}

constexpr int sizeOfSudokuStepComplexity() noexcept
{
    return 7;
}
