#include "SingleInColumn.h"
#include "SudokuIndexes.h"

using namespace std;

SudokuSolverStatus SingleInColumn::run(SudokuResultShared result)
{
    SudokuPlayerBoard &board = *result->getPlayerBoard();
    for (int col = 0; col < 9; col++)
    {
        SudokuValueSet values = board.getColumnValues(col);
        if (values.size() == 9)
        {
            // full
            continue;
        }

        SudokuValueSet missingValues = ~values;
        for (SudokuValue v : missingValues)
        {
            int freeIndex = -1;

            for (int row = 0; row < 9; row++)
            {
                int index = indexFromCoordinates(row, col);
                SudokuValue temp = board.getValue(index);
                if (temp != SudokuValue::EMPTY)
                {
                    continue;
                }

                if (board.getAllowedValues(index).contains(v))
                {
                    if (freeIndex != -1)
                    {
                        // second free cell
                        freeIndex = -2;
                        break;
                    }

                    freeIndex = index;
                }
            }

            if (freeIndex == -1)
            {
                // no free cell in column - fail the solution
                return SudokuSolverStatus::NO_SOLUTION;
            }

            if (freeIndex == -2)
            {
                // duplicate found
                continue;
            }

            // found free cell that is single in its column
            board.playValue(freeIndex, v);
            result->report(SudokuStepComplexity::EASY, SudokuStep::SINGLE_IN_COLUMN);
            return SudokuSolverStatus::SUCCEEDED;
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}