#pragma once

#include <bitset>
#include <array>
#include <stdexcept>
#include <string>
#include <sstream>
#include <iostream>

#include "SudokuBoard.h"
#include "SudokuIndexes.h"

class SudokuBoardImpl : public virtual SudokuBoard
{
  private:
    std::array<SudokuValue, BOARD_SIZE> _values = {};
    std::bitset<BOARD_SIZE> _readOnlyFlags;

    void appendValues(std::ostream &out) const;
    void appendRowValues(std::ostream &out) const;
    void appendColumnValues(std::ostream &out) const;
    void appendSquareValues(std::ostream &out) const;

    SudokuValue setValuePrivate(int index, SudokuValue value, bool isReadOnly);

    void checkIntegrity() const;

  protected:
    friend bool getIntegrityChecks();
    friend bool setIntegrityChecks(bool);
    static bool enableIntegrityChecks;

    SudokuBoardImpl() {}

    SudokuBoardImpl(const SudokuBoard &other);

    virtual SudokuValue setValueInternal(int index, SudokuValue value, bool isReadOnly)
    {
        return setValuePrivate(index, value, isReadOnly);
    }

    void assertIntegrity(bool condition) const
    {
        if (condition)
            return;

        throw std::logic_error("board has invalid state");
    }

  public:
    ~SudokuBoardImpl() {}

    SudokuValue getValue(int index) const override
    {
        return _values[index];
    }

    bool isReadOnlyCell(int index) const override
    {
        return _readOnlyFlags[index];
    }

    bool isEquivalentTo(SudokuBoardConstShared board) const override;

    bool containsAllValues(SudokuBoardConstShared board) const override;

    bool containsReadOnlyValues(SudokuBoardConstShared board) const override;

    std::string toString() const override
    {
        std::stringstream sstm;
        print(sstm, nullptr);
        return sstm.str();
    }

    void print(std::ostream &out, const char *format) const override;
};
