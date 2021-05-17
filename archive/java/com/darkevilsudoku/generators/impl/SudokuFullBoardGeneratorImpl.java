package com.darkevilsudoku.generators.impl;

import java.util.Random;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.board.SudokuIndexes;
import com.darkevilsudoku.board.SudokuValue;
import com.darkevilsudoku.board.SudokuValueSet;
import com.darkevilsudoku.generators.SudokuFullBoardGenerator;
import com.darkevilsudoku.utils.RandomUtils;

public class SudokuFullBoardGeneratorImpl extends SudokuFullBoardGenerator {

	private static void setSquareValues(SudokuBoard board, int square, int cellIndex, SudokuValue[] values,
			int startValues, int count, boolean isReadOnly) {
		while (count-- > 0) {
			board.setValue(SudokuIndexes.indexFromSquare(square, cellIndex), values[startValues++], isReadOnly);
			cellIndex++;
		}

		if (!board.isValid())
			throw new IllegalStateException();
	}

	private static boolean tryFillSquare(Random r, SudokuBoard board, int square) {
		int[] cellIndices = new int[9];
		for (int cell = 0; cell < 9; cell++)
			cellIndices[cell] = cell;

		RandomUtils.shuffle(r, cellIndices, 0, cellIndices.length);

		int allowedRetries = 1;
		boolean retry;
		do {
			retry = false;
			for (int i = 0; i < 9; i++) {
				int cell = cellIndices[i];
				int index = SudokuIndexes.indexFromSquare(square, cell);
				SudokuValueSet allowedValues = board.getAllowedValues(index);
				boolean valid = (allowedValues.size() > 0);
				if (valid) {
					int next = r.nextInt(allowedValues.size());
					SudokuValue value = allowedValues.getAt(next);
					board.setValue(index, value, false);
					if (!board.isValid())
						valid = false;
				}

				if (!valid) {
					if (allowedRetries-- == 0)
						return false;

					// reset the square
					for (cell = 0; cell < 9; cell++) {
						board.setValue(SudokuIndexes.indexFromSquare(square, cell), SudokuValue.EMPTY, false);
					}
					retry = true;
					break;
				}
			}
		} while (retry);

		return true;
	}

	@Override
	public SudokuBoard generate(Random r) {
		if (r == null) {
			r = new Random();
		}

		SudokuValue[] values = SudokuValue.values();

		SudokuBoard board = SudokuBoard.build();

		// populate squares 0, 4 (middle) and 8 (last)
		// since there is no intersection between those, safe to use totally
		// random values
		RandomUtils.shuffle(r, values, 0, values.length);

		setSquareValues(board, 0, 0, values, 0, 9, true);
		RandomUtils.shuffle(r, values, 0, values.length);
		setSquareValues(board, 4, 0, values, 0, 9, true);
		RandomUtils.shuffle(r, values, 0, values.length);
		setSquareValues(board, 8, 0, values, 0, 9, true);

		do {
			board.restart();

			// fill squares 1 to 7, skipping 4
			for (int square = 1; square <= 7; square++) {
				if (square == 4)
					continue; // skip

				if (!tryFillSquare(r, board, square)) {
					break;
				}
			}
		} while (board.getFreeCellCount() > 0);

		for (int i = 0; i < 81; i++) {
			board.setValue(i, board.getValue(i), true);
		}

		return board;
	}
}
