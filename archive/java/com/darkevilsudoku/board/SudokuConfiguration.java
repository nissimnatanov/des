package com.darkevilsudoku.board;

public final class SudokuConfiguration {
	private static boolean enableBoardIntegrityCheck = false;

	public static void setEnableBoardIntegrityCheck(boolean value) {
		enableBoardIntegrityCheck = value;
	}

	public static boolean getEnableBoardIntegrityCheck() {
		return enableBoardIntegrityCheck;
	}

}
