#include "SudokuBoardImpl.h"

using namespace std;

bool SudokuBoardImpl::enableIntegrityChecks = false;

SudokuBoardImpl::SudokuBoardImpl(const SudokuBoard &other)
{
    const SudokuBoardImpl *board = dynamic_cast<const SudokuBoardImpl *>(&other);
    if (board != nullptr)
    {
        _values = board->_values;
        _readOnlyFlags = board->_readOnlyFlags;
    }
    else
    {
        for (int i = 0; i < BOARD_SIZE; i++)
        {
            setValuePrivate(i, other.getValue(i), other.isReadOnlyCell(i));
        }
    }
}

bool SudokuBoardImpl::isEquivalentTo(std::shared_ptr<const SudokuBoard> board) const
{
    for (int i = 0; i < BOARD_SIZE; i++)
    {
        if (board->getValue(i) != getValue(i))
        {
            return false;
        }
    }

    return true;
}

bool SudokuBoardImpl::containsAllValues(SudokuBoardConstShared board) const
{
    for (int i = 0; i < BOARD_SIZE; i++)
    {
        if (!board->isEmptyCell(i) && (board->getValue(i) != getValue(i)))
        {
            return false;
        }
    }

    return true;
}

bool SudokuBoardImpl::containsReadOnlyValues(SudokuBoardConstShared board) const
{
    for (int i = 0; i < BOARD_SIZE; i++)
    {
        if (board->isReadOnlyCell(i) && (board->getValue(i) != getValue(i)))
        {
            return false;
        }
    }

    return true;
}

void SudokuBoardImpl::print(ostream &out, const char *format) const
{
    if (format == nullptr || format[0] == '\0')
    {
        format = "v";
    }

    for (int i = 0; format[i]; i++)
    {
        switch (tolower(format[i]))
        {
        case 'v':
            appendValues(out);
            break;
        case 'r':
            appendRowValues(out);
            break;
        case 'c':
            appendColumnValues(out);
            break;
        case 's':
            appendSquareValues(out);
            break;
        case 't':
            out << "Serialized: " << serializeBoard(this) << endl;
            break;
        default:
            throw logic_error("unsupported format character");
        }
    }
}

SudokuValue SudokuBoardImpl::setValuePrivate(int index, SudokuValue value, bool isReadOnly)
{
    if (getAccessMode() < AccessMode::PLAY)
    {
        throw logic_error("This game cannot be modified");
    }

    if ((getAccessMode() < AccessMode::EDIT) && (isReadOnly || isReadOnlyCell(index)))
    {
        throw logic_error("Edit mode is not allowed");
    }

    if ((value == SudokuValue::EMPTY) && isReadOnly)
    {
        throw logic_error("Empty cell cannot be read-only");
    }

    SudokuValue previousValue = getValue(index);
    _values[index] = value;
    _readOnlyFlags.set(index, isReadOnly);
    return previousValue;
}

void SudokuBoardImpl::appendValues(ostream &out) const
{
    out << endl;
    out << "╔═══════╦═══════╦═══════╗";
    out << endl;

    for (int row = 0; row < SEQUENCE_SIZE; row++)
    {
        if (row != 0 && (row % 3) == 0)
        {
            out << "╠═══════╬═══════╬═══════╣" << endl;
        }
        for (int col = 0; col < SEQUENCE_SIZE; col++)
        {
            if (col % 3 == 0)
            {
                out << "║ ";
            }
            int i = indexFromCoordinates(row, col);
            char digit = '0' + static_cast<int>(getValue(i));
            out << digit;

            if (isReadOnlyCell(i))
            {
                if (isValidCell(i))
                {
                    out << ' ';
                }
                else
                {
                    out << 'X';
                }
            }
            else
            {
                if (isValidCell(i))
                {
                    out << '.';
                }
                else
                {
                    out << '!';
                }
            }
        }

        out << "║";
        out << endl;
    }

    out << "╚═══════╩═══════╩═══════╝";
    out << endl;
}

void SudokuBoardImpl::appendRowValues(ostream &out) const
{
    out << "Rows:";
    for (int i = 0; i < SEQUENCE_SIZE; i++)
    {
        out << ' ' << getRowValues(i);
    }
    out << endl;
}

void SudokuBoardImpl::appendColumnValues(ostream &out) const
{
    out << "Columns:";
    for (int i = 0; i < SEQUENCE_SIZE; i++)
    {
        out << ' ' << getColumnValues(i);
    }
    out << endl;
}

void SudokuBoardImpl::appendSquareValues(ostream &out) const
{
    out << "Squares:";
    for (int i = 0; i < SEQUENCE_SIZE; i++)
    {
        out << ' ' << getSquareValues(i);
    }
    out << endl;
}
