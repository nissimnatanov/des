#include <algorithm>
#include <string>
#include "asserts.h"
#include "board/SudokuIndexesTest.h"
#include "board/SudokuValueTest.h"
#include "board/SudokuBoardTest.h"
#include "generators/SudokuSolutionGeneratorTest.h"
#include "SudokuBoardBenchmark.h"

using namespace std;

int main(int argc, char *argv[])
{
    cout << "Welcome to Dark Evil Sudoku Tests!" << std::endl;
    setIntegrityChecks(true);

    try
    {
        runSudokuIndexesTests();
        runSudokuValueTests();
        runSudokuBoardTests();
        runSudokuSolutionGeneratorTests();
        runSudokuBoardGeneratorBenchmark(SudokuLevel::EVIL, 10);
        runSudokuBoardSamplesBenchmark();
        cout << "Done!" << endl;
    }
    catch (assertion_error ae)
    {
        cerr << "Test failed!" << endl;
        cerr << "Test case: " << ae.testCase() << endl;
        cerr << "Location: " << ae.file() << ":" << ae.line() << endl;
    }
    catch (logic_error le)
    {
        cerr << "Test failed!" << endl;
        cerr << "Exception: " << le.what() << endl;
    }
    catch (...) {
        cerr << "Got exception, cannot understand its type!" << endl;
    }
    // assert(greeter.greeting().compare("Hello, World!") == 0);
    return 0;
}
