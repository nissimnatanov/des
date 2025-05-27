#include "SudokuBoardImpl.h"

using namespace std;

class SudokuOptimizedBoardImpl : public SudokuBoardImpl
{
private:
  array<SudokuValueSet, SEQUENCE_SIZE> _rowValues;
  array<SudokuValueSet, SEQUENCE_SIZE> _columnValues;
  array<SudokuValueSet, SEQUENCE_SIZE> _squareValues;
  array<int, SEQUENCE_SIZE + 1> _valueCounters;
  array<bool, BOARD_SIZE> _validFlags;
  array<SudokuValueSet, BOARD_SIZE> _userDisallowedValues;
  bool _valid;
  int _suppressRecalculateCount;
  bool _needToRecalculate;

  // Allowed values cache is special - the moment at least one value changes, it is easier to invalidate
  // all indexes instead of recalculating related ones.
  mutable bool _allowedValuesCacheValid;
  mutable array<SudokuValueSet, BOARD_SIZE> _allowedValuesCache;
  mutable bitset<BOARD_SIZE> _allowedValuesCacheValidFlags;

  void initializeFrom(const SudokuOptimizedBoardImpl *otherOrNull);

  void initializeEmpty()
  {
    initializeFrom(nullptr);
  }

  bool updateStats(int index, SudokuValue oldValue, SudokuValue newValue);

  void recalculateAllStats();
  void recalculateSequence(SudokuValueSet &sequenceSet, const sequence &sequenceIndexes);
  void markSequenceInvalid(SudokuValue v, const sequence &sequenceIndexes);
  SudokuValueSet getRelativeValues(int index) const;

  void checkIntegrity();

protected:
  void suppressRecalculations()
  {
    _suppressRecalculateCount++;
  }

  void resumeRecalculations()
  {
    --_suppressRecalculateCount;
    if (_suppressRecalculateCount == 0 && _needToRecalculate)
    {
      _needToRecalculate = false;
      recalculateAllStats();
    }
  }

  void addUserDisallowedValues(int index, SudokuValueSet values)
  {
    _userDisallowedValues.at(index) += values;
    _allowedValuesCacheValidFlags.reset(index);
  }

  void resetUserDisallowedValues(int index)
  {
    _userDisallowedValues.at(index) = SudokuValueSet::none();
    _allowedValuesCacheValidFlags.reset(index);
  }

public:
  SudokuOptimizedBoardImpl()
  {
    initializeEmpty();
    checkIntegrity();
  }

  SudokuOptimizedBoardImpl(const SudokuBoard &other);

  ~SudokuOptimizedBoardImpl() {}

  AccessMode getAccessMode() const override
  {
    return AccessMode::IMMUTABLE;
  }

  SudokuValueSet getRowValues(int row) const override
  {
    return _rowValues.at(row);
  }

  SudokuValueSet getColumnValues(int column) const override
  {
    return _columnValues.at(column);
  }

  SudokuValueSet getSquareValues(int square) const override
  {
    return _squareValues.at(square);
  }

  SudokuValueSet getAllowedValues(int index) const override;

  virtual int getValueCount(SudokuValue v) const override
  {
    return _valueCounters.at(static_cast<int>(v));
  }

  bool isValidCell(int index) const override
  {
    return _validFlags.at(index);
  }

  bool isValid() const override
  {
    return _valid;
  }

  SudokuValue setValueInternal(int index, SudokuValue value, bool isReadOnly) override;
};
