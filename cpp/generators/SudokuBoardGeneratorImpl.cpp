#include <algorithm>

#include "SudokuBoardGeneratorImpl.h"

SudokuBoardGeneratorShared newBoardGenerator()
{
    return make_shared<SudokuBoardGeneratorImpl>();
}

SudokuResultConstShared SudokuBoardGeneratorImpl::generate(
    SudokuLevel level, Random &r, SudokuSolutionOptional solutionOrNull)
{
    SudokuSolverShared solver = createSolver();

    if (level >= SudokuLevel::VERYHARD)
    {
        auto state = make_shared<SudokuBoardGeneratorState>(level, r, solutionOrNull, solver);
        return generateSlow(r, state);
    }
    else
    {
        SudokuBoardGeneratorState state(level, r, solutionOrNull, solver);
        return generate(r, state);
    }
}

SudokuResultConstShared SudokuBoardGeneratorImpl::generate(Random &r, SudokuBoardGeneratorState &state)
{
    bool keepGoing =
        state.removeWithRetries(
            r, /* freeAtLeast= */ 32, /* minToRemove= */ 3, /* maxToRemove= */ 8, /* maxRetries= */ 10) &&
        state.removeWithRetries(
            r, /* freeAtLeast= */ 48, /* minToRemove= */ 1, /* maxToRemove= */ 3, /* maxRetries= */ 15);

    if (keepGoing)
    {
        state.removeOneByOne(r);
    }

    // Done!
    return state.getLastResultOrSolve();
}

void sortByComplexityAndTrim(vector<SudokuBoardGeneratorStateConstShared> &states, size_t trimTo, int minComplexity)
{
    for (SudokuBoardGeneratorStateConstShared state : states)
    {
        // Ensure all have cached result.
        state->getLastResultOrSolve();
    }
    sort(states.begin(), states.end(),
         [](SudokuBoardGeneratorStateConstShared state1, SudokuBoardGeneratorStateConstShared state2) -> bool {
             SudokuResultConstShared r1 = state1->getLastResultOrSolve();
             SudokuResultConstShared r2 = state2->getLastResultOrSolve();
             return r1->getComplexity() > r2->getComplexity();
         });
    // get rid of duplicates
    states.erase(
        unique(states.begin(), states.end(), [](SudokuBoardGeneratorStateConstShared state1, SudokuBoardGeneratorStateConstShared state2) -> bool {
            SudokuResultConstShared r1 = state1->getLastResultOrSolve();
            SudokuResultConstShared r2 = state2->getLastResultOrSolve();
            return r1->getComplexity() == r2->getComplexity() &&
                   r1->getOriginalBoard()->isEquivalentTo(r2->getOriginalBoard());
        }),
        states.end());
    trimTo = min(trimTo, states.size());
    while (trimTo > 1 &&
           minComplexity > states[trimTo - 1]->getLastResultOrSolve()->getComplexity())
    {
        --trimTo;
    }
    if (states.size() > trimTo)
    {
        states.erase(states.begin() + trimTo, states.end());
    }
}

void getBestWithRetries(
    Random &r,
    vector<SudokuBoardGeneratorStateConstShared> &best,
    SudokuBoardGeneratorStateConstShared start,
    int count, int freeAtLeast,
    int minToRemove, int maxToRemove, int maxRetries)
{
    int newAdded = 0;
    bool betterFound = false;
    bool success = false;
    do
    {
        SudokuBoardGeneratorStateShared state = make_shared<SudokuBoardGeneratorState>(start);
        success = state->removeWithRetries(r, freeAtLeast, minToRemove, maxToRemove, maxRetries);
        if (success)
        {
            newAdded++;
            best.push_back(state);
            betterFound |= state->getLastResultOrSolve()->getComplexity() > start->getLastResultOrSolve()->getComplexity();
        }
    } while (success && newAdded < count);

    if (!betterFound)
    {
        best.push_back(make_shared<SudokuBoardGeneratorState>(start));
    }
}

SudokuBoardGeneratorStateConstShared SudokuBoardGeneratorImpl::removePhase0(
    Random &r, SudokuBoardGeneratorStateConstShared initialState, int freeCount)
{
    do
    {
        SudokuBoardGeneratorStateShared state = make_shared<SudokuBoardGeneratorState>(initialState);
        // Tests show that the fastest and the most efficient way to remove the first batch appears to be
        // a [2/3] combination.
        bool success = state->removeWithRetries(
            r, freeCount - 2, /* minToRemove= */ 2, /* maxToRemove= */ 3, /* maxRetries= */ 15);
        if (success)
        {
            success = state->removeWithRetries(
                r, freeCount, /* minToRemove= */ 1, /* maxToRemove= */ 1, /* maxRetries= */ 15);
        }
        if (success)
        {
            state->getLastResultOrSolve();
            return state;
        }
    } while (true);
}

