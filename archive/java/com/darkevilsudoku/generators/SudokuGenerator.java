package com.darkevilsudoku.generators;

import java.util.Random;
import com.darkevilsudoku.board.SudokuBoard;

public interface SudokuGenerator {
	SudokuBoard generate(Random r);

}
