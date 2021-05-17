package com.darkevilsudoku.board;

public final class SudokuFormatter {
	public static final int READONLY = 0x1;
	public static final int INVALID = 0x2;
	public static final int FULLWIDTH = 0x4;
	public static final int EMPTY_AS_ZERO = 0x8;
	public static final int ASCII_ONLY = 0x10;

	// empty cells
	private static final Character SPACE_CHAR = ' ';
	private static final Character FULLWIDTH_SPACE_CHAR = '\u3000';
	private static final Character ZERO_CHAR = '0';
	private static final Character FULLWIDTH_ZERO_CHAR = 0xFF10;

	// editable cells
	private static final Character UNDERLINE = '_';
	private static final Character UNDERLINE_OVERLAY = '\u0332';
	
	// wrong cells
	private static final Character X = 'X';
	private static final Character STRIKEOUT_OVERLAY = '\u0336';
	private static final Character SOLIDUS_OVERLAY = '\u0337';

	private static final String SPACE = SPACE_CHAR.toString();
	private static final String FULLWIDTH_SPACE = FULLWIDTH_SPACE_CHAR.toString();

	private static final String[] READONLY_DIGITS;
	private static final String[] FULLWIDTH_READONLY_DIGITS;

	private static final String[] WRITABLE_DIGITS;
	private static final String[] FULLWIDTH_WRITABLE_DIGITS;

	private static final String[] INVALID_READONLY_DIGITS;
	private static final String[] FULLWIDTH_INVALID_READONLY_DIGITS;

	private static final String[] INVALID_WRITABLE_DIGITS;
	private static final String[] FULLWIDTH_INVALID_WRITABLE_DIGITS;

	static {
		READONLY_DIGITS = new String[10];
		FULLWIDTH_READONLY_DIGITS = new String[10];

		WRITABLE_DIGITS = new String[10];
		FULLWIDTH_WRITABLE_DIGITS = new String[10];

		INVALID_READONLY_DIGITS = new String[10];
		FULLWIDTH_INVALID_READONLY_DIGITS = new String[10];

		INVALID_WRITABLE_DIGITS = new String[10];
		FULLWIDTH_INVALID_WRITABLE_DIGITS = new String[10];

		for (int i = 0; i < 10; i++) {
			READONLY_DIGITS[i] = new String(Character.toChars(ZERO_CHAR + i));
			FULLWIDTH_READONLY_DIGITS[i] = new String(Character.toChars(FULLWIDTH_ZERO_CHAR + i));

			WRITABLE_DIGITS[i] = READONLY_DIGITS[i] + UNDERLINE_OVERLAY;
			FULLWIDTH_WRITABLE_DIGITS[i] = FULLWIDTH_READONLY_DIGITS[i] + UNDERLINE_OVERLAY;

			INVALID_READONLY_DIGITS[i] = READONLY_DIGITS[i] + SOLIDUS_OVERLAY;
			FULLWIDTH_INVALID_READONLY_DIGITS[i] = FULLWIDTH_READONLY_DIGITS[i] + SOLIDUS_OVERLAY;

			INVALID_WRITABLE_DIGITS[i] = READONLY_DIGITS[i] + STRIKEOUT_OVERLAY;
			FULLWIDTH_INVALID_WRITABLE_DIGITS[i] = FULLWIDTH_READONLY_DIGITS[i] + STRIKEOUT_OVERLAY;
		}
	}

	public static boolean isValue(char c) {
		return (c > ZERO_CHAR && c <= (ZERO_CHAR + 9)) || (c > ZERO_CHAR && c <= (FULLWIDTH_ZERO_CHAR + 9));
	}

	public static SudokuValue getValue(char c) {
		if (c > ZERO_CHAR && c <= (ZERO_CHAR + 9)) {
			return SudokuValue.fromValue((byte) (c - ZERO_CHAR));
		}

		if (c > FULLWIDTH_ZERO_CHAR && c <= (FULLWIDTH_ZERO_CHAR + 9)) {
			return SudokuValue.fromValue((byte) (c - FULLWIDTH_ZERO_CHAR));
		}

		throw new IllegalArgumentException();
	}

	public static boolean isEmpty(char c, boolean ignoreSpaces) {
		if (c == ZERO_CHAR || c == FULLWIDTH_ZERO_CHAR) {
			return true;
		}

		if (!ignoreSpaces && (c == SPACE_CHAR || c == FULLWIDTH_SPACE_CHAR)) {
			return true;
		}

		return false;
	}

	public static boolean containsZeros(String s) {
		return s.contains(READONLY_DIGITS[0]) || s.contains(FULLWIDTH_READONLY_DIGITS[0]);
	}

	public static boolean isEditableIndicator(char c) {
		return c == UNDERLINE_OVERLAY || c == UNDERLINE;
	}

