#include "SudokuSolutionGenerator.h"

class SudokuSolutionGeneratorImpl : public SudokuSolutionGenerator
{
public:
    SudokuSolutionConstShared generate(Random& r) override;
};
