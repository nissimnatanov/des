package com.darkevilsudoku.board;

import com.darkevilsudoku.board.impl.SudokuBoardSerializerImpl;

public abstract class SudokuBoardSerializer {

	public abstract void serialize(SudokuBoard board, StringBuilder sb);
	public abstract SudokuBoard deserialize(String s);
	
	public String serialize(SudokuBoard board) {
		StringBuilder sb = new StringBuilder();
		serialize(board, sb);
		return sb.toString();
	}
	
	public static final class Builder {
		private boolean ignoreUnrecognizableChars = false;
		private boolean ignoreInvalidBoardValues = false;
		private boolean asciiOnly = false;
		private boolean compact = false;
		
		private Builder() {
		}
		
		public Builder setIgnoreUnrecognizableChars(boolean ignoreUnrecognizableChars) {
			this.ignoreUnrecognizableChars = ignoreUnrecognizableChars;
			return this;
		}
		
		public Builder setIgnoreInvalidBoardValues(boolean ignoreInvalidBoardValues) {
			this.ignoreInvalidBoardValues = ignoreInvalidBoardValues;
			return this;
		}

		public Builder setAsciiOnly(boolean asciiOnly) {
			this.asciiOnly = asciiOnly;
			return this;
		}
		
		public Builder setCompact(boolean compact) {
			this.compact = compact;
			return this;
		}
		
		public SudokuBoardSerializer build() {
			return new SudokuBoardSerializerImpl(ignoreUnrecognizableChars, ignoreInvalidBoardValues, asciiOnly, compact);
		}
	}
	
	public static Builder builder() {
		return new Builder();
	}

}
