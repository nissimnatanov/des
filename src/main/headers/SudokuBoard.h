#pragma once

#include <memory>
#include <stdexcept>
#include <string>
#include <iostream>

#include "SudokuValue.h"

class SudokuPlayerBoard;
class SudokuEditBoard;
class SudokuSolution;

class SudokuBoard
{
  public:
    SudokuBoard(const SudokuBoard &) = delete;
    SudokuBoard &operator=(const SudokuBoard &) = delete;

    static const int BOARD_SIZE = 81;
    static const int SEQUENCE_SIZE = 9; // length of the single row/column/square.
    // boards with less than 17 are mathematically proven to be illegal sudoku boards.
    static const int MAX_FREE_CELLS_FOR_VALID_BOARD = 64;

    enum class AccessMode
    {
        IMMUTABLE,
        PLAY,
        EDIT,
    };

    ~SudokuBoard() {}

    virtual AccessMode getAccessMode() const = 0;

    virtual SudokuValue getValue(int index) const = 0;

    bool isEmptyCell(int index) const
    {
        return getValue(index) == SudokuValue::EMPTY;
    }

    virtual bool isReadOnlyCell(int index) const = 0;

    virtual SudokuValueSet getRowValues(int row) const = 0;

    virtual SudokuValueSet getColumnValues(int column) const = 0;

    virtual SudokuValueSet getSquareValues(int square) const = 0;

    virtual SudokuValueSet getAllowedValues(int index) const = 0;

    bool isAllowedValue(int index, SudokuValue value) const
    {
        if (value == SudokuValue::EMPTY)
        {
            throw std::logic_error("Do not use empty value with isAllowedValue");
        }

        return getAllowedValues(index).contains(value);
    }

    SudokuValueSet getDisallowedValues(int index) const noexcept
    {
        return ~getAllowedValues(index);
    }

    virtual bool isEquivalentTo(std::shared_ptr<const SudokuBoard> board) const = 0;

    virtual bool containsAllValues(std::shared_ptr<const SudokuBoard> board) const = 0;

    virtual bool containsReadOnlyValues(std::shared_ptr<const SudokuBoard> board) const = 0;

    virtual int getValueCount(SudokuValue v) const = 0;

    int getFreeCellCount() const noexcept
    {
        return getValueCount(SudokuValue::EMPTY);
    }

    virtual bool isValidCell(int index) const = 0;

    virtual bool isValid() const = 0;

    virtual bool isSolved() const noexcept
    {
        return getFreeCellCount() == 0 && isValid();
    }

    virtual std::string toString() const = 0;

    virtual void print(std::ostream &out, const char *format) const = 0;

  protected:
    SudokuBoard() {}
};

inline std::ostream &operator<<(std::ostream &out, const SudokuBoard &board)
{
    board.print(out, "v");
    return out;
}

inline std::ostream &operator<<(std::ostream &out, const SudokuBoard *board)
{
    out << *board;
    return out;
}

class SudokuSolution : public virtual SudokuBoard
{
  public:
    SudokuSolution(const SudokuSolution &) = delete;
    SudokuSolution &operator=(const SudokuSolution &) = delete;

    ~SudokuSolution() {}

  protected:
    SudokuSolution() {}
};

class SudokuPlayerBoard : public virtual SudokuBoard
{
  public:
    SudokuPlayerBoard(const SudokuPlayerBoard &) = delete;
    SudokuPlayerBoard &operator=(const SudokuPlayerBoard &) = delete;

    ~SudokuPlayerBoard() {}

    virtual void playValue(int index, SudokuValue value) = 0;

    virtual void disallowValue(int index, SudokuValue value) = 0;

    virtual void disallowValues(int index, SudokuValueSet valueSet) = 0;

    virtual void restart() = 0;

  protected:
    SudokuPlayerBoard() {}
};

class SudokuEditBoard : public virtual SudokuPlayerBoard
{
  public:
    SudokuEditBoard(const SudokuEditBoard &) = delete;
    SudokuEditBoard &operator=(const SudokuEditBoard &) = delete;

    ~SudokuEditBoard() {}

    /** Sets read-only value. To set read-write value, use playValue. */
    virtual SudokuValue setValue(int index, SudokuValue value) = 0;

  protected:
    SudokuEditBoard() {}
};

using SudokuBoardShared = std::shared_ptr<SudokuBoard>;
using SudokuBoardConstShared = std::shared_ptr<const SudokuBoard>;
using SudokuEditBoardShared = std::shared_ptr<SudokuEditBoard>;
using SudokuEditBoardConstShared = std::shared_ptr<const SudokuEditBoard>;
using SudokuPlayerBoardShared = std::shared_ptr<SudokuPlayerBoard>;
using SudokuPlayerBoardConstShared = std::shared_ptr<const SudokuPlayerBoard>;
// Solution is always const.
using SudokuSolutionConstShared = std::shared_ptr<const SudokuSolution>;

bool getIntegrityChecks();
bool setIntegrityChecks(bool enable);

SudokuEditBoardShared newEditBoard();

const SudokuBoardConstShared cloneAsImmutable(SudokuBoardConstShared other);
SudokuPlayerBoardShared cloneToPlay(SudokuBoardConstShared other);
SudokuEditBoardShared cloneToEdit(SudokuBoardConstShared other);
const SudokuSolutionConstShared cloneAsSolution(SudokuBoardConstShared other);

inline std::ostream &operator<<(std::ostream &os, const SudokuBoard::AccessMode &am)
{
    switch (am)
    {
    case SudokuBoard::AccessMode::EDIT:
        os << "EDIT";
        break;
    case SudokuBoard::AccessMode::IMMUTABLE:
        os << "IMMUTABLE";
        break;
    case SudokuBoard::AccessMode::PLAY:
        os << "PLAY";
        break;
    default:
        throw std::logic_error("WRONG SudokuBoard::AccessMode");
    }
    return os;
}

class SudokuSolutionOptional
{
  private:
    SudokuSolutionConstShared solutionOrNull;

  public:
    SudokuSolutionOptional(const SudokuSolutionConstShared &solution) : solutionOrNull(solution) {}
    SudokuSolutionOptional() : solutionOrNull(nullptr) {}

    bool isPresent() const noexcept
    {
        return solutionOrNull != nullptr;
    }

    SudokuSolutionConstShared get() const
    {
        if (!isPresent())
        {
            throw std::logic_error("No solution is present .. use isPresent().");
        }
        return solutionOrNull;
    }

    static SudokuSolutionOptional empty() noexcept
    {
        return SudokuSolutionOptional();
    }
};

std::string serializeBoard(const SudokuBoard *board);
SudokuEditBoardShared deserializeSudokuBoard(std::string input);
