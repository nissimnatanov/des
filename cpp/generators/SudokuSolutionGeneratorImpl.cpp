#include <vector>

#include "SudokuIndexes.h"
#include "SudokuValue.h"
#include "SudokuSolutionGeneratorImpl.h"

using namespace std;

SudokuSolutionGeneratorShared newSolutionGenerator()
{
    return make_shared<SudokuSolutionGeneratorImpl>();
}

void setSquareValues(
    SudokuEditBoard &board,
    int square,
    vector<SudokuValue> &values)
{
  int cellIndex = 0;
  for (auto v : values)
  {
    board.setValue(indexFromSquare(square, cellIndex), v);
    cellIndex++;
  }

  if (!board.isValid())
  {
    throw logic_error("setSquareValues generated invalid board");
  }
}

bool tryFillSquare(Random &r, SudokuEditBoard &board, int square)
{
  array<int, 9> cellIndexes(getSequenceIndexes());

  // if there is at least one cell with no allowed values, no point to conitnue
  for (int cell : cellIndexes)
  {
      if (board.getAllowedValues(indexFromSquare(square, cell)).size() == 0)
      {
        return false;
      }
  }

  r.shuffle(cellIndexes);

  int allowedRetries = 4; // optimized for best performance.
  bool retry;
  do
  {
    retry = false;
    for (int cell : cellIndexes)
    {
      int index = indexFromSquare(square, cell);
      SudokuValueSet allowedValues = board.getAllowedValues(index);
      bool valid = (allowedValues.size() > 0);
      if (valid)
      {
        int next = r.nextIndex(allowedValues.size());
        SudokuValue value = allowedValues[next];
        board.playValue(index, value);
        if (!board.isValid())
        {
          valid = false;
        }
      }

      if (!valid)
      {
        if (allowedRetries-- == 0)
        {
          return false;
        }

        // reset the square
        for (cell = 0; cell < 9; cell++)
        {
          board.playValue(indexFromSquare(square, cell), SudokuValue::EMPTY);
        }

        retry = true;
        break;
      }
    }
  } while (retry);

  return true;
}

SudokuSolutionConstShared SudokuSolutionGeneratorImpl::generate(Random &r)
{
  SudokuEditBoardShared board = newEditBoard();

  vector<SudokuValue> values(SudokuValueSet::all().cbegin(), SudokuValueSet::all().cend());

  // populate squares 0, 4 (middle) and 8 (last)
  // since there is no intersection between those, safe to use totally
  // random values
  r.shuffle(values);
  setSquareValues(*board, 0, values);
  r.shuffle(values);
  setSquareValues(*board, 4, values);
  r.shuffle(values);
  setSquareValues(*board, 8, values);

  do
  {
    board->restart();

    // fill squares 1 to 7, skipping 4
    for (int square = 1; square <= 7; square++)
    {
      if (square == 4)
        continue; // skip

      if (!tryFillSquare(r, *board, square))
      {
        break;
      }
    }
  } while (board->getFreeCellCount() > 0);

  // reset all values as read-only.
  for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
  {
    board->setValue(i, board->getValue(i));
  }

  return cloneAsSolution(board);
}
