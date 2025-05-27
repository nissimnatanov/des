#pragma once

#include "SudokuPlayerBoardImpl.h"

class SudokuEditBoardImpl : public SudokuPlayerBoardImpl, public virtual SudokuEditBoard
{
  public:
    SudokuEditBoardImpl() {}
    ~SudokuEditBoardImpl() {}

    SudokuEditBoardImpl(const SudokuBoard &other) : SudokuPlayerBoardImpl(other) {}

    AccessMode getAccessMode() const override
    {
        return AccessMode::EDIT;
    }

    SudokuValue setValue(int index, SudokuValue value) override
    {
        if (value == SudokuValue::EMPTY)
        {
            return setValueInternal(index, value, /* isReadOnly= */ false);
        }

        else
        {
            return setValueInternal(index, value, /* isReadOnly= */ true);
        }
    }
};
