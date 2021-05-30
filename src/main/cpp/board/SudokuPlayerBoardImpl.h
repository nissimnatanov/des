#pragma once

#include "SudokuOptimizedBoardImpl.h"

class SudokuPlayerBoardImpl : public SudokuOptimizedBoardImpl, public virtual SudokuPlayerBoard
{
  public:
    SudokuPlayerBoardImpl() {}

    SudokuPlayerBoardImpl(const SudokuBoard &other) : SudokuOptimizedBoardImpl(other) {}

    ~SudokuPlayerBoardImpl() {}

    AccessMode getAccessMode() const override
    {
        return AccessMode::PLAY;
    }

    void playValue(int index, SudokuValue value) override
    {
        setValueInternal(index, value, /* isReadOnly= */ false);
    }

    void disallowValue(int index, SudokuValue value) override
    {
        disallowValues(index, SudokuValueSet(value));
    }

    void disallowValues(int index, SudokuValueSet values) override
    {
        addUserDisallowedValues(index, values);
    }

    void restart() override;
};
