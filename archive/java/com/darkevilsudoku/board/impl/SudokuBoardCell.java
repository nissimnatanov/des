package com.darkevilsudoku.board.impl;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.board.SudokuValue;
import com.darkevilsudoku.board.SudokuValueSet;

class SudokuBoardCell implements Cloneable {
	
	static SudokuBoardCell[] buildArray(int size) {
		SudokuBoardCell[] values = new SudokuBoardCell[SudokuBoard.BOARD_SIZE];
		for (int i = 0; i < values.length; i++) {
			values[i] = new SudokuBoardCell();
		}
		return values;
	}
	
	@Override
	public SudokuBoardCell clone() {
		SudokuBoardCell cell;
		try {
			cell = (SudokuBoardCell)super.clone();
		} catch (CloneNotSupportedException e) {
			throw new RuntimeException(e);
		}
		
		return cell;
	}
	
	private static final int VALID = 1;
	private static final int READONLY = 2;

	private SudokuValue value;
	private int state;
	private SudokuValueSet disallowedValues;

	SudokuBoardCell() {
		value = null;
		state = VALID;
		disallowedValues = SudokuValueSet.buildEmpty();
	}
	
	boolean isReadOnly() {
		return (state & READONLY) != 0;
	}

	void setReadOnly(boolean readOnly) {
		if (readOnly) {
			state |= READONLY;
		} else {
			state &= ~READONLY;
		}
	}

	boolean isValid() {
		return (state & VALID) != 0;
	}

	void setValid(boolean valid) {
		if (valid) {
			state |= VALID;
		} else {
			state &= ~VALID;
		}
	}

	boolean isEmpty() {
		return value == null;
	}

	private void validateEmpty() {
		if (!isEmpty())
			throw new RuntimeException();
	}
	
	boolean hasValue() {
		return getValue() != null;
	}

	SudokuValue getValue() {
		return value;
	}

	void setValue(SudokuValue value, boolean asReadOnly) {
		if (value == null)
			throw new NullPointerException();

		setReadOnly(asReadOnly);
		this.value = value;
	}

	SudokuValueSet getDisallowedValues() {
		validateEmpty();
		return this.disallowedValues;
	}

	SudokuValueSet getAllowedValues() {
		return this.getDisallowedValues().getComplement();
	}

	void disallow(SudokuValue value) {
		validateEmpty();
		disallowedValues = SudokuValueSet.add(disallowedValues, value);
	}

	void disallowValueMask(int valuesMask) {
		disallowedValues = SudokuValueSet.add(disallowedValues, valuesMask);
	}

	boolean isAllowed(SudokuValue value) {
		validateEmpty();
		return !disallowedValues.contains(value);
	}

	void reset() {
		value = null;
		state = 0;
		disallowedValues = SudokuValueSet.buildEmpty();
	}

	void resetAllowed() {
		validateEmpty();
		disallowedValues = SudokuValueSet.buildEmpty();
	}
}