	public static boolean isValidityIndicator(char c) {
		return c == STRIKEOUT_OVERLAY || c == SOLIDUS_OVERLAY || c == X;
	}

	public static String toString(SudokuValue v, int options) {
		boolean readOnly = (options & READONLY) != 0;
		boolean fullWidth = (options & FULLWIDTH) != 0;
		boolean asciiOnly = (options & ASCII_ONLY) != 0;
		boolean invalid = (options & INVALID) != 0;
		boolean emptyAsZero = (options & EMPTY_AS_ZERO) != 0;
		
		if (asciiOnly && fullWidth) {
			throw new IllegalArgumentException();
		}

		if (readOnly && SudokuValue.isEmpty(v)) {
			throw new IllegalArgumentException();
		}

		if (SudokuValue.isEmpty(v)) {
			if (asciiOnly) {
				if (emptyAsZero) {
					return READONLY_DIGITS[0];
				} 
				else {
					return SPACE;
				}
			}
			else {
				if (emptyAsZero) {
					return fullWidth ? FULLWIDTH_WRITABLE_DIGITS[0] : WRITABLE_DIGITS[0];
				} 
				else {
					return fullWidth ? FULLWIDTH_SPACE : SPACE;
				}
			}
		}

		String[] digits;

		if (readOnly || asciiOnly) {
			if (!invalid || asciiOnly) {
				digits = fullWidth ? FULLWIDTH_READONLY_DIGITS : READONLY_DIGITS;
			} else {
				digits = fullWidth ? FULLWIDTH_INVALID_READONLY_DIGITS : INVALID_READONLY_DIGITS;
			}
		} else {
			if (!invalid) {
				digits = fullWidth ? FULLWIDTH_WRITABLE_DIGITS : WRITABLE_DIGITS;
			} else {
				digits = fullWidth ? FULLWIDTH_INVALID_WRITABLE_DIGITS : INVALID_WRITABLE_DIGITS;
			}
		}
		String ret = digits[v.getValue()];
		if (asciiOnly) {
			// mark editable value with underscore and invalid cells with X
			if (invalid) {
				ret += X;
			}
			if (!readOnly) {
				ret += UNDERLINE;
			}
		}
		return ret;
	}

	public static void appendValues(StringBuilder sb, SudokuBoard board) {
		sb.append(System.lineSeparator());
		// using UNICODE characters here instead of ASCII
		sb.append("／－－－－－＋－－－－－＋－－－－－＼");
		sb.append(System.lineSeparator());

		for (int row = 0; row < SudokuBoard.SEQUENCE_SIZE; row++) {
			if (row != 0 && (row % 3) == 0) {
				sb.append("＋－－－－－＋－－－－－＋－－－－－＋");
				sb.append(System.lineSeparator());
			}
			for (int col = 0; col < SudokuBoard.SEQUENCE_SIZE; col++) {
				if (col % 3 == 0) {
					sb.append("｜");
				}
				int i = SudokuIndexes.indexFromCoordinates(row, col);
				int options = FULLWIDTH;
				if (board.isReadOnlyCell(i)) {
					options |= READONLY;
				}
				if (!board.isValidCell(i)) {
					options |= INVALID;
				}

				sb.append(toString(board.getValue(i), options));
				if (col % 3 != 2) {
					sb.append(FULLWIDTH_SPACE);
				}
			}

			sb.append("｜");
			sb.append(System.lineSeparator());
		}

		sb.append("＼－－－－－＋－－－－－＋－－－－－／");
		sb.append(System.lineSeparator());
	}

	public static void appendRowValues(StringBuilder sb, SudokuBoard board) {
		sb.append(System.lineSeparator());
		sb.append("Row values:");
		sb.append(System.lineSeparator());
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			sb.append(" [r");
			sb.append(i);
			sb.append("] ");
			sb.append(board.getRowValues(i));
			sb.append(System.lineSeparator());
		}
	}

	public static void appendColumnValues(StringBuilder sb, SudokuBoard board) {
		sb.append(System.lineSeparator());
		sb.append("Column values:");
		sb.append(System.lineSeparator());
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			sb.append(" [c");
			sb.append(i);
			sb.append("] ");
			sb.append(board.getColumnValues(i));
			sb.append(System.lineSeparator());
		}
	}

	public static void appendSquareValues(StringBuilder sb, SudokuBoard board) {
		sb.append(System.lineSeparator());
		sb.append("Square values:");
		sb.append(System.lineSeparator());
		for (int i = 0; i < SudokuBoard.SEQUENCE_SIZE; i++) {
			sb.append(" [s");
			sb.append(i);
			sb.append("] ");
			sb.append(board.getSquareValues(i));
			sb.append(System.lineSeparator());
		}
	}

	private SudokuFormatter() {
	}
}
