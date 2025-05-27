#pragma once

#include "SudokuStepComplexity.h"

enum class SudokuLevel
{
    UNKNOWN,
    EASY,
    MEDIUM,
    HARD,
    VERY_HARD,
    EVIL,
    DARK_EVIL,
    NIGHTMARE,
    BLACKHOLE,
};

enum class SudokuLevelBar : unsigned long
{
    EASY_BAR = 100,
    MEDIUM_BAR = 250,
    HARD_BAR = 500,
    VERY_HARD_BAR = 1000,
    EVIL_BAR = 2000,
    DARK_EVIL_BAR = 5000,
    NIGHTMARE_BAR = 15000,
};

inline bool operator<=(unsigned long complexity, SudokuLevelBar bar)
{
    return complexity <= static_cast<unsigned long>(bar);
}

inline SudokuLevel fromComplexity(unsigned long complexity)
{
    if (complexity <= SudokuLevelBar::EASY_BAR)
    {
        return SudokuLevel::EASY;
    }
    else if (complexity <= SudokuLevelBar::MEDIUM_BAR)
    {
        return SudokuLevel::MEDIUM;
    }
    else if (complexity <= SudokuLevelBar::HARD_BAR)
    {
        return SudokuLevel::HARD;
    }
    else if (complexity <= SudokuLevelBar::VERY_HARD_BAR)
    {
        return SudokuLevel::VERY_HARD;
    }
    else if (complexity <= SudokuLevelBar::EVIL_BAR)
    {
        return SudokuLevel::EVIL;
    }
    else if (complexity <= SudokuLevelBar::DARK_EVIL_BAR)
    {
        return SudokuLevel::DARK_EVIL;
    }
    else if (complexity <= SudokuLevelBar::NIGHTMARE_BAR)
    {
        return SudokuLevel::NIGHTMARE;
    }
    else
    {
        return SudokuLevel::BLACKHOLE;
    }
}

inline const char *nameOf(const SudokuLevel &level)
{
    switch (level)
    {
    case SudokuLevel::UNKNOWN:
        return "Unknown";
    case SudokuLevel::EASY:
        return "Easy";
    case SudokuLevel::MEDIUM:
        return "Medium";
    case SudokuLevel::HARD:
        return "Hard";
    case SudokuLevel::VERY_HARD:
        return "VeryHard";
    case SudokuLevel::EVIL:
        return "Evil";
    case SudokuLevel::DARK_EVIL:
        return "DarkEvil";
    case SudokuLevel::NIGHTMARE:
        return "Nightmare";
    case SudokuLevel::BLACKHOLE:
        return "BlackHole";
    default:
        return "WRONG SudokuLevel";
    }
}

inline std::ostream &operator<<(std::ostream &os, const SudokuLevel &level)
{
    os << nameOf(level);
    return os;
}
