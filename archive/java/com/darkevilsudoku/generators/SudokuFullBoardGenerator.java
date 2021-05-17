package com.darkevilsudoku.generators;

import java.util.Random;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.generators.impl.SudokuFullBoardGeneratorImpl;

public abstract class SudokuFullBoardGenerator {
	public abstract SudokuBoard generate(Random r);


	private static final SudokuFullBoardGenerator instance = new SudokuFullBoardGeneratorImpl();
	public static final SudokuFullBoardGenerator getInstance() {
		return instance;
	}

	public static SudokuBoard newBoard(Random r) {
		return getInstance().generate(r);
	}

	public static SudokuBoard newBoard() {
		return getInstance().generate(null);
	}
}
