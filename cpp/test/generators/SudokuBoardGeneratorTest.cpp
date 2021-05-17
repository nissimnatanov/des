#include <iostream>
#include <chrono>
#include <map>

#include "SudokuBoardGeneratorTest.h"
#include "SudokuSolver.h"
#include "asserts.h"

using namespace std;

struct LevelStats
{
    int total_count;
    int total_complexity;
    int total_value_count;

    LevelStats() : total_count(0), total_complexity(0), total_value_count(0) {}
};

void generate()
{
    map<SudokuLevel, LevelStats> allStats;

    int total = 0;
    auto start = chrono::steady_clock::now();
    auto last = start;
    int less_than_20 = 0;
    int timeouts = 0;
    while (true)
    {
        SudokuBoardGeneratorShared generator = newBoardGenerator();
        SudokuResultConstShared result = generator->generate(SudokuLevel::BLACKHOLE);
        if (result->getStatus() == SudokuSolverStatus::SUCCEEDED)
        {
            total++;
            LevelStats &stats = allStats[result->getLevel()];
            stats.total_count++;
            stats.total_complexity += result->getComplexity();
            stats.total_value_count += 81 - result->getOriginalBoard()->getFreeCellCount();

            bool print = (result->getLevel() >= SudokuLevel::NIGHTMARE);
            if (result->getOriginalBoard()->getFreeCellCount() > 61)
            {
                less_than_20++;
                print = true;
            }
            if (print)
            {
                cout << result->getOriginalBoard()->toString();
                cout << result;
            }
        }
        else if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
        {
            timeouts++;
            cout << "TIMEOUT detected: " << serializeBoard(result->getOriginalBoard().get()) << endl;
        }

        auto end = std::chrono::steady_clock::now();
        double elapsed_seconds = chrono::duration_cast<chrono::duration<double>>(end - last).count();
        if (elapsed_seconds > 60)
        {
            double elapsed_minutes_since_start =
                chrono::duration_cast<chrono::duration<double>>(end - start).count() / 60;
            last = start + chrono::minutes((int)elapsed_minutes_since_start);
            cerr << "Status [" << elapsed_minutes_since_start << "] { " << endl;
            cerr << "  All: total = " << total
                 << ", 19- = " << less_than_20
                 << ", timeouts = " << timeouts
                 << endl;
            for (auto s : allStats)
            {
                int complexity_average = 0;
                if (s.second.total_count != 0)
                {
                    complexity_average = s.second.total_complexity / s.second.total_count;
                }
                double value_count_average = s.second.total_value_count;
                if (s.second.total_value_count != 0)
                {
                    value_count_average /= s.second.total_count;
                }
                cerr << "  " << s.first << ": total = " << s.second.total_count
                     << ", complexity average " << complexity_average
                     << ", value count average " << value_count_average
                     << endl;
            }
            cerr << "}" << endl;
        }
    }
}

void runSudokuBoardGeneratorTests()
{
    cerr << "-------------------------------------" << endl;
    cerr << "Running SudokuBoardGenerator tests..." << endl;

    bool prevIntegrityChecks = setIntegrityChecks(false);
    generate();
    setIntegrityChecks(prevIntegrityChecks);
}
