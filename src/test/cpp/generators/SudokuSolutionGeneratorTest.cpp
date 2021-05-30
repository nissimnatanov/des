#include <iostream>
#include <chrono>
#include "SudokuSolutionGeneratorTest.h"
#include "asserts.h"

using namespace std;

void singleBoard()
{
    SudokuSolutionGeneratorShared generator = newSolutionGenerator();
    Random r;
    SudokuSolutionConstShared b = generator->generate(r);
    b->print(cerr, "vrcs");
}

void timeTest()
{
    bool prevIntegrityChecks = setIntegrityChecks(false);
    const int board_count = 100;
    SudokuSolutionGeneratorShared generator = newSolutionGenerator();

    Random r;

    auto start = std::chrono::steady_clock::now();

    for (int i = 0; i < board_count; i++)
    {
        SudokuSolutionConstShared b = generator->generate(r);
    }

    auto end = std::chrono::steady_clock::now();
    double elapsed_seconds = chrono::duration_cast<chrono::duration<double>>(end - start).count();
    cerr << "Time to generate " << board_count << " solutions: "
              << elapsed_seconds << " seconds." << endl;

    setIntegrityChecks(prevIntegrityChecks);
}

void runSudokuSolutionGeneratorTests()
{
    cerr << "Running SudokuSolutionGenerator tests..." << endl;
    singleBoard();
    timeTest();
}
