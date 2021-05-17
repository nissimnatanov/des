#include "SudokuEditBoardImpl.h"

using namespace std;

SudokuEditBoardShared newEditBoard()
{
    return make_shared<SudokuEditBoardImpl>();
}
