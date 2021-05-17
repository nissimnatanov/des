package test.darkevilsudoku.board;

import static org.junit.Assert.*;

import org.junit.Test;

import com.darkevilsudoku.board.SudokuValue;
import com.darkevilsudoku.board.SudokuValueSet;

public class SudokuValueSetTests {

	@Test
	public final void testHashCode() {
		SudokuValueSet vs1 = SudokuValueSet.build(SudokuValue.FIVE);
		vs1 = SudokuValueSet.add(vs1, SudokuValue.EIGHT);

		SudokuValueSet vs2 = SudokuValueSet.build(SudokuValue.EIGHT);
		vs2 = SudokuValueSet.add(vs2, SudokuValue.FIVE);

		assertEquals(vs1.hashCode(), vs2.hashCode());
	}

	@Test
	public final void testBuildEmpty() {
		SudokuValueSet vs1 = SudokuValueSet.buildEmpty();

		validateSet(vs1, "");
	}

	@Test
	public final void testBuild() {
		SudokuValueSet vs1 = SudokuValueSet.build(SudokuValue.FIVE);

		assertEquals(vs1.size(), 1);
		assertEquals(vs1.getCombined(), 5);

		validateSet(vs1, "5");
	}

	@Test
	public final void testBuildFull() {
		SudokuValueSet vs1 = SudokuValueSet.buildFull();

		validateSet(vs1, "123456789");
	}

	@Test
	public final void testAdd() {
		SudokuValueSet vs = SudokuValueSet.buildFull();
		vs = SudokuValueSet.add(vs, SudokuValue.EIGHT);
		validateSet(vs, "123456789");

		vs = SudokuValueSet.buildEmpty();

		vs = SudokuValueSet.add(vs, SudokuValue.EIGHT);
		validateSet(vs, "8");

		vs = SudokuValueSet.add(vs, SudokuValue.EIGHT);
		validateSet(vs, "8");

		vs = SudokuValueSet.add(vs, SudokuValue.FIVE);
		validateSet(vs, "58");

		vs = SudokuValueSet.add(vs, SudokuValue.SEVEN);
		validateSet(vs, "578");

		vs = SudokuValueSet.add(vs, SudokuValue.ONE);
		validateSet(vs, "1578");

		vs = SudokuValueSet.add(vs, SudokuValue.THREE);
		validateSet(vs, "13578");

		vs = SudokuValueSet.add(vs, SudokuValue.SIX);
		validateSet(vs, "135678");

		vs = SudokuValueSet.add(vs, SudokuValue.NINE);
		validateSet(vs, "1356789");

		vs = SudokuValueSet.add(vs, SudokuValue.TWO);
		validateSet(vs, "12356789");

		vs = SudokuValueSet.add(vs, SudokuValue.FOUR);
		validateSet(vs, "123456789");
	}

	@Test
	public final void testAdd2() {
		SudokuValueSet vs1 = SudokuValueSet.buildEmpty();
		SudokuValueSet vs2 = SudokuValueSet.buildFull();
		SudokuValueSet vs = SudokuValueSet.add(vs1, vs2);
		
		validateSet(vs1, "");
		validateSet(vs, "123456789");
		
		vs1 = getOddSet();
		vs2 = SudokuValueSet.build(SudokuValue.FOUR);
		vs = SudokuValueSet.add(vs1, vs2);
		validateSet(vs, "134579");
	}
	
