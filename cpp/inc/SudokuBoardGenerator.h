#pragma once

#include "SudokuSolutionGenerator.h"
#include "SudokuResult.h"
#include "SudokuLevel.h"

class SudokuBoardGenerator
{
public:
  SudokuBoardGenerator() = default;
  SudokuBoardGenerator(const SudokuBoardGenerator &) = delete;
  SudokuBoardGenerator &operator=(const SudokuBoardGenerator &) = delete;

  virtual SudokuResultConstShared generate(
      SudokuLevel level, Random &r, SudokuSolutionOptional solutionOrNull) = 0;

  SudokuResultConstShared generate(SudokuLevel level)
  {
    Random r;
    return generate(level, r, SudokuSolutionOptional::empty());
  }
};

using SudokuBoardGeneratorShared = std::shared_ptr<SudokuBoardGenerator>;

SudokuBoardGeneratorShared newBoardGenerator();
