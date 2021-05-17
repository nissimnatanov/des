package com.darkevilsudoku.board.impl;

import com.darkevilsudoku.board.*;

import java.util.FormatFlagsConversionMismatchException;

public class SudokuBoardImpl implements SudokuBoard {

	private SudokuBoardCell[] values;
	private SudokuValueSet[] columnValues;
	private SudokuValueSet[] rowValues;
	private SudokuValueSet[] squareValues;
	private boolean editMode;
	private boolean readOnly;
	private boolean valid;
	private SudokuValueSet[] userDisallowed;
	private int[] valueStats;

	public SudokuBoardImpl() {
		values = SudokuBoardCell.buildArray(BOARD_SIZE);

		editMode = true;
		valid = true;

		valueStats = new int[10]; // 0 is for empty cells, 1 - 9 are for digits
		valueStats[0] = BOARD_SIZE;

		rowValues = SudokuValueSet.buildArray(SudokuBoard.SEQUENCE_SIZE);
		columnValues = SudokuValueSet.buildArray(SudokuBoard.SEQUENCE_SIZE);
		squareValues = SudokuValueSet.buildArray(SudokuBoard.SEQUENCE_SIZE);

		userDisallowed = SudokuValueSet.buildArray(SudokuBoard.BOARD_SIZE);

		checkIntegrity();
	}

	@Override
	public SudokuBoard clone() {
		SudokuBoardImpl newBoard;
		try {
			newBoard = (SudokuBoardImpl) super.clone();
		} catch (CloneNotSupportedException e) {
			throw new RuntimeException(e);
		}

		newBoard.values = (SudokuBoardCell[]) this.values.clone();
		for (int i = 0; i < BOARD_SIZE; i++) {
			newBoard.values[i] = newBoard.values[i].clone();
		}

		newBoard.columnValues = (SudokuValueSet[]) this.columnValues.clone();
		newBoard.rowValues = (SudokuValueSet[]) this.rowValues.clone();
		newBoard.squareValues = (SudokuValueSet[]) this.squareValues.clone();

		newBoard.userDisallowed = (SudokuValueSet[]) this.userDisallowed.clone();

		newBoard.valueStats = (int[]) this.valueStats.clone();

		newBoard.checkIntegrity();
		return newBoard;
	}

	@Override
	public boolean hasValue(int index) {
		return values[index].hasValue();
	}

	@Override
	public SudokuValue getValue(int index) {
		return values[index].getValue();
	}
	
	@Override
	public SudokuValueSet getAllowedValues(int index) {
		return values[index].getAllowedValues();
	}

	@Override
	public SudokuValueSet getRowValues(int row) {
		return rowValues[row];
	}

	@Override
	public SudokuValueSet getColumnValues(int column) {
		return columnValues[column];
	}

	@Override
	public SudokuValueSet getSquareValues(int square) {
		return squareValues[square];
	}

	@Override
	public int getValueCount(SudokuValue v) {
		int vi = 0;
		if (!SudokuValue.isEmpty(v)) {
			vi = v.getValue();
		}
		return valueStats[vi];
	}

	@Override
	public int getFreeCellCount() {
		return getValueCount(null);
	}

	@Override
	public boolean isReadOnly() {
		return readOnly;
	}

	@Override
	public boolean isReadOnlyCell(int index) {
		return values[index].isReadOnly();
	}

	@Override
	public boolean isValidCell(int index) {
		return values[index].isValid();
	}

	@Override
	public boolean isValid() {
		return valid;
	}

	@Override
	public boolean isSolved() {
		return getFreeCellCount() == 0 && isValid();
	}

	private void incValueCount(SudokuValue v) {
		int vi = 0;
		if (!SudokuValue.isEmpty(v)) {
			vi = v.getValue();
		}
		++valueStats[vi];
	}

	private void decValueCount(SudokuValue v) {
		int vi = 0;
		if (!SudokuValue.isEmpty(v)) {
			vi = v.getValue();
		}
		--valueStats[vi];
	}

	private void assertIntegrity(boolean condition) {
		if (condition)
			return;

		throw new IllegalStateException();
	}

	@Override
	public String toString() {
		StringBuilder sb = new StringBuilder();
		SudokuFormatter.appendValues(sb, this);
		return sb.toString();
	}

