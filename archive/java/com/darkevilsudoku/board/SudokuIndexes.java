package com.darkevilsudoku.board;

import java.util.ArrayList;
import java.util.List;

public final class SudokuIndexes {
	private static final int[][] rows;
	private static final int[][] columns;
	private static final int[][] squares;
	private static final int[][] related;
	private static final int[] indexToSquare;
	private static final int[] indexToCellIndexSquare;

	static {
		rows = new int[SudokuBoard.SEQUENCE_SIZE][];
		columns = new int[SudokuBoard.SEQUENCE_SIZE][];
		squares = new int[SudokuBoard.SEQUENCE_SIZE][];
		related = new int[SudokuBoard.BOARD_SIZE][];
		indexToSquare = new int[SudokuBoard.BOARD_SIZE];
		indexToCellIndexSquare = new int[SudokuBoard.BOARD_SIZE];

		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			rows[i] = new int[SudokuBoard.SEQUENCE_SIZE];
			columns[i] = new int[SudokuBoard.SEQUENCE_SIZE];
			squares[i] = new int[SudokuBoard.SEQUENCE_SIZE];

			for (int j = 0; j < SudokuBoard.SEQUENCE_SIZE; j++) {
				// set row indexes
				// 0 .. 8; 9 .. 17; ...
				rows[i][j] = i * SudokuBoard.SEQUENCE_SIZE + j;

				// set column indexes
				// { 0, 9, 18 ...} {1, 10, 19 ...} ...
				columns[i][j] = i + j * SudokuBoard.SEQUENCE_SIZE;

				// set square indexes
				// { 0, 1, 2, 9, 10, 11, ...}
				int index = 0;
				index += (i / 3) * 27; // note that i / 3 gives integer value so
										// the expression cannot be replaced
										// with i * 9
				index += (j / 3) * 9;
				index += ((i % 3) * 3 + (j % 3));

				squares[i][j] = index;

				// set map between index to square
				indexToSquare[index] = i;
				indexToCellIndexSquare[index] = j;
			}
		}

		// get the related indexes for each cell
		for (int i = 0; i < SudokuBoard.BOARD_SIZE; i++) {
			List<Integer> relatedList = new ArrayList<Integer>(20);
			for (int relatedIndex : rows[i / 9]) {
				if (relatedIndex == i)
					continue;
				if (!relatedList.contains(relatedIndex))
					relatedList.add(relatedIndex);
			}
			for (int relatedIndex : columns[i % 9]) {
				if (relatedIndex == i)
					continue;
				if (!relatedList.contains(relatedIndex))
					relatedList.add(relatedIndex);
			}
			for (int relatedIndex : squares[squareFromIndex(i)]) {
				if (relatedIndex == i)
					continue;
				if (!relatedList.contains(relatedIndex))
					relatedList.add(relatedIndex);
			}

			if (relatedList.size() != 20)
				throw new RuntimeException();
			related[i] = new int[20];
			for (int j = 0; j < 20; j++)
				related[i][j] = relatedList.get(j);
		}
	}

	public static int rowFromIndex(int index) {
		return index / 9;
	}

	public static int columnFromIndex(int index) {
		return index % 9;
	}

	public static int squareFromIndex(int index) {
		return indexToSquare[index];
	}

	public static int cellIndexInSquareFromIndex(int index) {
		return indexToCellIndexSquare[index];
	}

	public static int indexFromCoordinates(int row, int col) {
		return row * 9 + col;
	}

	public static int indexFromSquare(int squareIndex, int cellIndex) {
		return squares[squareIndex][cellIndex];
	}

	// TODO replace with immutable
	public static int[] getRelated(int index) {
		return related[index];
	}

	// TODO replace with immutable
	public static int[] getRowIndexes(int row) {
		return rows[row];
	}

	// TODO replace with immutable
	public static int[] getColumnIndexes(int column) {
		return columns[column];
	}

	// TODO replace with immutable
	public static int[] getSquareIndexes(int square) {
		return squares[square];
	}
	
	private SudokuIndexes() {
	}
}