	@Test
	public final void testRemove() {
		SudokuValueSet vs = SudokuValueSet.buildFull();
		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT);
		validateSet(vs, "12345679");

		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT);
		validateSet(vs, "12345679");

		vs = SudokuValueSet.remove(vs, SudokuValue.FIVE);
		validateSet(vs, "1234679");

		vs = SudokuValueSet.buildEmpty();

		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT);
		validateSet(vs, "");
	}

	@Test
	public final void testRemove2() {
		SudokuValueSet vs = SudokuValueSet.buildFull();
		vs = SudokuValueSet.remove(vs, SudokuValueSet.build(SudokuValue.EIGHT));
		validateSet(vs, "12345679");

		vs = SudokuValueSet.remove(vs, SudokuValueSet.build(SudokuValue.SEVEN));
		validateSet(vs, "1234569");

		vs = SudokuValueSet.remove(vs, getOddSet());
		validateSet(vs, "246");

		vs = SudokuValueSet.remove(vs, SudokuValueSet.buildEmpty());
		validateSet(vs, "246");
	}

	@Test
	public final void testRemoveMask() {
		SudokuValueSet vs = SudokuValueSet.buildFull();
		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT.getMask());
		validateSet(vs, "12345679");

		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT.getMask() | SudokuValue.FIVE.getMask());
		validateSet(vs, "1234679");

		vs = SudokuValueSet.remove(vs, SudokuValue.THREE.getMask());
		validateSet(vs, "124679");

		vs = SudokuValueSet.buildEmpty();

		vs = SudokuValueSet.remove(vs, SudokuValue.FULL_MASK);
		validateSet(vs, "");
	}

	@Test
	public final void testComplement() {
		SudokuValueSet vs = SudokuValueSet.buildFull();
		SudokuValueSet vs_complement = vs.getComplement();
		validateSet(vs_complement, "");

		vs = SudokuValueSet.remove(vs, SudokuValue.EIGHT);
		validateSet(vs_complement, "");
		vs_complement = vs.getComplement();
		validateSet(vs.getComplement(), "8");
	}

	@Test
	public final void testEqualsObject() {
		SudokuValueSet vs1 = SudokuValueSet.buildFull();
		SudokuValueSet vs2 = SudokuValueSet.buildFull();
		assertTrue(vs1.equals(vs2));
		assertTrue(vs1.equals((Object) vs2));
		assertFalse(vs1.equals(null));
		assertFalse(vs1.equals(SudokuValueSet.buildEmpty()));
		assertFalse(vs1.equals(SudokuValueSet.build(SudokuValue.SEVEN)));

		vs1 = SudokuValueSet.add(vs1, SudokuValue.SIX);
		assertTrue(vs1.equals(vs2));

		vs1 = SudokuValueSet.remove(vs1, SudokuValue.SIX);
		assertFalse(vs1.equals(vs2));

		vs1 = SudokuValueSet.remove(vs1, SudokuValue.ONE);
		vs2 = SudokuValueSet.remove(vs2, SudokuValue.ONE);
		vs2 = SudokuValueSet.remove(vs2, SudokuValue.SIX);
		assertTrue(vs1.equals(vs2));
	}
	
	private final SudokuValueSet getOddSet() {
		int oddMask = 
				SudokuValue.ONE.getMask() | 
				SudokuValue.THREE.getMask() |
				SudokuValue.FIVE.getMask() |
				SudokuValue.SEVEN.getMask() | 
				SudokuValue.NINE.getMask();
		return SudokuValueSet.buildFromMask(oddMask);
	}

	private final void validateSet(SudokuValueSet vs, String expected) {
		if (expected == null)
			assertNull(vs);
		else
			assertNotNull(vs);

		assertEquals("size does not match", expected.length(), vs.size());
		assertEquals("size of complement does not match", 9 - expected.length(), vs.getComplement().size());
		String expectedCombined = expected.length() > 0 ? expected : "0";
		assertEquals("toString does not match", expectedCombined, vs.toString());
		assertEquals("combined does not match", Integer.parseInt(expectedCombined), vs.getCombined());

		for(int i=0; i<vs.size(); i++) {
			SudokuValue v = vs.getAt(i);
			assertTrue("unexpected value " + v.toString(), i < expected.length());
			assertEquals("enumerated value does not match", expected.substring(i, i + 1), v.toString());
			i++;
		}
	}
}
