#pragma once

#include <memory>
#include "SudokuBoardIndexManager.h"
#include "SudokuSolver.h"

using namespace std;

class SudokuBoardGeneratorState
{
private:
  SudokuLevel _level;
  SudokuSolutionConstShared _solution;
  SudokuEditBoardShared _board;
  SudokuBoardIndexManager _indexes;
  SudokuSolverShared _solver;
  bool _proven;
  mutable SudokuResultConstShared _lastResult;
  shared_ptr<const SudokuBoardGeneratorState> _source;

  void checkIntegrity() const;
  SudokuResultConstShared solveInternal() const;

  SudokuBoardGeneratorState(const SudokuBoardGeneratorState &source) = default;

public:
  SudokuBoardGeneratorState() = delete;
  SudokuBoardGeneratorState &operator=(const SudokuBoardGeneratorState &) = delete;
  SudokuBoardGeneratorState(SudokuLevel level, Random &r, SudokuSolverShared solver, SudokuBoardConstShared board);

  SudokuBoardGeneratorState(shared_ptr<const SudokuBoardGeneratorState> source);
  SudokuBoardGeneratorState(SudokuLevel level, Random &r, SudokuSolutionOptional solutionOrNull, SudokuSolverShared solver);

  SudokuLevel level() const noexcept
  {
    return _level;
  }

  SudokuSolutionConstShared solution() const noexcept
  {
    return _solution;
  }

  SudokuEditBoardConstShared board() const noexcept
  {
    return _board;
  }

  SudokuBoardIndexManager &indexes() noexcept
  {
    return _indexes;
  }

  const SudokuBoardIndexManager &indexes() const noexcept
  {
    return _indexes;
  }

  shared_ptr<const SudokuBoardGeneratorState> source() const noexcept
  {
    return _source;
  }

  SudokuResultConstShared solve() const;

  const SudokuResultConstShared getLastResultOrSolve() const
  {
    return _lastResult ? _lastResult : solve();
  }

  bool tryProve();

  bool removeWithRetries(Random &r, int freeAtLeast, int minToRemove, int maxToRemove, int maxRetries);
  void removeOneByOne(Random &r);
  void restoreSimpleValues(Random &r, int minComplexity);

  bool tryRemove(Random &r, int minToRemove, int maxToRemove, int maxRetries);
  bool tryRemoveOne(Random &r);
  bool tryRemoveBatch(Random &r, int batchSize);
  bool shouldContinue(Random &r);

  void mergeWith(SudokuBoardConstShared board);
};

using SudokuBoardGeneratorStateShared = shared_ptr<SudokuBoardGeneratorState>;
using SudokuBoardGeneratorStateConstShared = shared_ptr<const SudokuBoardGeneratorState>;

SudokuBoardGeneratorStateShared merge(
  SudokuBoardGeneratorStateConstShared state1, SudokuBoardGeneratorStateConstShared state2);
