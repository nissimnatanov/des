#pragma once

#include "SudokuBoardImpl.h"

class SudokuSolutionImpl : public SudokuBoardImpl, public virtual SudokuSolution
{
  public:
    SudokuSolutionImpl(const SudokuBoard &other) : SudokuBoardImpl(other)
    {
        if (!other.isSolved())
        {
            throw std::logic_error("solution can only be created from solved boards");
        }
    }

    ~SudokuSolutionImpl() {}

    SudokuSolutionImpl() = delete;
    SudokuSolutionImpl(const SudokuSolutionImpl &) = delete;
    SudokuSolutionImpl(SudokuSolutionImpl &&) = delete;
    SudokuSolutionImpl &operator=(const SudokuSolutionImpl &) = delete;
    SudokuSolutionImpl &operator=(SudokuSolutionImpl &&) = delete;

    AccessMode getAccessMode() const override
    {
        return AccessMode::IMMUTABLE;
    }

    SudokuValueSet getRowValues(int row) const override
    {
        return SudokuValueSet::all();
    }

    SudokuValueSet getColumnValues(int column) const override
    {
        return SudokuValueSet::all();
    }

    SudokuValueSet getSquareValues(int square) const override
    {
        return SudokuValueSet::all();
    }

    SudokuValueSet getAllowedValues(int index) const override
    {
        return SudokuValueSet::none();
    }

    virtual int getValueCount(SudokuValue v) const override
    {
        return (v == SudokuValue::EMPTY) ? 0 : SEQUENCE_SIZE;
    }

    bool isValidCell(int index) const override
    {
        return true;
    }

    bool isValid() const override
    {
        return true;
    }
};
