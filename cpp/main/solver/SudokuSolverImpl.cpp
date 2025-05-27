#include <array>
#include <memory>

#include "SudokuSolverImpl.h"
#include "SingleInSquare.h"
#include "SingleInRow.h"
#include "SingleInColumn.h"
#include "TheOnlyChoiceInCell.h"
#include "IdentifyPairs.h"
#include "TrialAndError.h"

using namespace std;

const shared_ptr<SudokuSolverAlgorithm> recursiveAlgorithm = make_shared<TrialAndError>();

const array<const shared_ptr<SudokuSolverAlgorithm>, 6> _algorithms = {
    make_shared<SingleInSquare>(),
    make_shared<SingleInRow>(),
    make_shared<SingleInColumn>(),
    make_shared<TheOnlyChoiceInCell>(),

    // elimination algorithms
    make_shared<IdentifyPairs>(),

    recursiveAlgorithm};

SudokuSolutionOptional SudokuSolverImpl::inferSolution(SudokuSolutionOptional solution, SudokuBoardConstShared board)
{
    if (solution.isPresent())
    {
        if (!solution.get()->containsReadOnlyValues(board))
        {
            throw logic_error("Given solution does not match the board to be solved.");
        }
        return solution;
    }
    else if (board->isSolved())
    {
        return SudokuSolutionOptional(cloneAsSolution(board));
    }
    else
    {
        return SudokuSolutionOptional::empty();
    }
}

SudokuResultConstShared SudokuSolverImpl::run(
    const SudokuSolverOptions &options,
    SudokuBoardConstShared inputBoard,
    SudokuSolutionOptional solutionOrNull)
{
    if (!inputBoard)
    {
        throw logic_error("Missing board.");
    }

    solutionOrNull = inferSolution(solutionOrNull, inputBoard);
    return runInternal(options, inputBoard, solutionOrNull);
}

SudokuResultConstShared SudokuSolverImpl::runInternal(
    const SudokuSolverOptions &options,
    SudokuBoardConstShared inputBoard,
    SudokuSolutionOptional solutionOrNull)
{
    SudokuResultShared result = make_shared<SudokuResult>(inputBoard, options);
    if (solutionOrNull.isPresent())
    {
        result->addSolution(solutionOrNull.get());
    }

    if (inputBoard->getFreeCellCount() > SudokuBoard::MAX_FREE_CELLS_FOR_VALID_BOARD)
    {
        // Boards with less than 17 values are mathematically proven to be wrong.
        result->complete(SudokuSolverStatus::LESS_THAN_17);
        return result;
    }

    int missing_digits = 0;
    for (SudokuValue v : SudokuValueSet::all())
    {
        if (inputBoard->getValueCount(v) == 0)
        {
            missing_digits++;
        }
    }
    if (missing_digits >= 2)
    {
        // There is no point to even try solving boards with two or more values missing.
        result->complete(SudokuSolverStatus::TWO_OR_MORE_VALUES_MISSING);
        return result;
    }

    while (result->getPlayerBoard()->getFreeCellCount() > 0)
    {
        SudokuSolverStatus stepStatus = runStep(result);
        if (stepStatus != SudokuSolverStatus::SUCCEEDED)
        {
            result->complete(stepStatus);
            break;
        }
    }

    // Check all solutions have the board values.
    for (SudokuSolutionConstShared resultSolution : result->getSolutions())
    {
        if (!resultSolution->containsReadOnlyValues(result->getPlayerBoard()))
        {
            throw logic_error("Result solution does not match the board to be solved.");
        }
    }

    // Check if we can infer solution from the board.
    if (result->getPlayerBoard()->isSolved())
    {
        // Looks like there is a solution in the board itself.
        result->addSolution(cloneAsSolution(result->getPlayerBoard()));
    }

    switch (result->getStatus())
    {
    case SudokuSolverStatus::UNKNOWN:
        // Board is solved, yet the solver did not set status yet. Infer it
        // from the solutions, if available.
        if (result->getPlayerBoard()->isSolved())
        {
            if (result->getSolutions().size() >= 2)
            {
                result->complete(SudokuSolverStatus::TWO_SOLUTIONS);
            }
            else if (result->getSolutions().size() == 1)
            {
                result->complete(SudokuSolverStatus::SUCCEEDED);
            }
            else
            {
                throw logic_error("solved board is added as a solution above");
            }
        }
        break;
    case SudokuSolverStatus::SUCCEEDED:
        if (!result->getPlayerBoard()->isSolved())
        {
            throw logic_error("status is success but board is not set");
        }
        if (result->getSolutions().size() >= 2)
        {
            // Solver can resolve board ignoring solution, and reach a different
            // one.
            // Override solver's decision.
            result->complete(SudokuSolverStatus::TWO_SOLUTIONS);
        }
        else if (result->getSolutions().size() == 0)
        {
            // See above: solved board is added as a solution.
            throw logic_error("solved board is added as a solution above");
        }
        break;
    case SudokuSolverStatus::TWO_SOLUTIONS:
        if (result->getSolutions().size() < 2)
        {
            throw logic_error("status is two solutions but second solution is missing");
        }
        break;

    default:
        break;
    }

    return result;
}

SudokuSolverStatus SudokuSolverImpl::runStep(SudokuResultShared result)
{
    if (result->getOptions().getAction() == SudokuSolverAction::SOLVE &&
        result->getOptions().getCurrentRecursionDepth() != 0 &&
        result->getOptions().getCurrentRecursionDepth() != result->getOptions().getMaxRecursionDepth())
    {
        // When we need to determine the level, BFS recursion is used (max recursion depth starts from 1
        // and then increases to the max - see TrialAndError::run.
        SudokuSolverStatus status = recursiveAlgorithm->run(this, result);
        if (status != SudokuSolverStatus::UNKNOWN)
        {
            return status;
        }
        if (result->isTimedOut())
        {
            return SudokuSolverStatus::TIMEOUT;
        }
    }
    else
    {
        for (auto algorithm : _algorithms)
        {
            SudokuSolverStatus status = algorithm->run(this, result);
            if (status != SudokuSolverStatus::UNKNOWN)
            {
                return status;
            }
            if (result->isTimedOut())
            {
                return SudokuSolverStatus::TIMEOUT;
            }
        }
    }

    return SudokuSolverStatus::UNKNOWN;
}
