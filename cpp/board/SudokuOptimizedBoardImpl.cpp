#include <set>
#include "SudokuOptimizedBoardImpl.h"

using namespace std;

SudokuOptimizedBoardImpl::SudokuOptimizedBoardImpl(const SudokuBoard &other) : SudokuBoardImpl(other)
{
    const SudokuOptimizedBoardImpl *board = dynamic_cast<const SudokuOptimizedBoardImpl *>(&other);
    if (board != nullptr)
    {
        initializeFrom(board);
    }
    else
    {
        _suppressRecalculateCount = 0;
        _needToRecalculate = false;
        recalculateAllStats();
    }

    checkIntegrity();
}

void SudokuOptimizedBoardImpl::initializeFrom(const SudokuOptimizedBoardImpl *otherOrNull)
{
    if (otherOrNull != nullptr)
    {
        if (otherOrNull->_suppressRecalculateCount != 0)
        {
            throw logic_error("Cannot copy board that has calculations suppressed.");
        }
        _rowValues = otherOrNull->_rowValues;
        _columnValues = otherOrNull->_columnValues;
        _squareValues = otherOrNull->_squareValues;
        _valueCounters = otherOrNull->_valueCounters;
        _validFlags = otherOrNull->_validFlags;
        _userDisallowedValues = otherOrNull->_userDisallowedValues;
        _valid = otherOrNull->_valid;
        _allowedValuesCacheValid = otherOrNull->_allowedValuesCacheValid;
        if (_allowedValuesCacheValid)
        {
            _allowedValuesCache = otherOrNull->_allowedValuesCache;
            _allowedValuesCacheValidFlags = otherOrNull->_allowedValuesCacheValidFlags;
        }
    }
    else
    {
        _rowValues.fill(SudokuValueSet::none());
        _columnValues.fill(SudokuValueSet::none());
        _squareValues.fill(SudokuValueSet::none());
        _userDisallowedValues.fill(SudokuValueSet::none());
        _valueCounters.fill(0);
        _valueCounters.at(static_cast<int>(SudokuValue::EMPTY)) = BOARD_SIZE;
        _validFlags.fill(true);
        _valid = true;
        _allowedValuesCacheValid = false;
    }
    _suppressRecalculateCount = 0;
    _needToRecalculate = false;
}

SudokuValue SudokuOptimizedBoardImpl::setValueInternal(int index, SudokuValue value, bool isReadOnly)
{
    SudokuValue oldValue = SudokuBoardImpl::setValueInternal(index, value, isReadOnly);
    bool needToRecalculateAll = updateStats(index, oldValue, value);
    if (needToRecalculateAll)
    {
        recalculateAllStats();
    }
    return oldValue;
}

bool SudokuOptimizedBoardImpl::updateStats(int index, SudokuValue oldValue, SudokuValue newValue)
{
    if (oldValue == newValue)
    {
        return false;
    }

    int row = rowFromIndex(index);
    int column = columnFromIndex(index);
    int square = squareFromIndex(index);

    // if board was valid before and the new value does not appear in relative cells,
    // there is no need to re-validate the board.
    if (newValue != SudokuValue::EMPTY)
    {
        _valid = _valid && !getRelativeValues(index).contains(newValue);
    }

    if (!_valid)
    {
        // recalculate all - do not care about performance for this case ...
        return true;
    }

    if (oldValue != SudokuValue::EMPTY)
    {
        _rowValues.at(row) -= oldValue;
        _columnValues.at(column) -= oldValue;
        _squareValues.at(square) -= oldValue;
    }
    --_valueCounters.at(static_cast<int>(oldValue));

    if (newValue != SudokuValue::EMPTY)
    {
        _rowValues.at(row) += newValue;
        _columnValues.at(column) += newValue;
        _squareValues.at(square) += newValue;
    }
    ++_valueCounters.at(static_cast<int>(newValue));
    _allowedValuesCacheValid = false;
    return false;
}

SudokuValueSet SudokuOptimizedBoardImpl::getAllowedValues(int index) const
{
    if (getValue(index) != SudokuValue::EMPTY)
    {
        return SudokuValueSet::none();
    }

    if (_allowedValuesCacheValid && _allowedValuesCacheValidFlags[index])
    {
        return _allowedValuesCache.at(index);
    }

    if (!_allowedValuesCacheValid)
    {
        _allowedValuesCacheValidFlags.reset();
        _allowedValuesCacheValid = true;
    }

    SudokuValueSet disallowedValues = getRelativeValues(index) + _userDisallowedValues.at(index);
    SudokuValueSet allowedValues = ~disallowedValues;
    _allowedValuesCacheValidFlags[index] = true;
    _allowedValuesCache[index] = allowedValues;
    return allowedValues;
}