static const bool debug_print = false;

static const SudokuBoardGeneratorImpl::Phase phases[] = {
    // freeCount, minComplexity, enforceMinComplexity, generateCount, selectCount, restore, merge
    {40, 100, true, 1000, 5, true, true},
    {42, 220, true, 20, 5, true, false},
    {44, 450, true, 30, 5, true, false},
    {46, 700, true, 40, 5, true, false},
    {48, 1000, true, 50, 5, true, false},
    {50, 1300, true, 60, 5, true, false},
    {53, 1700, true, 70, 5, true, false},
    {53, 1700, true, 80, 5, true, false},
    {56, 2000, true, 90, 5, true, false},
    {56, 2000, true, 100, 5, true, false},
    // Stop restoration
    {56, 2000, true, 100, 5, false, false},
};

static const int phase_count = sizeof(phases) / sizeof(phases[0]);

vector<SudokuBoardGeneratorStateConstShared> SudokuBoardGeneratorImpl::generatePhase0Batch(
    Random &r,
    SudokuBoardGeneratorStateConstShared initialState,
    const Phase &phase)
{
    vector<SudokuBoardGeneratorStateConstShared> best;
    bool done = false;
    int count = 0;
    do
    {
        SudokuBoardGeneratorStateConstShared state = removePhase0(r, initialState, phase.freeCount);

        best.push_back(state);
        count++;
        if (count < phase.generateCount)
        {
            done = false;
        }
        else
        {
            sortByComplexityAndTrim(best, /* trimTo= */ phase.selectCount, phase.minComplexity);
            if (best.size() == 0)
            {
                done = false;
            }
            else
            {
                done = !phase.enforceMinComplexity ||
                       phase.minComplexity <= best[0]->getLastResultOrSolve()->getComplexity();
            }
        }
    } while (!done);
    return best;
}

vector<SudokuBoardGeneratorStateConstShared> SudokuBoardGeneratorImpl::generatePhaseNBatch(
    int phi,
    Random &r,
    const vector<SudokuBoardGeneratorStateConstShared> &sources,
    const Phase &phase)
{
    vector<SudokuBoardGeneratorStateConstShared> next;
    for (SudokuBoardGeneratorStateConstShared b : sources)
    {
        int maxRetries = min(30, b->indexes().remainedSize());
        int minToRemove = 1;
        int maxToRemove = 2;
        if (b->indexes().remainedSize() > 44)
        {
            minToRemove++;
            maxToRemove++;
        }
        getBestWithRetries(
            r, next, b,
            phase.generateCount, phase.freeCount,
            minToRemove, maxToRemove,
            maxRetries);
    }
    sortByComplexityAndTrim(next, /* trimTo= */ phase.selectCount, phase.minComplexity);
    return move(next);
}

void printBest(int iter, const char *state, int minComplexity, const vector<SudokuBoardGeneratorStateConstShared> &best)
{
    if (debug_print)
    {
        int complexityFirst = best.at(0)->getLastResultOrSolve()->getComplexity();
        int freeCellsFirst = best.at(0)->getLastResultOrSolve()->getOriginalBoard()->getFreeCellCount();
        cout << "Iter[" << iter << ", " << state << "]"
             << " minComplexity=" << minComplexity << ", size=" << best.size() << ", complexity=["
             << complexityFirst << " #" << freeCellsFirst;
        if (best.size() > 1)
        {
            int complexityLast = best.at(best.size() - 1)->getLastResultOrSolve()->getComplexity();
            int freeCellsLast = best.at(best.size() - 1)->getLastResultOrSolve()->getOriginalBoard()->getFreeCellCount();
            cout << ", " << complexityLast << " #" << freeCellsLast;
        }
        cout << "]" << endl;
    }
}

int getMaxCandidateSize(int depth)
{
    // 5 * 5 * 5 * 4 * 3 * 2 = 3,000
    if (depth < 4)
    {
        return 5;
    }
    if (depth < 7)
    {
        return 8 - depth;
    }

    return 1;
}

