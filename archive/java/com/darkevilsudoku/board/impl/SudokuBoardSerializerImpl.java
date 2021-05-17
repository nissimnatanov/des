package com.darkevilsudoku.board.impl;

import com.darkevilsudoku.board.SudokuBoard;
import com.darkevilsudoku.board.SudokuBoardSerializer;
import com.darkevilsudoku.board.SudokuFormatter;
import com.darkevilsudoku.board.SudokuValue;

public class SudokuBoardSerializerImpl extends SudokuBoardSerializer {

	private boolean ignoreUnrecognizableChars;
	private boolean ignoreInvalidBoardValues;
	private boolean asciiOnly;
	private boolean compact;

	public SudokuBoardSerializerImpl(
			boolean ignoreUnrecognizableChars, 
			boolean ignoreInvalidBoardValues, 
			boolean asciiOnly,
			boolean compact) {
		this.ignoreUnrecognizableChars = ignoreUnrecognizableChars;
		this.ignoreInvalidBoardValues = ignoreInvalidBoardValues;
		this.asciiOnly = asciiOnly;
		this.compact = compact;
	}
	
	public boolean getIgnoreUnrecognizableChars() {
		return ignoreUnrecognizableChars;
	}
	
	public boolean getIgnoreInvalidBoardValues() {
		return ignoreInvalidBoardValues;
	}
	
	public boolean getCompact() {
		return compact;
	}
	
	private static char getCompactChar(int emptyCount) {
		return(char)('A' + emptyCount - 1);
	}
	
	private static int getCompactSpaceCount(char c) {
		if (c >= 'A' || c <= 'I') {
			return c - 'A' + 1;
		}
		if (c >= 'a' || c <= 'i') {
			return c - 'a' + 1;
		}
		return 0;
	}

	/**
	 * Simple serialization: - empty cell => "0" - value cells: read-only, valid
	 * => "D" where D is the digit that represents the value: 1..9 editable,
	 * valid => "D̲" where D is the same as above followed by combining
	 * underline char. read-only, invalid => "D̷" - combined with solidus
	 * editable, invalid => "D̵" - combined with strikeout
	 */
	@Override
	public void serialize(SudokuBoard board, StringBuilder sb) {
		int emptyCount = 0;
		
		for (int i = 0; i < SudokuBoard.BOARD_SIZE; i++) {
			if (compact) {
				if (!board.hasValue(i)) {
					emptyCount++;
					if (emptyCount == 9) {
						// reached max
						sb.append(getCompactChar(9));
						emptyCount = 0;
					}
					continue;
				}
				else {
					if (emptyCount > 0) {
						// reached value
						sb.append(emptyCount);
						emptyCount = 0;
					}
				}
			}
			
			int options = SudokuFormatter.EMPTY_AS_ZERO;
			if (board.isReadOnlyCell(i)) {
				options |= SudokuFormatter.READONLY;
			}
			if (!ignoreInvalidBoardValues && !board.isValidCell(i)) {
				options |= SudokuFormatter.INVALID;
			}
			
			if (asciiOnly) {
				options |= SudokuFormatter.ASCII_ONLY;
			}
			else {
				options |= SudokuFormatter.FULLWIDTH;
			}
			
			sb.append(SudokuFormatter.toString(board.getValue(i), options));
		}
		
		if (emptyCount > 0) {
			// left spaces at the end
			sb.append(emptyCount);
			emptyCount = 0;
		}
	}

	@Override
	public SudokuBoard deserialize(String s) {
		SudokuBoard board = SudokuBoard.build();
		boolean ignoreSpaces = false;

		if (SudokuFormatter.containsZeros(s)) {

			// spaces will be totally ignored if zeros are present
			ignoreSpaces = true;
		}

		int boardIndex = 0;

		for (int i = 0; i < s.length(); i++) {
			char c = s.charAt(i);

			boolean isValue = SudokuFormatter.isValue(c);
			boolean isSingleEmptyCell = !isValue && SudokuFormatter.isEmpty(c, ignoreSpaces);
			int compactSpaceCount = getCompactSpaceCount(c);

			if (isValue || isSingleEmptyCell) {
				SudokuValue v = SudokuFormatter.getValue(c);
				boolean readOnly = (v != null);

				// got a digit or a valid empty indicator, check if it is
				// read-only or editable
				// for empty, editable mark can be there, yet it is useless
				if ((i + 1) < s.length() && SudokuFormatter.isEditableIndicator(s.charAt(i + 1))) {
					i++;
					readOnly = false;
				}
				if (boardIndex >= SudokuBoard.BOARD_SIZE) {
					throw new IllegalArgumentException("redundant cell character [" + c + "] @ index " + i);
				}

				board.setValue(boardIndex++, v, readOnly);
				continue;
			}

			if (compactSpaceCount > 0) {
				if ((boardIndex + compactSpaceCount) >= SudokuBoard.BOARD_SIZE) {
					throw new IllegalArgumentException("space count overflow [" + c + "] @ index " + i);
				}
				
				// advance, excluding current
				boardIndex += compactSpaceCount - 1;
				continue;
			}
			
			if (Character.isSpaceChar(c)) {
				// ignore other potential spaces chars like new line, tab, and
				// others
				// if zeros are present, spaces will be ignored here as well
				continue;
			}

			if (SudokuFormatter.isValidityIndicator(c)) {
				// silently ignore, the board will control validity indicators
				continue;
			}
			
			if (!ignoreUnrecognizableChars) {
				throw new IllegalArgumentException("unrecognizable character [" + c + "] @ index " + i);
			}
		}

		if (boardIndex != SudokuBoard.BOARD_SIZE) {
			throw new IllegalArgumentException("incomplete board, only " + boardIndex + " cells were read");
		}

		return board;
	}
}