SudokuValueSet SudokuOptimizedBoardImpl::getRelativeValues(int index) const
{
    return getRowValues(rowFromIndex(index)) +
           getColumnValues(columnFromIndex(index)) +
           getSquareValues(squareFromIndex(index));
}

void SudokuOptimizedBoardImpl::recalculateAllStats()
{
    if (_suppressRecalculateCount > 0)
    {
        // wait for resume
        _needToRecalculate = true;
        return;
    }

    // either board is invalid or just copied from another board
    initializeEmpty();

    for (int i = 0; i < BOARD_SIZE; i++)
    {
        SudokuValue v = getValue(i);
        if (v != SudokuValue::EMPTY)
        {
            --_valueCounters.at(static_cast<int>(SudokuValue::EMPTY));
            ++_valueCounters.at(static_cast<int>(v));
        }
    }

    // _rowValues; _columnValues; _squareValues; _valueCounters; _validFlags; _valid;
    for (int seq = 0; seq < SEQUENCE_SIZE; seq++)
    {
        recalculateSequence(_rowValues.at(seq), getRowIndexes(seq));
        recalculateSequence(_columnValues.at(seq), getColumnIndexes(seq));
        recalculateSequence(_squareValues.at(seq), getSquareIndexes(seq));
    }
}
void SudokuOptimizedBoardImpl::recalculateSequence(SudokuValueSet &sequenceSet, const sequence &sequenceIndexes)
{
    set<SudokuValue> dupValues;
    for (int i : sequenceIndexes)
    {
        SudokuValue v = getValue(i);
        if (v == SudokuValue::EMPTY)
        {
            continue;
        }
        if (sequenceSet.contains(v))
        {
            dupValues.insert(v);
        }
        else
        {
            sequenceSet += v;
        }
    }
    if (dupValues.size() > 0)
    {
        for (SudokuValue v : dupValues)
        {
            markSequenceInvalid(v, sequenceIndexes);
        }
    }
}

void SudokuOptimizedBoardImpl::markSequenceInvalid(SudokuValue v, const sequence &sequenceIndexes)
{
    _valid = false;
    vector<int> readOnly;
    bool foundReadWrite = false;

    for (int i : sequenceIndexes)
    {
        if (getValue(i) != v)
        {
            continue;
        }
        if (isReadOnlyCell(i))
        {
            readOnly.push_back(i);
        }
        else
        {
            foundReadWrite = true;
            _validFlags.at(i) = false;
        }
    }

    if (!foundReadWrite || readOnly.size() > 1)
    {
        for (int i : readOnly)
        {
            _validFlags.at(i) = false;
        }
    }
}

void SudokuOptimizedBoardImpl::checkIntegrity()
{
    if (!enableIntegrityChecks)
    {
        return;
    }

    array<int, SEQUENCE_SIZE + 1> valueCounters{};

    for (int i = 0; i < BOARD_SIZE; i++)
    {
        SudokuValue value = getValue(i);

        valueCounters.at(static_cast<int>(value)) += 1;

        // ensure value does not appear
        if (value != SudokuValue::EMPTY)
        {
            // check this value is disallowed in other places
            auto relatedIndicies = getRelated(i);
            for (int related : relatedIndicies)
            {
                if (getValue(related) == SudokuValue::EMPTY)
                {
                    assertIntegrity(!isAllowedValue(related, value));
                }
                else if (getValue(related) == value)
                {
                    // ensure one of them is marked as wrong
                    if (!isReadOnlyCell(related))
                    {
                        assertIntegrity(!isValidCell(related));
                    }
                    if (!isReadOnlyCell(i))
                    {
                        assertIntegrity(!isValidCell(i));
                    }
                    if (isReadOnlyCell(related) && isReadOnlyCell(i))
                    {
                        assertIntegrity(!isValidCell(related));
                        assertIntegrity(!isValidCell(i));
                    }
                }
            }
        }
        else
        {
            // check that disallowed values are a union of row/column/square
            int row = rowFromIndex(i);
            int col = columnFromIndex(i);
            int square = squareFromIndex(i);

            SudokuValueSet disallowedValuesExpected =
                getRowValues(row) +
                getColumnValues(col) +
                getSquareValues(square) +
                _userDisallowedValues.at(i);

            assertIntegrity(getDisallowedValues(i) == disallowedValuesExpected);
        }
    }

    for (int v = 0; v <= SEQUENCE_SIZE; v++)
    {
        assertIntegrity(getValueCount(static_cast<SudokuValue>(v)) == valueCounters.at(v));
    }
}
