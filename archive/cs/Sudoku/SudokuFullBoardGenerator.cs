using System;
using System.Collections.Generic;
using System.Text;
using System.IO;

namespace Sudoku
{
    /// <summary>
    /// Holds list of full Sudoku boards used as a initial template to the game generator
    /// </summary>
    public class SudokuFullBoardGenerator
    {
        private static void SetSquareValues(SudokuBoard board, int square, int cellIndex, byte[] values, int startValues, int count, bool isReadOnly)
        {
            while (count-- > 0)
            {
                board.SetValue(SudokuBoard.GetIndexFromSquare(square, cellIndex), values[startValues++], isReadOnly);
                cellIndex++;
            }

            if (!board.IsValid)
                throw new InvalidOperationException();
        }

        private static bool TryFillSquare(Random r, SudokuBoard board, int square)
        {
            int[] cellIndices = new int[9];
            for (int cell = 0; cell < 9; cell++)
                cellIndices[cell] = cell;
            
            r.Shuffle(cellIndices);

            int allowedRetries = 1;
            bool retry;
            do
            {
                retry = false;
                for (int i = 0; i < 9; i++)
                {
                    int cell = cellIndices[i];
                    int index = SudokuBoard.GetIndexFromSquare(square, cell);
                    SudokuValueSet allowedValues = board.GetAllowedValues(index);
                    bool valid = (allowedValues.Count > 0);
                    if (valid)
                    {
                        int next = r.Next(allowedValues.Count);
                        byte value = allowedValues[next];
                        board.SetValue(index, value, false);
                        if (!board.IsValid)
                            valid = false;
                    }
                    
                    if (!valid)
                    {
                        if (allowedRetries-- == 0)
                            return false;

                        // reset the square
                        for (cell = 0; cell < 9; cell++)
                        {
                            board.SetValue(SudokuBoard.GetIndexFromSquare(square, cell), 0, false);
                        }
                        retry = true;
                        break;
                    }
                }
            } while (retry);

            return true;
        }

        public static SudokuBoard GetRandomBoard(Random r)
        {
            r = r ?? new Random();

            byte[] values = new byte[9];
            for (byte v = 1; v <= 9; v++)
                values[v - 1] = v;

            SudokuBoard board = new SudokuBoard(true);

            // populate squares 0, 4 (middle) and 8 (last)
            // since there is no intersection between those, safe to use totally random values
            r.Shuffle(values);
            SetSquareValues(board, 0, 0, values, 0, 9, true);
            r.Shuffle(values);
            SetSquareValues(board, 4, 0, values, 0, 9, true);
            r.Shuffle(values);
            SetSquareValues(board, 8, 0, values, 0, 9, true);

            do
            {
                board.RestartGame();

                // fill squares 1 to 7, skipping 4
                for (int square = 1; square <= 7; square++)
                {
                    if (square == 4)
                        continue; // skip

                    if (!TryFillSquare(r, board, square))
                    {
                        break;
                    }
                }
            } while (board.FreeCellsCount > 0);

            for (int i = 0; i < 81; i++)
            {
                board.SetValue(i, board[i], true);
            }

            return board;
        }
    }
}
