#pragma once

#include <array>
#include <vector>
#include <map>
#include <memory>
#include <chrono>

#include "SudokuBoard.h"
#include "SudokuStep.h"
#include "SudokuStepComplexity.h"
#include "SudokuLevel.h"
#include "SudokuSolverStatus.h"
#include "SudokuSolverAction.h"
#include "SudokuSolverOptions.h"

using SteadyClock = std::chrono::steady_clock;
using SteadyTimePoint = std::chrono::time_point<SteadyClock>;

class SudokuResult
{
  private:
    SudokuBoardConstShared _originalBoard;
    SudokuPlayerBoardShared _playerBoard;
    const SudokuSolverOptions _options;
    unsigned long _complexity;
    std::array<std::array<int, sudokuStepSize()>, sizeOfSudokuStepComplexity()> _steps;
    SudokuSolverStatus _status;
    SteadyTimePoint _startTime, _endTime;

    std::vector<SudokuSolutionConstShared> _solutions;

  public:
    SudokuResult(const SudokuResult &) = delete;
    SudokuResult &operator=(const SudokuResult &) = delete;

    SudokuResult(SudokuBoardConstShared inputBoard, SudokuSolverOptions options)
        : _options(options), _complexity(0), _status(SudokuSolverStatus::UNKNOWN)
    {
        _originalBoard = cloneAsImmutable(inputBoard);
        _playerBoard = cloneToPlay(inputBoard);
        _startTime = SteadyClock::now();
        _endTime = SteadyTimePoint::max();

        for (int sci = 0; sci < _steps.size(); sci++)
        {
            auto &stepCounters = _steps.at(sci);
            stepCounters.fill(0);
        }
    }

    const SudokuSolverOptions getOptions() const noexcept
    {
        return _options;
    }

    SudokuBoardConstShared getOriginalBoard() const noexcept
    {
        return _originalBoard;
    }

    SudokuPlayerBoardShared getPlayerBoard() const noexcept
    {
        return _playerBoard;
    }

    auto elapsed() const noexcept
    {
        SteadyTimePoint end = _endTime;
        if (end == SteadyTimePoint::max())
        {
            end = SteadyClock::now();
        }
        return end - _startTime;
    }

    bool isTimedOut() const noexcept
    {
        return (elapsed() > _options.getMaxSolverTime());
    }

    void report(SudokuStepComplexity c, SudokuStep s)
    {
        report(c, s, 1);
    }

    /* SudokuStepComplexity must have operator*(int), indexOf() and nameOf() methods. Names must be UNIQUE! */
    void report(SudokuStepComplexity c, SudokuStep s, int count)
    {
        _complexity += c * count;

        auto &stepMap = _steps[indexOf(c)];
        stepMap[indexOf(s)] += count;
    }

    void merge(const SudokuResult &result);

    unsigned long getComplexity() const noexcept
    {
        return _complexity;
    }

    SudokuSolverStatus getStatus() const noexcept
    {
        return _status;
    }

    void complete(SudokuSolverStatus status) noexcept
    {
        _status = status;
        _endTime = SteadyClock::now();
    }

    bool addSolution(SudokuSolutionConstShared solution);

    const auto &getSolutions() const
    {
        return _solutions;
    }

    const auto &getSteps() const
    {
        return _steps;
    }

    SudokuLevel getLevel() const
    {
        if (getOptions().getAction() != SudokuSolverAction::SOLVE)
        {
            // Level is accurate only for SOLVE.
            return SudokuLevel::UNKNOWN;
        }
        return fromComplexity(getComplexity());
    }
};

std::ostream &operator<<(std::ostream &os, const SudokuResult &result);

using SudokuResultShared = std::shared_ptr<SudokuResult>;
using SudokuResultConstShared = std::shared_ptr<const SudokuResult>;

inline std::ostream &operator<<(std::ostream &os, const SudokuResult *result)
{
    return os << *result;
}
