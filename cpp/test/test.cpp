#include <algorithm>
#include <string>
#include "asserts.h"
#include "board/SudokuIndexesTest.h"
#include "board/SudokuValueTest.h"
#include "board/SudokuBoardTest.h"
#include "generators/SudokuSolutionGeneratorTest.h"
#include "generators/SudokuBoardGeneratorTest.h"
#include "solver/SudokuSolverTest.h"

using namespace std;

void test(string testCase)
{
    transform(testCase.begin(), testCase.end(), testCase.begin(), ::tolower);
    if (testCase.empty() ||
        testCase.compare("indexes") == 0)
    {
        runSudokuIndexesTests();
    }
    if (testCase.empty() ||
        testCase.compare("value") == 0)
    {
        runSudokuValueTests();
    }
    if (testCase.empty() ||
        testCase.compare("board") == 0)
    {
        runSudokuBoardTests();
    }
    if (testCase.empty() ||
        testCase.compare("solution_generator") == 0)
    {
        runSudokuSolutionGeneratorTests();
    }
    if (testCase.empty() ||
        testCase.compare("solver") == 0)
    {
        runSudokuSolverTests();
    }
    if (testCase.empty() ||
        testCase.compare("board_generator") == 0)
    {
        runSudokuBoardGeneratorTests();
    }
}
std::map<int, std::map<int, int>> tt;

int main(int argc, char *argv[])
{
    cout << "Welcome to Dark Evil Sudoku Tests!" << std::endl;
    setIntegrityChecks(true);
    const char *testCase = "";
    if (argc > 1)
    {
        testCase = argv[2];
        if (argc > 3)
        {
            cerr << "Only one test case supported." << endl;
            return 1;
        }
    }
    try
    {
        test(testCase);
        cout << "Done!" << endl;
    }
    catch (assertion_error ae)
    {
        cerr << "Test failed!" << endl;
        cerr << "Test case: " << ae.testCase() << endl;
        cerr << "Location: " << ae.file() << ":" << ae.line() << endl;
    }
    return 0;
}