	@Override
	public String toString(String format) {
		if (format == null || format.length() == 0)
			return toString();

		StringBuilder sb = new StringBuilder();
		for (int i = 0; i < format.length(); i++) {
			char c = format.charAt(i);
			switch (Character.toLowerCase(c)) {
			case ' ':
				sb.append(System.lineSeparator());
				break;
			case 'v':
				SudokuFormatter.appendValues(sb, this);
				break;
			case 'r':
				SudokuFormatter.appendRowValues(sb, this);
				break;
			case 'c':
				SudokuFormatter.appendColumnValues(sb, this);
				break;
			case 's':
				SudokuFormatter.appendSquareValues(sb, this);
				break;
			default:
				// not sure what the right type is here ...
				throw new FormatFlagsConversionMismatchException(format, c);
			}
		}

		return sb.toString();
	}

	private void checkIntegrity() {
		if (!SudokuConfiguration.getEnableBoardIntegrityCheck())
			return;

		for (int i = 0; i < BOARD_SIZE; i++) {
			SudokuValue value = this.getValue(i);

			// ensure value does not appear
			if (!SudokuValue.isEmpty(value)) {
				// check this value is disallowed in other places
				int[] relatedIndicies = SudokuIndexes.getRelated(i);
				for (int related : relatedIndicies) {
					if (values[related].isEmpty()) {
						assertIntegrity(!values[related].isAllowed(value));
					} else if (values[related].getValue() == value) {
						// ensure one of them is marked as invalid
						if (!values[related].isReadOnly())
							assertIntegrity(!values[related].isValid());
						if (!values[i].isReadOnly())
							assertIntegrity(!values[i].isValid());
						if (values[related].isReadOnly() && values[i].isReadOnly()) {
							assertIntegrity(!values[related].isValid());
							assertIntegrity(!values[i].isValid());
						}
					}
				}
			} else {
				// check that disallowed values are a union of row/column/square
				int row = SudokuIndexes.rowFromIndex(i);
				int col = SudokuIndexes.columnFromIndex(i);
				int square = SudokuIndexes.squareFromIndex(i);

				int disallowedValuesExpected = SudokuValue.combineMask(rowValues[row].getMask(),
						columnValues[col].getMask(), squareValues[square].getMask(), userDisallowed[i].getMask());

				assertIntegrity(values[i].getDisallowedValues().getMask() == disallowedValuesExpected);
			}
		}

		for (int sequenceIndex = 0; sequenceIndex < SudokuBoard.SEQUENCE_SIZE; sequenceIndex++) {
			checkAggregateValues(squareValues, sequenceIndex, SudokuIndexes.getSquareIndexes(sequenceIndex));
			checkAggregateValues(rowValues, sequenceIndex, SudokuIndexes.getRowIndexes(sequenceIndex));
			checkAggregateValues(columnValues, sequenceIndex, SudokuIndexes.getColumnIndexes(sequenceIndex));
		}
	}

	private void checkAggregateValues(SudokuValueSet[] stats, int sequenceIndex, int[] indicies) {
		// need to recalculate the statistics from scratch for the given
		// sequence as there might be duplicates there before and/or now...
		int mask = 0;
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			SudokuValue value = getValue(indicies[i]);
			if (!SudokuValue.isEmpty(value))
				mask = SudokuValue.combineMask(mask, value.getMask());
		}

