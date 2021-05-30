#include "SingleInSquare.h"
#include "SudokuIndexes.h"

using namespace std;

SudokuSolverStatus SingleInSquare::run(SudokuResultShared result)
{
    SudokuPlayerBoard &board = *result->getPlayerBoard();
    for (int square = 0; square < 9; square++)
    {
        SudokuValueSet values = board.getSquareValues(square);
        if (values.size() == 9)
        {
            // full
            continue;
        }

        SudokuValueSet missingValues = ~values;
        for (SudokuValue v : missingValues)
        {
            int freeIndex = -1;

            for (int cell = 0; cell < 9; cell++)
            {
                int index = indexFromSquare(square, cell);
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
                // no free cell in square - fail the solution
                return SudokuSolverStatus::NO_SOLUTION;
            }

            if (freeIndex == -2)
            {
                // duplicate found
                continue;
            }

            // found free cell that is single in its square
            board.playValue(freeIndex, v);
            result->report(SudokuStepComplexity::EASY, SudokuStep::SINGLE_IN_SQUARE);
            return SudokuSolverStatus::SUCCEEDED;
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}