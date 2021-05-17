using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Sudoku;
using System.IO;
using System.Diagnostics;

namespace SudokuTest
{
    [TestClass]
    public class SudokuFullBoardGeneratorTest
    {
        [TestMethod]
        public void testSpeed()
        {
            generateGames(100);
        }

        void generateGames(int count)
        {
            Random r = new Random();
            for(int i=0; i< count; i++)
            {
                SudokuBoard b = SudokuFullBoardGenerator.GetRandomBoard(r);
                Assert.IsTrue(b.IsValid);
                Assert.AreEqual(b.FreeCellsCount, 0);
            }
        }

        [TestMethod]
        public void solveBoardWith17Values()
        {
            using (var r = new StreamReader(@"F:\Sudoku\CSharp\SudokuTest\games\sudoku17.txt"))
            {
                string line;
                while ((line = r.ReadLine()) != null)
                {
                    if (line.Length == 0)
                        continue;

                    SudokuBoard b = SudokuBoard.LoadFromText(line);
                    var res = SudokuSolver.Prove(b);
                    if (res.Result.Result == SudokuAlgorithmResult.Succeeded)
                    {
                        res = SudokuSolver.Solve(b);
                        if (res.GameLevel >= SudokuLevel.Medium)
                        {
                            string path = @"D:\Old\Projects\Sudoku\games\17\" + b.SolutionWeight + "_" + Guid.NewGuid().ToString() + ".sudoku";
                            b.SaveToFile(path);
                        }
                    }
                    else
                    {
                        Trace.WriteLine("Wrong: " + line);
                    }
                }
            }
        }
    }
}