void findAllPossibleBoards(
    Random &r,
    SudokuBoardGeneratorStateConstShared source,
    vector<SudokuBoardGeneratorStateConstShared> &best,
    int depth)
{
    // First, remove all indexes that lead to unsolvable board.
    vector<int> mustStayIndexes;
    vector<pair<SudokuBoardGeneratorStateShared, int>> candidates;
    for (int index : source->indexes().remained())
    {
        SudokuBoardGeneratorStateShared state = make_shared<SudokuBoardGeneratorState>(source);
        state->indexes().prioritizeIndex(index);
        if (state->tryRemoveBatch(r, 1))
        {
            candidates.push_back(make_pair(state, index));
        }
        else
        {
            mustStayIndexes.push_back(index);
        }
    }

    if (candidates.empty())
    {
        // Leaf board - stop now.
        best.push_back(source);
        return;
    }

    int max_candidates = getMaxCandidateSize(depth);

    if (candidates.size() > max_candidates)
    {
        // recursion depth will blow, limit it.
        r.shuffle(candidates);
        candidates.erase(candidates.begin() + max_candidates, candidates.end());
    }

    // We have left with states with valid boards, try to fine tune one-by-one only those boards, removing
    // indexes previously found as "must stay".
    for (pair<SudokuBoardGeneratorStateShared, int> candidate : candidates)
    {
        for (int index : mustStayIndexes)
        {
            candidate.first->indexes().removeIndex(index);
        }

        // recursivly find more boards.
        findAllPossibleBoards(r, candidate.first, best, depth + 1);
        mustStayIndexes.push_back(candidate.second);
    }
    // once in a while, get rid of duplicates
    if (depth % 3 == 0)
    {
        sortByComplexityAndTrim(best, best.size(), 0);
    }
}

bool addMergedBoards(vector<SudokuBoardGeneratorStateConstShared> &best)
{
    bool createdMerged = false;
    // capture size before adding new ones.
    int originalSize = best.size();

    for (int i = 0; (i + 1) < originalSize; i++)
    {
        for (int j = i + 1; j < originalSize; j++)
        {
            SudokuBoardGeneratorStateShared mergeState = merge(best.at(i), best.at(j));
            if (mergeState->tryProve())
            {
                best.push_back(mergeState);
                createdMerged = true;
            }
        }
    }
    return createdMerged;
}

void restore(Random &r, vector<SudokuBoardGeneratorStateConstShared> &best, int minComplexity)
{
    for (int i = 0; i < best.size(); i++)
    {
        SudokuBoardGeneratorStateShared state = make_shared<SudokuBoardGeneratorState>(best.at(i));
        state->restoreSimpleValues(r, minComplexity);
        best.at(i) = state;
    }
}

SudokuResultConstShared SudokuBoardGeneratorImpl::generateSlow(
    Random &r, SudokuBoardGeneratorStateConstShared initialState)
{
    vector<SudokuBoardGeneratorStateConstShared> current;

    for (int phi = 0; phi < phase_count; phi++)
    {
        const Phase &phase = phases[phi];
        if (phi == 0)
        {
            current = generatePhase0Batch(r, initialState, phase);
        }
        else
        {
            current = generatePhaseNBatch(phi, r, current, phase);
            if (phase.enforceMinComplexity &&
                (phase.minComplexity > current[0]->getLastResultOrSolve()->getComplexity()))
            {
                if (debug_print)
                {
                    printBest(phi, "failed", phase.minComplexity, current);
                    cout << "----" << endl;
                }
                // do not bother if requirements were not met.
                return current[0]->getLastResultOrSolve();
            }
        }
        printBest(phi, "base", phase.minComplexity, current);

        if (phase.restore)
        {
            restore(r, current, phase.minComplexity);
            printBest(phi, "restored", phase.minComplexity, current);
        }
        if (phase.merge)
        {
            bool mergedAtLeastOne = addMergedBoards(current);
            printBest(phi, "merged(before sort)", phase.minComplexity, current);
            if (mergedAtLeastOne)
            {
                sortByComplexityAndTrim(current, /* trimTo= */ phase.selectCount, phase.minComplexity);
                printBest(phi, "merged", phase.minComplexity, current);
                if (phase.restore)
                {
                    restore(r, current, phase.minComplexity);
                    printBest(phi, "merged/restored", phase.minComplexity, current);
                }
            }
        }
    }

    {
        vector<SudokuBoardGeneratorStateConstShared> last;
        for (SudokuBoardGeneratorStateConstShared bLast : current)
        {
            findAllPossibleBoards(r, bLast, last, 1);
        }
        sortByComplexityAndTrim(last, /* trimTo= */ 3, /* minComplexity= */ 0);
        current = move(last);
    }
    printBest(phase_count, "last", 0, current);

    for (SudokuBoardGeneratorStateConstShared b : current)
    {
        auto result = b->getLastResultOrSolve();
        if (result->getLevel() < SudokuLevel::DARKEVIL)
        {
            break;
        }

        string board = serializeBoard(result->getOriginalBoard().get());
        cout << result->getLevel() << " (" << result->getComplexity() << "): " << board << endl;

        // print generation chain
        shared_ptr<const SudokuBoardGeneratorState> bNext = b->source();
        cout << "         " << result->getComplexity();
        while (bNext)
        {
            if (bNext->board()->isSolved())
            {
                // OK for solution to be in the chain.
                break;
            }
            cout << "," << bNext->getLastResultOrSolve()->getComplexity();
            bNext = bNext->source();
        }
        cout << endl;
    }
    if (debug_print)
    {
        cout << "----" << endl;
    }
    return current.at(0)->getLastResultOrSolve();
}
