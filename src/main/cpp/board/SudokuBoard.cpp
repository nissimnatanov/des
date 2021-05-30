#include <sstream>

#include "SudokuSolutionImpl.h"
#include "SudokuEditBoardImpl.h"

using namespace std;

bool getIntegrityChecks()
{
    return SudokuBoardImpl::enableIntegrityChecks;
}

bool setIntegrityChecks(bool enable)
{
    bool prev = SudokuBoardImpl::enableIntegrityChecks;
    SudokuBoardImpl::enableIntegrityChecks = enable;
    return prev;
}

const SudokuBoardConstShared cloneAsImmutable(SudokuBoardConstShared other)
{
    if (other->getAccessMode() == SudokuBoard::AccessMode::IMMUTABLE)
    {
        return other;
    }

    return make_shared<SudokuOptimizedBoardImpl>(*other);
}

SudokuPlayerBoardShared cloneToPlay(SudokuBoardConstShared other)
{
    return make_shared<SudokuPlayerBoardImpl>(*other);
}

SudokuEditBoardShared cloneToEdit(SudokuBoardConstShared other)
{
    return make_shared<SudokuEditBoardImpl>(*other);
}

const SudokuSolutionConstShared cloneAsSolution(SudokuBoardConstShared other)
{
    const SudokuSolution *solution = dynamic_cast<const SudokuSolution *>(other.get());
    if (solution != nullptr)
    {
        return shared_ptr<const SudokuSolution>(other, solution);
    }
    else
    {
        return make_shared<const SudokuSolutionImpl>(*other);
    }
}

string serializeBoard(const SudokuBoard *board)
{
    ostringstream s;
    int spacesSeen = 0;
    for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
    {
        SudokuValue v = board->getValue(i);
        if (v == SudokuValue::EMPTY)
        {
            ++spacesSeen;
            if (spacesSeen > 26)
            {
                spacesSeen -= 26;
                s << 'Z';
            }
        }
        else
        {
            if (spacesSeen > 0)
            {
                char spaceSub = 'A' + (spacesSeen - 1);
                spacesSeen = 0;
                s << spaceSub;
            }
            char digit = '0' + static_cast<int>(v);
            s << digit;
            if (!board->isReadOnlyCell(i))
            {
                s << '_'; // provided by the player
            }
        }
    }

    if (spacesSeen > 0)
    {
        char spaceSub = 'A' + (spacesSeen - 1);
        spacesSeen = 0;
        s << spaceSub;
    }
    return s.str();
}

SudokuEditBoardShared deserializeSudokuBoard(std::string input)
{
    SudokuEditBoardShared board = newEditBoard();
    int i = 0;
    for (char c : input)
    {
        if (c == ' ' || c == '\t' || c == '\r' || c == '\n')
        {
            continue;
        }
        if (c >= 'a' && c <= 'z')
        {
            i += 1 + c - 'a';
        }
        else if (c >= 'A' && c <= 'Z')
        {
            i += 1 + c - 'A';
        }
        else if (c == '0')
        {
            ++i;
        }
        else if (c >= '1' && c <= '9')
        {
            SudokuValue v = static_cast<SudokuValue>(c - '0');
            board->setValue(i, v);
            ++i;
        }
        else if (c == '_')
        {
            SudokuValue v = board->getValue(i);
            if (v == SudokuValue::EMPTY)
            {
                throw logic_error("misplaced _ sign");
            }
            else if (!board->isReadOnlyCell(i))
            {
                throw logic_error("duplicate _ sign");
            }
            else
            {
                board->setValue(i, SudokuValue::EMPTY);
                board->playValue(i, v);
            }
        }
        else
        {
            throw logic_error(string("Invalid board character '") + c + "', at index: " + to_string(i));
        }
    }

    if (i != SudokuBoard::BOARD_SIZE)
    {
        throw logic_error(string("Incomplete board, stopped at index: ") + to_string(i));
    }

    return board;
}