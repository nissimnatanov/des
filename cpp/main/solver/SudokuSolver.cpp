#include "SudokuSolverImpl.h"

using namespace std;

SudokuSolverShared createSolver()
{
    return make_shared<SudokuSolverImpl>();
}
