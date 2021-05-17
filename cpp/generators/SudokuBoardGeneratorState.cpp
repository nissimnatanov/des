#include "SudokuBoardGeneratorState.h"
#include "SudokuSolutionGenerator.h"

bool shouldContinueByLevel(SudokuLevel desiredLevel, Random &r)
{
    switch (desiredLevel)
    {
    case SudokuLevel::EASY:
        // For easy games - keep trying (otherwise, game can be too easy).
        return r.percentProbability(95);
    case SudokuLevel::MEDIUM:
        // For medium games - keep trying a bit less.
        return r.percentProbability(75);
    case SudokuLevel::HARD:
        // For hard games - continue in half of the cases..
        return r.percentProbability(50);
    case SudokuLevel::VERYHARD:
        // For very hard games - make it even harder, but stop sometimes.
        return r.percentProbability(75);
    default:
        // For harder games, keep going until overflows...
        return true;
    }
}

SudokuBoardGeneratorState::SudokuBoardGeneratorState(
    SudokuLevel level, Random &r, SudokuSolutionOptional solutionOrNull, SudokuSolverShared solver)
    : _level(level), _solver(solver), _proven(false)
{
    if (solutionOrNull.isPresent())
    {
        _solution = solutionOrNull.get();
    }
    else
    {
        SudokuSolutionGeneratorShared solutionGenerator = newSolutionGenerator();
        _solution = solutionGenerator->generate(r);
    }

    _board = cloneToEdit(_solution);
    // Mark all cells as read-only.
    for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
    {
        _board->setValue(i, _board->getValue(i));
    }

    checkIntegrity();
}

SudokuBoardGeneratorState::SudokuBoardGeneratorState(
    SudokuLevel level, Random &r, SudokuSolverShared solver, SudokuBoardConstShared board)
    : _level(level), _solver(solver), _proven(false)
{
    SudokuSolverOptions options(SudokuSolverAction::SOLVE);
    auto res = solver->run(options, board);
    if (res->getStatus() != SudokuSolverStatus::SUCCEEDED)
    {
        throw logic_error("Cannot fine-tune board with no solution.");
    }
    _solution = res->getSolutions().at(0);
    _board = cloneToEdit(board);
    // Mark all cells as read-only.
    for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
    {
        if (_board->isEmptyCell(i))
        {
            // remove empty cell from indexes
            _indexes.removeIndex(i);
        }
        else
        {
            _board->setValue(i, _board->getValue(i));
        }
    }

    checkIntegrity();
}

SudokuBoardGeneratorState::SudokuBoardGeneratorState(shared_ptr<const SudokuBoardGeneratorState> source)
    : SudokuBoardGeneratorState(*source.get())
{
    // must clone the board
    _board = cloneToEdit(_board);
    _lastResult.reset();
    // do not reset _proven - it is reset if board is modified.
    // remember the source.
    _source = source;

    checkIntegrity();
}

void SudokuBoardGeneratorState::checkIntegrity() const
{
    if (!getIntegrityChecks())
    {
        return;
    }

    vector<int> all;

    for (int index : _indexes.remained())
    {
        if (_board->getValue(index) == SudokuValue::EMPTY)
        {
            throw logic_error("remained index points to an empty value");
        }
        all.push_back(index);
    }

    for (int index : _indexes.reserved())
    {
        // reserved values can point to either cell
        all.push_back(index);
    }

    for (int index : _indexes.removed())
    {
        // removed values can point to either cell
        all.push_back(index);
    }

    if (all.size() != SudokuBoard::BOARD_SIZE)
    {
        throw logic_error("total size of remained/reserved/removed is not same as board size");
    }

    sort(all.begin(), all.end());
    for (int i = 0; i < all.size(); i++)
    {
        if (i != all.at(i))
        {
            throw logic_error("one or more index values is wrong");
        }
    }
}

