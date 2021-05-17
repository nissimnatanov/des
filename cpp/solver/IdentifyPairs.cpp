#include "IdentifyPairs.h"
#include "SudokuIndexes.h"

using namespace std;

SudokuSolverStatus IdentifyPairs::run(SudokuResultShared result)
{
    SudokuPlayerBoard &board = *result->getPlayerBoard();
    for (int index = 0; index < 81; index++)
    {
        if (!board.isEmptyCell(index))
        {
            continue;
        }

        SudokuValueSet allowedValues = board.getAllowedValues(index);
        if (allowedValues.size() != 2)
        {
            continue;
        }

        int peerIndex = findPeer(board, allowedValues, getRelated(index));
        if (peerIndex < 0)
        {
            continue;
        }

        // found peer
        int eliminationCount = 0;
        int row = rowFromIndex(index);
        if (row == rowFromIndex(peerIndex))
        {
            // same row
            if (tryEliminate(board, index, peerIndex, allowedValues, getRowIndexes(row)))
            {
                eliminationCount++;
            }
        }

        int col = columnFromIndex(index);
        if (col == columnFromIndex(peerIndex))
        {
            // same column
            if (tryEliminate(board, index, peerIndex, allowedValues, getColumnIndexes(col)))
            {
                eliminationCount++;
            }
        }

        int square = squareFromIndex(index);
        if (square == squareFromIndex(peerIndex))
        {
            // same square
            if (tryEliminate(board, index, peerIndex, allowedValues, getSquareIndexes(square)))
            {
                eliminationCount++;
            }
        }

        if (eliminationCount > 0)
        {
            result->report(SudokuStepComplexity::HARD, SudokuStep::IDENTIFY_PAIRS);
            return SudokuSolverStatus::SUCCEEDED;
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}

bool IdentifyPairs::tryEliminate(
    SudokuPlayerBoard &board, int ignore1, int ignore2, SudokuValueSet allowedValues, const sequence &indexes)
{
    bool found = false;
    for (int temp = 0; temp < 9; temp++)
    {
        int index = indexes.at(temp);
        if (index == ignore1 || index == ignore2)
        {
            continue;
        }

        if (!board.isEmptyCell(index))
        {
            continue;
        }

        SudokuValueSet tempAllowedValues = board.getAllowedValues(index);
        if ((tempAllowedValues & allowedValues).size() > 0)
        {
            // found a cell that we can remove values - turn them off
            board.disallowValues(index, allowedValues);
            found = true;
        }
    }

    return found;
}

int IdentifyPairs::findPeer(SudokuPlayerBoard &board, SudokuValueSet allowedValues, const related &indexes)
{
    for (int peerIndex : indexes)
    {
        if (!board.isEmptyCell(peerIndex))
        {
            continue;
        }
        SudokuValueSet peerAllowedValues = board.getAllowedValues(peerIndex);
        if (peerAllowedValues == allowedValues)
        {
            return peerIndex;
        }
    }
    return -1;
}
