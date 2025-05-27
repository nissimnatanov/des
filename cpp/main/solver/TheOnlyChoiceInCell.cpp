#include "TheOnlyChoiceInCell.h"

SudokuSolverStatus TheOnlyChoiceInCell::run(SudokuResultShared result)
{
  SudokuPlayerBoard &board = *result->getPlayerBoard();
  for (int index = 0; index < 81; index++)
  {
    SudokuValue v = board.getValue(index);
    if (v != SudokuValue::EMPTY)
    {
      continue; // already set
    }

    SudokuValueSet allowedValues = board.getAllowedValues(index);

    if (allowedValues.size() > 1)
    {
      continue;
    }

    if (allowedValues.size() == 0)
    {
      // no value can be used for that specific cell
      return SudokuSolverStatus::NO_SOLUTION;
    }

    // exact one value can be used for this cell
    board.playValue(index, allowedValues[0]);
    result->report(SudokuStepComplexity::MEDIUM, SudokuStep::THE_ONLY_CHOICE);
    return SudokuSolverStatus::SUCCEEDED;
  }

  return SudokuSolverStatus::UNKNOWN;
}
