package com.darkevilsudoku.board;

public enum SudokuValue {
	ONE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE;

	public static final int FULL_MASK = 0x1FF;
	public static final SudokuValue EMPTY = null;

	public String toString() {
		return Integer.toString(this.ordinal() + 1);
	}	
	public static boolean isEmpty(SudokuValue v) {
		return v == EMPTY;
	}

	public byte getValue() {
		return (byte) (1 + this.ordinal());
	}

	public int getMask() {
		return maskOf(getValue());
	}
	
	public static SudokuValue fromValue(byte v1) {
		switch (v1) {
		case 0:
			return EMPTY;
		case 1:
			return ONE;
		case 2:
			return TWO;
		case 3:
			return THREE;
		case 4:
			return FOUR;
		case 5:
			return FIVE;
		case 6:
			return SIX;
		case 7:
			return SEVEN;
		case 8:
			return EIGHT;
		case 9:
			return NINE;
		default:
			throw new IllegalArgumentException();
		}
	}

	public static int maskOf(byte value) {
		if (value < 1 || value > 9)
			throw new IllegalArgumentException();
		return 1 << (value - 1);
	}

	public static int combineMask(int mask1, int mask2) {
		return mask1 | mask2;
	}

	public static int combineMask(int mask1, int mask2, int mask3) {
		return mask1 | mask2 | mask3;
	}

	public static int combineMask(int mask1, int mask2, int mask3, int mask4) {
		return mask1 | mask2 | mask3 | mask4;
	}

	public static int removeMask(int mask, int maskToRemove) {
		return mask & complementMask(maskToRemove);
	}

	public static int complementMask(int mask) {
		return FULL_MASK & ~mask;
	}

	public SudokuValue next() {
		if (getValue() == 9) {
			throw new IllegalStateException();
		}

		return SudokuValue.fromValue((byte)(getValue() + 1));
	}

	public boolean isLast() {
		return this == NINE;
	}
}
