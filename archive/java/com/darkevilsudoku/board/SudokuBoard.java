package com.darkevilsudoku.board;

import com.darkevilsudoku.board.impl.SudokuBoardImpl;

public interface SudokuBoard extends Cloneable {
	public static final int BOARD_SIZE = 81;
	public static final int SEQUENCE_SIZE = 9; // length of the single row, column or square

	public SudokuBoard clone();
	
	public boolean hasValue(int index);
	public SudokuValue getValue(int index);
	public SudokuValueSet getAllowedValues(int index);

	public SudokuValueSet getRowValues(int row);
	public SudokuValueSet getColumnValues(int column);
	public SudokuValueSet getSquareValues(int square);

	public int getValueCount(SudokuValue v);
	public int getFreeCellCount();

	public boolean isReadOnly();
	public boolean isReadOnlyCell(int index);
	public boolean isValidCell(int index);	
	public boolean isValid();
	public boolean isSolved();
	
	public void setValue(int index, SudokuValue value, boolean isReadOnly);
	public void restart();

	public String toString(String format);
	
	public static SudokuBoard build() {
		return new SudokuBoardImpl();
	}
}