SudokuResultConstShared SudokuBoardGeneratorState::solve() const
{
    if (_lastResult)
    {
        throw logic_error("Do not solve twice on the same state, clone instead!");
    }
    if (!_proven)
    {
        throw logic_error("Should NEVER solve before prove!");
    }

    _lastResult = solveInternal();
    return _lastResult;
}

SudokuResultConstShared SudokuBoardGeneratorState::solveInternal() const
{
    SudokuSolverOptions options(SudokuSolverAction::SOLVE);
    auto result = _solver->run(options, _board);
    if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
    {
        cerr << "TIMEOUT detected during solve on previously proven board!" << endl;
        cerr << serializeBoard(result->getOriginalBoard().get()) << endl;
        cerr << result->getOriginalBoard() << endl;
    }
    else if (result->getStatus() != SudokuSolverStatus::SUCCEEDED)
    {
        cerr << serializeBoard(board().get());
        throw logic_error("failed to solve previosly proven board");
    }

    checkIntegrity();
    return result;
}

bool SudokuBoardGeneratorState::tryProve()
{
    if (_lastResult)
    {
        throw logic_error("Do not solve and then prove on the same state, clone instead!");
    }
    if (_proven)
    {
        throw logic_error("Do not prove twice on the same state (even after clone)!");
    }

    SudokuSolverOptions options(SudokuSolverAction::PROVE);
    auto result = _solver->run(options, _board, solution());
    if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
    {
        cerr << "TIMEOUT detected!" << endl;
        cerr << serializeBoard(result->getOriginalBoard().get()) << endl;
        cerr << result->getOriginalBoard() << endl;
    }
    checkIntegrity();
    _proven = result->getStatus() == SudokuSolverStatus::SUCCEEDED;
    return _proven;
}

bool SudokuBoardGeneratorState::shouldContinue(Random &r)
{
    if (level() >= SudokuLevel::BLACKHOLE)
    {
        // for BLACKHOLE: do not bother checking the level - let it go till the end!
        return true;
    }

    if (board()->getFreeCellCount() > 32)
    {
        // Too early even for simple games.
        return true;
    }

    SudokuResultConstShared result = solve();
    if (result->getStatus() != SudokuSolverStatus::SUCCEEDED)
    {
        return false;
    }

    if (result->getLevel() < level())
    {
        // Keep going.
        return true;
    }

    if (result->getLevel() > level())
    {
        // Overflow, stop.
        return false;
    }

    return shouldContinueByLevel(level(), r);
}

void SudokuBoardGeneratorState::restoreSimpleValues(Random &r, int minComplexity)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }
    if (!_proven)
    {
        throw logic_error("Should NEVER restoreSimpleValues before prove!");
    }

    indexes().shuffleRemoved(r);

    SudokuResultConstShared lastResult = solveInternal();

    for (int index : indexes().removed())
    {
        if (!_board->isEmptyCell(index))
        {
            // ignore cells with value, their indexes were marked as 'removed' because removing them leads to
            // unsolvable board.
            continue;
        }
        _board->setValue(index, solution()->getValue(index));
        SudokuResultConstShared result = solveInternal();
        if (result->getComplexity() < minComplexity ||
            (lastResult->getComplexity() - result->getComplexity()) > 5)
        {
            // Not a good choice for restoration - remove it
            _board->setValue(index, SudokuValue::EMPTY);
        }
        else
        {
            indexes().restoreRemoved(index);
            lastResult = result;
        }
    }
    checkIntegrity();
    _lastResult = lastResult;
}

bool SudokuBoardGeneratorState::removeWithRetries(
    Random &r, int freeAtLeast, int minToRemove, int maxToRemove, int maxRetries)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }

    bool keepGoing = true;

    // Remove fast, with up to 10 retries.
    while (keepGoing && board()->getFreeCellCount() < freeAtLeast)
    {
        keepGoing = tryRemove(r, minToRemove, maxToRemove, maxRetries);
        keepGoing = keepGoing && shouldContinue(r);
    }

    // Done!
    checkIntegrity();
    return keepGoing;
}