		assertIntegrity(stats[sequenceIndex].getMask() == mask);
	}

	@Override
	public void setValue(int index, SudokuValue value, boolean isReadOnly) {
		if (readOnly) {
			throw new IllegalStateException("This game is read-only");
		}

		if (!editMode && (values[index].isReadOnly() || isReadOnly)) {
			throw new IllegalArgumentException("Edit mode is not allowed");
		}

		if (SudokuValue.isEmpty(value) && isReadOnly) {
			throw new IllegalArgumentException("Empty cell cannot be read-only");
		}

		SudokuValue previousValue = values[index].getValue();
		if (previousValue == value) {
			values[index].setReadOnly(isReadOnly);
			return;
		}

		if (SudokuValue.isEmpty(value)) {
			values[index].reset();
		} else {
			values[index].setValue(value, isReadOnly);
		}

		decValueCount(previousValue);
		incValueCount(value);

		// recalculate state of related row/column/square
		int row = SudokuIndexes.rowFromIndex(index);
		int column = SudokuIndexes.columnFromIndex(index);
		int square = SudokuIndexes.squareFromIndex(index);

		if (!valid) {
			// fully recalculate the state
			valid = true;

			recalculateAggregateValues(rowValues, row, SudokuIndexes.getRowIndexes(row));
			recalculateAggregateValues(columnValues, column, SudokuIndexes.getColumnIndexes(column));
			recalculateAggregateValues(squareValues, square, SudokuIndexes.getSquareIndexes(square));

			valid = recalculateAllowedValuesFromAggregates();
		} else {
			updateAggregateValues(rowValues, SudokuIndexes.getRowIndexes(row), index, row, previousValue);
			updateAggregateValues(columnValues, SudokuIndexes.getColumnIndexes(column), index, column, previousValue);
			updateAggregateValues(squareValues, SudokuIndexes.getSquareIndexes(square), index, square, previousValue);

			valid = updateAllowedValuesFromAggregates(index, previousValue);
		}

		checkIntegrity();
	}

	private boolean disallowValue(int[] indexes, int modifiedIndex, SudokuValue value) {
		boolean valid = true;
		for (int i = 0; i < indexes.length; i++) {
			int target = indexes[i];
			if (target == modifiedIndex)
				throw new IllegalArgumentException("should not have modified index in related indicies");

			if (values[target].isEmpty()) {
				values[target].disallow(value);
			} else if (values[target].getValue() == value) {
				// collision detected with new value
				valid = false;

				if (values[target].isReadOnly()) {
					values[modifiedIndex].setValid(false);
					if (values[modifiedIndex].isReadOnly()) {
						// both read-only? mark both as invalid
						values[target].setValid(false);
					}
				} else {
					values[target].setValid(false);
					if (!values[modifiedIndex].isReadOnly()) {
						values[modifiedIndex].setValid(false);
					}
				}
			}
		}

		return valid;
	}

	private boolean validate(int[] indexes) {
		boolean valid = true;
		if (indexes.length != SudokuBoard.SEQUENCE_SIZE)
			throw new IllegalArgumentException();

		int seenValuesMask = 0;
		int dupeValuesMask = 0;

		// pass 1 - get mask of existing values
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			int target = indexes[i];
			SudokuValue value = values[target].getValue();
			if (SudokuValue.isEmpty(value))
				continue;

			int mask = value.getMask();
			if ((mask & seenValuesMask) == 0) {
				seenValuesMask |= mask;
			} else {
				dupeValuesMask |= mask;
			}
		}

		if (dupeValuesMask != 0) {
			valid = false;

			// pass 2 - mark duplicate cells as invalid
			boolean atLeastOneInvalid = false;
			boolean checkReadOnly = true;
			do {
				for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
					int target = indexes[i];
					if (values[target].isEmpty() || (checkReadOnly && values[target].isReadOnly()))
						continue;

					int mask = values[target].getValue().getMask();
					if ((mask & dupeValuesMask) != 0) {
						values[target].setValid(false);
						atLeastOneInvalid = true;
					}
				}
				if (!checkReadOnly)
					break;// no third retry
				else
					checkReadOnly = false; // if we have to go second time
			} while (!atLeastOneInvalid);

			if (checkReadOnly && !atLeastOneInvalid)
				throw new IllegalStateException("should never happen");
		}

		return valid;
	}

	private void updateAllowedValuesFromAggregates(int index) {
		if (values[index].isEmpty()) {
			int row = SudokuIndexes.rowFromIndex(index);
			int col = SudokuIndexes.columnFromIndex(index);
			int square = SudokuIndexes.squareFromIndex(index);

			int seenValuesMask = rowValues[row].getMask() | columnValues[col].getMask()
					| squareValues[square].getMask();
			seenValuesMask |= userDisallowed[index].getMask();

			values[index].resetAllowed();
			values[index].disallowValueMask(seenValuesMask);
		}
	}

	private void updateAllowedValuesFromAggregates(int[] indices) {
		for (int i = 0; i < indices.length; i++) {
			updateAllowedValuesFromAggregates(indices[i]);
		}
	}

	private boolean recalculateAllowedValuesFromAggregates() {
		// first, update all allowed values for empty cells and reset Valid flag
		// for non-empty ones
		for (int i = 0; i < BOARD_SIZE; i++) {
			values[i].setValid(true);
			updateAllowedValuesFromAggregates(i);
		}

		boolean valid = true;
		// validate all
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			valid &= validate(SudokuIndexes.getRowIndexes(i));
			valid &= validate(SudokuIndexes.getColumnIndexes(i));
			valid &= validate(SudokuIndexes.getSquareIndexes(i));
		}

		return valid;
	}

	private boolean updateAllowedValuesFromAggregates(int modifiedIndex, SudokuValue previousValue) {
		if (!valid)
			throw new IllegalStateException("Use RecalculateAllowedValuesFromAggregates instead");

		boolean valid = true;
		if (SudokuValue.isEmpty(previousValue)) {
			// previous state is healthy and new value is added in empty spot -
			// safe to perform delta update only
			SudokuValue newValue = values[modifiedIndex].getValue();
			if (!SudokuValue.isEmpty(newValue)) {
				valid &= disallowValue(SudokuIndexes.getRelated(modifiedIndex), modifiedIndex, newValue);
			}
		} else {
			updateAllowedValuesFromAggregates(SudokuIndexes.getRelated(modifiedIndex));

			// update for self as well (if empty)
			updateAllowedValuesFromAggregates(modifiedIndex);

			int row = SudokuIndexes.rowFromIndex(modifiedIndex);
			int column = SudokuIndexes.columnFromIndex(modifiedIndex);
			int square = SudokuIndexes.squareFromIndex(modifiedIndex);

			valid &= validate(SudokuIndexes.getRowIndexes(row));
			valid &= validate(SudokuIndexes.getColumnIndexes(column));
			valid &= validate(SudokuIndexes.getSquareIndexes(square));
		}

		return valid;
	}

	private void updateAggregateValues(SudokuValueSet[] stats, int[] indicies, int modifiedIndex, int sequenceIndex,
			SudokuValue previousValue) {
		if (!valid)
			throw new IllegalStateException("Use RecalculateAggregateValues instead");

		// previous state is healthy - safe to perform delta update only
		if (!SudokuValue.isEmpty(previousValue)) {
			stats[sequenceIndex] = SudokuValueSet.remove(stats[sequenceIndex], previousValue);
		}
		SudokuValue newValue = values[modifiedIndex].getValue();
		if (!SudokuValue.isEmpty(newValue)) {
			stats[sequenceIndex] = SudokuValueSet.add(stats[sequenceIndex], newValue);
		}
	}

	private void recalculateAggregateValues() {
		for (int sequenceIndex = 0; sequenceIndex < SudokuBoard.SEQUENCE_SIZE; sequenceIndex++) {
			recalculateAggregateValues(rowValues, sequenceIndex, SudokuIndexes.getRowIndexes(sequenceIndex));
			recalculateAggregateValues(columnValues, sequenceIndex, SudokuIndexes.getColumnIndexes(sequenceIndex));
			recalculateAggregateValues(squareValues, sequenceIndex, SudokuIndexes.getSquareIndexes(sequenceIndex));
		}
	}

	private void recalculateAggregateValues(SudokuValueSet[] stats, int sequenceIndex, int[] indexes) {
		// need to recalculate the statistics from scratch for the given
		// sequence as there
		// might be duplicates there before and/or now...
		int mask = 0;
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			SudokuBoardCell cell = values[indexes[i]];
			if (cell.hasValue()) {
				mask |= cell.getValue().getMask();
			}
		}

		stats[sequenceIndex] = SudokuValueSet.buildFromMask(mask);
	}

	@Override
	public void restart() {
		for (int i = 0; i < values.length; i++) {
			if (!values[i].isReadOnly()) {
				decValueCount(values[i].getValue());
				values[i].reset();
				incValueCount(values[i].getValue());
			}
		}

		for (int i = 0; i < userDisallowed.length; i++) {
			userDisallowed[i] = SudokuValueSet.buildEmpty();
		}

		recalculateAggregateValues();
		valid = recalculateAllowedValuesFromAggregates();

		checkIntegrity();
	}
}
