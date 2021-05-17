package test.darkevilsudoku.board;

import static org.junit.Assert.*;

import java.util.logging.Logger;

import org.junit.Test;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.board.SudokuConfiguration;
import com.darkevilsudoku.board.SudokuBoardSerializer;
import com.darkevilsudoku.board.SudokuValue;
import com.darkevilsudoku.board.SudokuValueSet;
import com.darkevilsudoku.board.impl.SudokuBoardSerializerImpl;

public class SudokuBoardTests {
	private Logger logger = Logger.getGlobal();

	@Test
	public void testEmptyBoard() {
		SudokuConfiguration.setEnableBoardIntegrityCheck(true);
		
		SudokuBoard b = SudokuBoard.build();
		logger.info(b.toString());

		assertEquals(81, b.getFreeCellCount());
		assertEquals(SudokuValueSet.buildEmpty(), b.getRowValues(5));
		assertEquals(SudokuValueSet.buildEmpty(), b.getColumnValues(1));
		assertEquals(SudokuValueSet.buildEmpty(), b.getSquareValues(0));
		
		assertEquals(null, b.getValue(22));
		assertEquals(0, b.getValueCount(SudokuValue.FIVE));
		assertFalse(b.isReadOnlyCell(44));
		assertFalse(b.isReadOnlyCell(55));
		assertTrue(b.isValid());
		assertTrue(b.isValidCell(11));
		assertFalse(b.isSolved());

		b.restart();
	}
	
	@Test
	public void testBoardWithSomeValues() {
		SudokuConfiguration.setEnableBoardIntegrityCheck(true);

		SudokuBoard b = SudokuBoard.build();
		b.setValue(44, SudokuValue.EIGHT, false);
		b.setValue(67, SudokuValue.FIVE, true);
		b.setValue(80, SudokuValue.SEVEN, false);
		logger.info(b.toString());

		assertEquals(78, b.getFreeCellCount());
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getRowValues(7));
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getColumnValues(4));
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getSquareValues(7));

		assertEquals(SudokuValueSet.build(SudokuValue.EIGHT), b.getRowValues(4));
		assertEquals(SudokuValueSet.build(SudokuValue.EIGHT, SudokuValue.SEVEN), b.getColumnValues(8));
		assertEquals(SudokuValueSet.build(SudokuValue.EIGHT), b.getSquareValues(5));

		assertEquals(SudokuValueSet.buildEmpty(), b.getRowValues(5));
		assertEquals(SudokuValueSet.buildEmpty(), b.getColumnValues(1));
		assertEquals(SudokuValueSet.buildEmpty(), b.getSquareValues(0));
		
		assertEquals(null, b.getValue(22));
		assertEquals(SudokuValue.FIVE, b.getValue(67));
		assertEquals(SudokuValue.EIGHT, b.getValue(44));
		assertEquals(null, b.getValue(79));
		
		assertEquals(1, b.getValueCount(SudokuValue.FIVE));
		assertFalse(b.isReadOnlyCell(0));
		assertFalse(b.isReadOnlyCell(44));
		assertTrue(b.isReadOnlyCell(67));
		assertFalse(b.isReadOnlyCell(55));
		assertTrue(b.isValid());
		assertTrue(b.isValidCell(11));
		assertFalse(b.isSolved());
		
		b.restart();
		assertEquals(80, b.getFreeCellCount());
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getRowValues(7));
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getColumnValues(4));
		assertEquals(SudokuValueSet.build(SudokuValue.FIVE), b.getSquareValues(7));

		assertEquals(SudokuValueSet.buildEmpty(), b.getRowValues(4));
		assertEquals(SudokuValueSet.buildEmpty(), b.getColumnValues(8));
		assertEquals(SudokuValueSet.buildEmpty(), b.getSquareValues(5));

		assertEquals(SudokuValueSet.buildEmpty(), b.getRowValues(0));
		assertEquals(SudokuValueSet.buildEmpty(), b.getColumnValues(5));
		assertEquals(SudokuValueSet.buildEmpty(), b.getSquareValues(1));
		
		assertEquals(null, b.getValue(22));
		assertEquals(SudokuValue.FIVE, b.getValue(67));
		assertEquals(null, b.getValue(44));
		assertEquals(null, b.getValue(79));
		
		assertEquals(1, b.getValueCount(SudokuValue.FIVE));
		assertEquals(0, b.getValueCount(SudokuValue.EIGHT));

		assertFalse(b.isReadOnlyCell(0));
		assertFalse(b.isReadOnlyCell(33));
		assertFalse(b.isReadOnlyCell(44));
		assertFalse(b.isReadOnlyCell(55));
		assertTrue(b.isReadOnlyCell(67));
		assertTrue(b.isValid());
		assertTrue(b.isValidCell(11));
		assertFalse(b.isSolved());
	}
	
	@Test
	public void testInvalidBoardSimple() {
		SudokuConfiguration.setEnableBoardIntegrityCheck(true);

		SudokuBoard b = SudokuBoard.build();
		b.setValue(44, SudokuValue.EIGHT, false);
		
		b.setValue(67, SudokuValue.FIVE, true);
		b.setValue(77, SudokuValue.FIVE, true);

		b.setValue(80, SudokuValue.EIGHT, false);
		logger.info(b.toString());
		SudokuBoardSerializer serializer = SudokuBoardSerializerImpl.builder().build();
		logger.info(serializer.serialize(b));

		serializer = SudokuBoardSerializer.builder()
				.setAsciiOnly(true)
				.setCompact(true)
				.build();
		logger.info(serializer.serialize(b));
	}
}
