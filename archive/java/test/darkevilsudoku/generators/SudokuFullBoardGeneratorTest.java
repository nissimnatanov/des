package test.darkevilsudoku.generators;

import static org.junit.Assert.*;

import java.util.Random;

import org.junit.Test;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.board.SudokuConfiguration;
import com.darkevilsudoku.generators.SudokuFullBoardGenerator;

public class SudokuFullBoardGeneratorTest {
	
	@Test(timeout = 1200)
	public void testNewBoard_Speed() {
		SudokuConfiguration.setEnableBoardIntegrityCheck(false);
		generateBoards(100);
	}
	
	@Test
	public void testNewBoard_Accuracy() {
		SudokuConfiguration.setEnableBoardIntegrityCheck(true);
		try {
			generateBoards(100);
		}
		finally {
			SudokuConfiguration.setEnableBoardIntegrityCheck(false);
		}
	}
	
	private void generateBoards(int count) {
		Random r = new Random();
		for(int i=0; i< count; i++) {
			
			SudokuBoard b =  SudokuFullBoardGenerator.newBoard(r);
			
			assertTrue(b.isValid());
			assertTrue(b.getFreeCellCount() == 0);
		}
	}
}