void SudokuBoardGeneratorState::removeOneByOne(Random &r)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }

    bool keepGoing = true;

    // Remove the last pieces one by one.
    while (keepGoing && board()->getFreeCellCount() < SudokuBoard::MAX_FREE_CELLS_FOR_VALID_BOARD)
    {
        keepGoing = tryRemoveOne(r);
        keepGoing = keepGoing && shouldContinue(r);
    }
    checkIntegrity();
}

bool SudokuBoardGeneratorState::tryRemove(Random &r, int minToRemove, int maxToRemove, int maxRetries)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }

    if (minToRemove > maxToRemove || minToRemove < 1)
    {
        throw logic_error("minToRemove and maxToRemove are out of range");
    }

    int retries = 0;
    while (retries < maxRetries && indexes().remainedSize() > 0)
    {
        int allowedToRemove = SudokuBoard::MAX_FREE_CELLS_FOR_VALID_BOARD - board()->getFreeCellCount();
        if (allowedToRemove == 0)
        {
            checkIntegrity();
            return false;
        }
        indexes().shuffleRemained(r);
        int currentBatch = r.nextInClosedRange(minToRemove, maxToRemove);
        currentBatch = min(currentBatch, indexes().remainedSize());
        currentBatch = min(currentBatch, allowedToRemove);
        if (tryRemoveBatch(r, currentBatch))
        {
            checkIntegrity();
            return true;
        }

        ++retries;
    }

    checkIntegrity();
    return false;
}

bool SudokuBoardGeneratorState::tryRemoveOne(Random &r)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }

    while (indexes().remainedSize() > 0)
    {
        if (tryRemoveBatch(r, 1))
        {
            checkIntegrity();
            return true;
        }
        // if batch has only one index and it fails, it is removed from remained list.
    }

    checkIntegrity();
    return false;
}

bool SudokuBoardGeneratorState::tryRemoveBatch(Random &r, int batchSize)
{
    if (_lastResult)
    {
        throw logic_error("Do not modify state after solve, clone instead!");
    }

    if (getIntegrityChecks())
    {
        // just make sure board is still valid (this is an expensive call...).
        if (!tryProve())
        {
            throw logic_error("do not use invalid boards as an input here");
        }
    }

    if (indexes().remainedSize() < batchSize)
    {
        throw logic_error("not enough indexes to remove");
    }

    indexes().reserve(batchSize);
    _proven = false;
    for (int index : indexes().reserved())
    {
        _board->setValue(index, SudokuValue::EMPTY);
    }

    bool success = tryProve();
    if (success)
    {
        indexes().removeReserved();
    }
    else
    {
        for (int index : indexes().reserved())
        {
            _board->setValue(index, solution()->getValue(index));
        }

        if (batchSize == 1)
        {
            // We removed only one and board failed - no point in trying same index again.
            indexes().removeReserved();
        }
        else
        {
            indexes().revertReserved();
        }
    }
    checkIntegrity();
    return success;
}

void SudokuBoardGeneratorState::mergeWith(SudokuBoardConstShared other)
{
    for (int index = 0; index < SudokuBoard::BOARD_SIZE; index++)
    {
        SudokuValue currentValue = _board->getValue(index);
        if (currentValue == SudokuValue::EMPTY)
        {
            continue;
        }

        SudokuValue otherValue = other->getValue(index);
        if (otherValue != SudokuValue::EMPTY && otherValue != currentValue)
        {
            throw logic_error("Values of current board do not match the values of merged one.");
        }

        if (otherValue == SudokuValue::EMPTY)
        {
            if (_indexes.tryRemoveIndex(index))
            {
                _board->setValue(index, SudokuValue::EMPTY);
            }
        }
    }
    _proven = false;
}

SudokuBoardGeneratorStateShared merge(
    SudokuBoardGeneratorStateConstShared state1, SudokuBoardGeneratorStateConstShared state2)
{
    SudokuBoardGeneratorStateShared merge = make_shared<SudokuBoardGeneratorState>(state1);

    merge->mergeWith(state2->board());
    return merge;
}
