#pragma once

#include "SudokuBoard.h"
#include "Random.h"

class SudokuSolutionGenerator
{
  public:
    SudokuSolutionGenerator() = default;
    SudokuSolutionGenerator(const SudokuSolutionGenerator &) = delete;
    SudokuSolutionGenerator &operator=(const SudokuSolutionGenerator &) = delete;

    virtual SudokuSolutionConstShared generate(Random &r) = 0;

    SudokuSolutionConstShared generate()
    {
        Random r;
        return generate(r);
    }
};

using SudokuSolutionGeneratorShared = std::shared_ptr<SudokuSolutionGenerator>;

SudokuSolutionGeneratorShared newSolutionGenerator();
