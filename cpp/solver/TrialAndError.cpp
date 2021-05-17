#include <algorithm>

#include "TrialAndError.h"
#include "SudokuSolver.h"

using namespace std;

void reportRecursiveStep(SudokuResult &result)
{
    SudokuStepComplexity complexity;
    switch (result.getOptions().getCurrentRecursionDepth())
    {
    case 0:
        complexity = SudokuStepComplexity::RECURSION1;
        break;
    case 1:
        complexity = SudokuStepComplexity::RECURSION2;
        break;
    case 2:
        complexity = SudokuStepComplexity::RECURSION3;
        break;
    case 3:
    default:
        complexity = SudokuStepComplexity::RECURSION4;
        break;
    }
    result.report(complexity, SudokuStep::TRIAL_AND_ERROR);
}

SudokuSolverStatus TrialAndError::run(SudokuSolver *solver, SudokuResultShared result)
{
    SudokuPlayerBoard &board = *result->getPlayerBoard();
    if (result->getOptions().getCurrentRecursionDepth() >= result->getOptions().getMaxRecursionDepth())
    {
        // this algorithm activates recursion
        return SudokuSolverStatus::UNKNOWN;
    }

    vector<int> sortedIndexes;
    for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
    {
        if (board.isEmptyCell(i))
        {
            sortedIndexes.push_back(i);
        }
    }

    sort(sortedIndexes.begin(), sortedIndexes.end(),
         [board_p = &board](int o1, int o2)->bool {
             // Sort behavior may be different between machines if array has 'equal elements'. Make sure
             // there are no such to keep algorithm stable between different implementations of sort.
             int combined1 = board_p->getAllowedValues(o1).combined();
             int combined2 = board_p->getAllowedValues(o2).combined();

             return (combined1 == combined2) ? (o1 < o2) : (combined1 < combined2);
         });

    if (result->getOptions().getAction() == SudokuSolverAction::SOLVE &&
        result->getOptions().getCurrentRecursionDepth() == 0)
    {
        return useLayeredRecursion(solver, result, sortedIndexes);
    }
    else
    {
        SudokuSolverOptions childOptions(result->getOptions());
        childOptions.setCurrentRecursionDepth(childOptions.getCurrentRecursionDepth() + 1);
        if (result->getOptions().getAction() == SudokuSolverAction::PROVE)
        {
            childOptions.setAction(SudokuSolverAction::SOLVE_FAST);
        }
        return tryAndCheckAlgorithm(solver, result, sortedIndexes, childOptions);
    }
}

SudokuSolverStatus TrialAndError::useLayeredRecursion(
    SudokuSolver *solver,
    SudokuResultShared result,
    const std::vector<int> &sortedIndexes)
{
    SudokuSolverOptions childOptions(result->getOptions());
    childOptions.setCurrentRecursionDepth(childOptions.getCurrentRecursionDepth() + 1);

    for (int maxRecursionDepth = result->getOptions().getCurrentRecursionDepth() + 1;
         maxRecursionDepth <= result->getOptions().getMaxRecursionDepth();
         maxRecursionDepth++)
    {
        childOptions.setMaxRecursionDepth(maxRecursionDepth);
        SudokuSolverStatus status = tryAndCheckAlgorithm(solver, result, sortedIndexes, childOptions);
        if (status != SudokuSolverStatus::UNKNOWN)
        {
            return status;
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}

SudokuSolverStatus TrialAndError::tryAndCheckAlgorithm(
    SudokuSolver *solver,
    SudokuResultShared result,
    const std::vector<int> &sortedIndexes,
    const SudokuSolverOptions &childOptions)
{
    SudokuPlayerBoard &board = *result->getPlayerBoard();

    for (int index : sortedIndexes)
    {
        if (!board.isEmptyCell(index))
        {
            throw logic_error("cell was supposed to be empty");
        }

        SudokuValueSet allowedValues = board.getAllowedValues(index);
        for (SudokuValue testValue : allowedValues)
        {
            board.playValue(index, testValue);
            SudokuSolutionOptional solution =
                result->getSolutions().size() > 0
                    ? SudokuSolutionOptional(result->getSolutions().at(0))
                    : SudokuSolutionOptional::empty();
            SudokuResultConstShared testResult = solver->run(childOptions, result->getPlayerBoard(), solution);
            board.playValue(index, SudokuValue::EMPTY);

            if (testResult->getStatus() == SudokuSolverStatus::UNKNOWN)
            {
                continue;
            }
            else if (testResult->getStatus() == SudokuSolverStatus::TIMEOUT)
            {
                return SudokuSolverStatus::TIMEOUT;
            }

            reportRecursiveStep(*result);
            result->merge(*testResult);

            if (testResult->getStatus() == SudokuSolverStatus::NO_SOLUTION)
            {
                // when settings this value, the solution fails
                // thus, we can safely remove it
                board.disallowValue(index, testValue);
                return SudokuSolverStatus::SUCCEEDED;
            }

            if (testResult->getStatus() == SudokuSolverStatus::TWO_SOLUTIONS)
            {
                // no need to continue from this point if we detect more
                // than one solution is available
                for (SudokuSolutionConstShared newSolution : testResult->getSolutions())
                {
                    result->addSolution(newSolution);
                }
                return SudokuSolverStatus::TWO_SOLUTIONS;
            }

            if (testResult->getStatus() == SudokuSolverStatus::SUCCEEDED)
            {
                // Solution succeeded, remember it (if new).
                SudokuSolutionConstShared solution = testResult->getSolutions().at(0);
                result->addSolution(solution);

                if (result->getOptions().getAction() != SudokuSolverAction::PROVE)
                {
                    // if asked to solve, copy the values from
                    // the solution and stop
                    for (int ci = 0; ci < SudokuBoard::BOARD_SIZE; ci++)
                    {
                        if (board.isEmptyCell(ci))
                        {
                            board.playValue(ci, solution->getValue(ci));
                        }
                    }
                    return result->getSolutions().size() >= 2
                               ? SudokuSolverStatus::TWO_SOLUTIONS
                               : SudokuSolverStatus::SUCCEEDED;
                }
                else if (result->getOptions().getAction() == SudokuSolverAction::PROVE)
                {
                    if (result->getSolutions().size() >= 2)
                    {
                        // Got solution that is different from previously found or
                        // provided source board => two solutions are available.
                        return SudokuSolverStatus::TWO_SOLUTIONS;
                    }

                    continue;
                }
            }

            throw logic_error("unexpected state of TrialAndError");
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}
