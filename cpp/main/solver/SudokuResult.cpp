#include "SudokuResult.h"

using namespace std;

void SudokuResult::merge(const SudokuResult &result)
{
  _complexity += result._complexity;

  for (int sci = 0; sci < result._steps.size(); sci++)
  {
    const auto &stepCounters = result._steps.at(sci);
    auto &target = _steps[sci];
    for (int si = 0; si < stepCounters.size(); si++)
    {
      target[si] += stepCounters[si];
    }
  }
}

bool SudokuResult::addSolution(SudokuSolutionConstShared solution)
{
  for (auto current : _solutions)
  {
    if (current->isEquivalentTo(solution))
    {
      return false;
    }
  }
  _solutions.push_back(solution);
  return true;
}

ostream &operator<<(ostream &os, const SudokuResult &result)
{
  os << "Sudoku Result {" << endl;
  os << "  action: " << result.getOptions().getAction() << endl;
  os << "  status: " << result.getStatus() << endl;
  if (result.getOptions().getAction() == SudokuSolverAction::SOLVE)
  {
    os << "  level (complexity): " << result.getLevel() << " (" << result.getComplexity() << ")" << endl;
  }
  os << "  board:  " << serializeBoard(result.getOriginalBoard().get()) << endl;
  double elapsed_microseconds = chrono::duration_cast<std::chrono::microseconds>(result.elapsed()).count();
  os << "  elapsed (microseconds): " << elapsed_microseconds << endl;
  os << "  value count: " << (81 - result.getOriginalBoard()->getFreeCellCount()) << endl;

  const auto &steps = result.getSteps();
  if (steps.size() > 0)
  {
    os << "  steps: {" << endl;
    for (int sci = 0; sci < steps.size(); sci++)
    {
      const auto &stepCounters = steps[sci];
      SudokuStepComplexity stepComplexity = stepComplexityFromIndex(sci);
      for (int si = 0; si < stepCounters.size(); si++)
      {
        auto count = stepCounters[si];
        if (count == 0)
        {
          continue;
        }
        auto step = sudokuStepFromIndex(si);
        auto weight = stepComplexity * 1;
        os << "    " << stepComplexity << " [" << step << ", " << weight << "] X " << count
           << ": " << (weight * count) << endl;
      }
    }
    os << "  }" << endl;
  }
  else
  {
    os << "No steps available." << endl;
  }
  os << "}" << endl;

  return os;
}
