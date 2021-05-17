#include "SudokuPlayerBoardImpl.h"

using namespace std;

void SudokuPlayerBoardImpl::restart()
{
    suppressRecalculations();

    for (int i = 0; i < BOARD_SIZE; i++)
    {
        if (!isReadOnlyCell(i))
        {
            playValue(i, SudokuValue::EMPTY);
        }
    }

    resumeRecalculations();
}
