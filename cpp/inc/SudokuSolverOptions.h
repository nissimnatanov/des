#pragma once

#include <chrono>

#include "SudokuSolverAction.h"

using namespace std::literals::chrono_literals;

class SudokuSolverOptions
{
  private:
    SudokuSolverAction _action;
    int _currentRecursionDepth, _maxRecursionDepth;
    std::chrono::seconds _maxSolverTime;

  public:
    SudokuSolverOptions() = delete;
    // OK with copy/move c-tor and equal operator.

    SudokuSolverOptions(SudokuSolverAction action)
        : _action(action),
          _currentRecursionDepth(0),
          _maxRecursionDepth(17),
          _maxSolverTime(10s)
    {
    }

    SudokuSolverAction getAction() const noexcept
    {
        return _action;
    }

    int getCurrentRecursionDepth() const noexcept
    {
        return _currentRecursionDepth;
    }

    int getMaxRecursionDepth() const noexcept
    {
        return _maxRecursionDepth;
    }

    std::chrono::seconds getMaxSolverTime() const noexcept
    {
        return _maxSolverTime;
    }

    SudokuSolverOptions &setAction(SudokuSolverAction action) noexcept
    {
        _action = action;
        return *this;
    }

    SudokuSolverOptions &setCurrentRecursionDepth(int depth) noexcept
    {
        _currentRecursionDepth = depth;
        return *this;
    }
    SudokuSolverOptions &setMaxRecursionDepth(int depth) noexcept
    {
        _maxRecursionDepth = depth;
        return *this;
    }
    void setMaxSolverTime(std::chrono::seconds maxTime) noexcept
    {
        _maxSolverTime = maxTime;
    }
};
