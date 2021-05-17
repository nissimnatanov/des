using System;
using System.Collections.Generic;
using System.Text;
using System.Diagnostics;

namespace Sudoku
{
    enum SudokuAlgorithmWeight
    {
        None = 0,

        // algorithm levels
        Easy = 1,           // single in square
        Medium = 3,         // single in row/column
        Hard = 10,          // identify pairs
        
        // trial and error
        VeryHard = 30,      // single recursion
        Evil = 500,         // second-level recursion
        DarkEvil =  10000,  // third-level recursion
        BlackHole = 250000  // 4th level of recursion
    }

    enum SudokuClassificationLimits
    {
        // game levels
        EasyGameLimit = 65,
        MediumGameLimit = 100,
        HardGameLimit = 200,
        VeryHardGameLimit = 500,
        EvilGameLimit = 2000,
        DarkEvilGameLimit = 10000,
    }

    interface ISudokuAlgorithm
    {
        SudokuAlgorithmResultDetails Run(SudokuSolver solver);
    }

    interface ISudokuEliminationAlgorithm
    {
        SudokuAlgorithmResultDetails Eliminate(SudokuSolver solver);
    }
    interface ISudokuRecursiveAlgorithm
    {
        SudokuAlgorithmResultDetails Run(SudokuSolver solver);
    }
    #region sudoku solver

    public enum SudokuSolverMode
    {
        Prove,
        Solve,
        Hint
    }

    public class SudokuSolver
    {
        public struct SolutionInfo
        {
            public SolutionInfo(SudokuAlgorithmResultDetails result, SudokuLevel level, int totalWeight, Dictionary<int, int> stats)
            {
                this.Result = result;
                this.GameLevel = level;
                this.TotalWeight = totalWeight;
                this.Statistics = stats;
            }

            public readonly SudokuAlgorithmResultDetails Result;
            public readonly SudokuLevel GameLevel;
            public readonly int TotalWeight;
            public readonly IReadOnlyDictionary<int, int> Statistics;
        }

        static readonly IReadOnlyList<ISudokuAlgorithm> SolvingAlgorithms;
        static readonly IReadOnlyList<ISudokuEliminationAlgorithm> EliminationAlgorithms;
        static readonly ISudokuRecursiveAlgorithm RecursiveAlgorithm;

        static SudokuSolver()
        {
            var SolvingAlgorithmsTemp = new List<ISudokuAlgorithm>();
            SolvingAlgorithmsTemp.Add(new SudokuAlgorithmSingleInSquare());
            SolvingAlgorithmsTemp.Add(new SudokuAlgorithmSingleInRow());
            SolvingAlgorithmsTemp.Add(new SudokuAlgorithmSingleInColumn());
            SolvingAlgorithmsTemp.Add(new SudokuAlgorithmTheOnlyChoiceInCell());
            SolvingAlgorithms = SolvingAlgorithmsTemp;

            var EliminationAlgorithmsTemp = new List<ISudokuEliminationAlgorithm>();
            EliminationAlgorithmsTemp.Add(new SudokuAlgorithmIdentifyPairs());
            EliminationAlgorithms = EliminationAlgorithmsTemp;

            RecursiveAlgorithm = new SudokuAlgorithmTrialAndError();
        }
        
        internal SudokuBoard Board { get; private set; }
        internal SudokuBoard SourceBoard { get; set; }
        internal int CurrentRecursionLevel { get; private set; }
        internal int MaxRecursionLevel { get; private set; }
        internal SudokuSolverMode Mode { get; private set; }
        internal Dictionary<int, int> Statistics { get; private set; }

        private SudokuSolver(SudokuBoard board, SudokuBoard sourceBoard, int currentRecursionLevel, int maxRecursionLevel, SudokuSolverMode mode)
        {
            this.Board = board;
            this.SourceBoard = sourceBoard;
            this.CurrentRecursionLevel = currentRecursionLevel;
            this.MaxRecursionLevel = maxRecursionLevel;
            this.Mode = mode;
        }
        
        internal SolutionInfo Solve()
        {
            if (!this.Board.IsValid)
                return new SudokuSolver.SolutionInfo(
                    SudokuAlgorithmResultDetails.Failed("Invalid cell values"), SudokuLevel.Unclassified, 0, null);

            if (this.Board.FreeCellsCount == 0)
                return new SolutionInfo(SudokuAlgorithmResultDetails.Succeeded(-1, 0, 0), SudokuLevel.Unclassified, 0, null);

            if (this.Board.FreeCellsCount > 64)
                return new SolutionInfo(
                    SudokuAlgorithmResultDetails.Failed("Board has less than 17 values is not solvable"), SudokuLevel.Unclassified, 0, null);

            Statistics = new Dictionary<int, int>();
            bool algorithmFound = true;
            SudokuAlgorithmResultDetails firstSuccessResult = SudokuAlgorithmResultDetails.Unknown();
            int totalWeight = 0;
            
            do {
                algorithmFound = false;
                for (int i = -0; i < SolvingAlgorithms.Count; i++)
                {
                    ISudokuAlgorithm alg = SolvingAlgorithms[i];
                    var res = alg.Run(this);
                    if (res.Result == SudokuAlgorithmResult.Unknown)
                        continue;

                    if (res.Result == SudokuAlgorithmResult.Succeeded)
                    {
                        ReportStats(res.Weight);
                        totalWeight += (int)res.Weight;
                        if (Mode == SudokuSolverMode.Hint)
                        {
                            // one is enough
                            return new SolutionInfo(res, SudokuLevel.Unclassified, totalWeight, this.Statistics);
                        }

                        if (firstSuccessResult.Result == SudokuAlgorithmResult.Unknown)
                            firstSuccessResult = res;

                        algorithmFound = true;
                        break;
                    }
                    else
                    {
                        // failed
                        return new SolutionInfo(res, SudokuLevel.Unclassified, totalWeight, this.Statistics);
                    }
                }

                if (!algorithmFound)
                {
                    // try elimination algorithms
                    for (int i = 0; i < EliminationAlgorithms.Count; i++)
                    {
                        ISudokuEliminationAlgorithm alg = EliminationAlgorithms[i];
                        var res = alg.Eliminate(this);
                        if (res.Result == SudokuAlgorithmResult.Succeeded)
                        {
                            // elimination succeeded
                            ReportStats(res.Weight);
                            totalWeight += (int)res.Weight;
                            algorithmFound = true;
                            break;
                        }
                        else if (res.Result == SudokuAlgorithmResult.Failed)
                        {
                            // no point to continue
                            return new SolutionInfo(SudokuAlgorithmResultDetails.Failed("Could not eliminate values"), SudokuLevel.Unclassified, totalWeight, this.Statistics);
                        }
                        else if (res.Result == SudokuAlgorithmResult.TwoSolutions)
                        {
                            // no point to continue
                            return new SolutionInfo(SudokuAlgorithmResultDetails.TwoSolutions(), SudokuLevel.Unclassified, totalWeight, this.Statistics);
                        }
                    }
                }

                if (!algorithmFound)
                {
                    // try recursive algorithm
                    var res = RecursiveAlgorithm.Run(this);
                    if (res.Result == SudokuAlgorithmResult.Succeeded)
                    {
                        // recursive call succeeded
                        // no success report - algorithm reports stats
                        totalWeight += (int)res.Weight;
                        algorithmFound = true;
                        if (Mode == SudokuSolverMode.Hint && res.ValueIndex >= 0 && res.Value > 0)
                        {
                            // value found
                            return new SolutionInfo(res, SudokuLevel.Unclassified, totalWeight, this.Statistics);
                        }
                        else
                        {
                            continue;
                        }
                    }
                    else if (res.Result == SudokuAlgorithmResult.Failed)
                    {
                        // no point to continue
                        return new SolutionInfo(SudokuAlgorithmResultDetails.Failed("Could not eliminate values"), SudokuLevel.Unclassified, totalWeight, this.Statistics);
                    }
                    else if (res.Result == SudokuAlgorithmResult.TwoSolutions)
                    {
                        // no point to continue
                        return new SolutionInfo(SudokuAlgorithmResultDetails.TwoSolutions(), SudokuLevel.Unclassified, totalWeight, this.Statistics);
                    }
                }

                if (!algorithmFound)
                {
                    return new SolutionInfo(SudokuAlgorithmResultDetails.Unknown(), SudokuLevel.Unclassified, totalWeight, this.Statistics);
                }
            } while (Board.FreeCellsCount > 0) ;

            // define the level of the algorithm
            SudokuLevel level;

            if (totalWeight <= (int)SudokuClassificationLimits.EasyGameLimit)
            {
                if (Board.FreeCellsCount < 55)
                    level = SudokuLevel.Easy;
                else
                    level = SudokuLevel.Medium;
            }
            else if (totalWeight <= (int)SudokuClassificationLimits.MediumGameLimit)
                level = SudokuLevel.Medium;
            else if (totalWeight <= (int)SudokuClassificationLimits.HardGameLimit)
                level = SudokuLevel.Hard;
            else if (totalWeight <= (int)SudokuClassificationLimits.VeryHardGameLimit)
                level = SudokuLevel.VeryHard;
            else if (totalWeight <= (int)SudokuClassificationLimits.EvilGameLimit)
                level = SudokuLevel.Evil;
            else if (totalWeight <= (int)SudokuClassificationLimits.DarkEvilGameLimit)
                level = SudokuLevel.DarkEvil;
            else
                level = SudokuLevel.BlackHole;

            return new SolutionInfo(firstSuccessResult, level, totalWeight, this.Statistics);
        }

        internal void ReportStats(int weight, int count = 1)
        {
            int current = 0;
            if (Statistics.TryGetValue(weight, out current))
            {
                Statistics[weight] = current + count;
            }
            else
            {
                Statistics[weight] = count;
            }
        }

        internal void ReportStats(IReadOnlyDictionary<int, int> stats)
        {
            if (stats != null)
            {
                foreach (var entry in stats)
                {
                    ReportStats(entry.Key, entry.Value);
                }
            }
        }

        public static SolutionInfo RecursiveSolve(SudokuBoard board, SudokuBoard sourceBoard, int currentRecursionLevel, int maxRecursionLevel)
        {
            SudokuSolver recSolver = new SudokuSolver(board, sourceBoard, currentRecursionLevel, maxRecursionLevel, SudokuSolverMode.Solve);

            SolutionInfo info = recSolver.Solve();
            if (board.EditMode && info.Result.Result == SudokuAlgorithmResult.Succeeded)
                board.SetLevel(info.GameLevel, info.TotalWeight);
            return info;
        }

        public static SolutionInfo Solve(SudokuBoard board, SudokuBoard sourceBoard = null, int maxRecursionLevel = 5)
        {
            SudokuSolver solver = new SudokuSolver(board, null, 0, maxRecursionLevel, SudokuSolverMode.Solve);

            SolutionInfo info = solver.Solve();
            if (board.EditMode && info.Result.Result == SudokuAlgorithmResult.Succeeded)
                board.SetLevel(info.GameLevel, info.TotalWeight);
            return info;
        }
        
        public static SolutionInfo Prove(SudokuBoard board, SudokuBoard sourceBoard = null, int maxRecursionLevel = 5)
        {
            SudokuSolver solver = new SudokuSolver(board, sourceBoard, 0, maxRecursionLevel, SudokuSolverMode.Prove);
            SolutionInfo info = solver.Solve();
            if (board.EditMode && info.Result.Result == SudokuAlgorithmResult.Succeeded)
                board.SetLevel(info.GameLevel, info.TotalWeight);
            return info;
        }

        public static SolutionInfo Hint(SudokuBoard board, int maxRecursionLevel = 3)
        {
            SudokuSolver solver = new SudokuSolver(board, null, 0, maxRecursionLevel, SudokuSolverMode.Hint);
            return solver.Solve();
        }
    }

    #endregion
        
    #region single in square

    class SudokuAlgorithmSingleInSquare : ISudokuAlgorithm
    {
        public SudokuAlgorithmSingleInSquare()
        { }

        public SudokuAlgorithmResultDetails Run(SudokuSolver solver)
        {
            var board = solver.Board;
            // check squares
            for (int square = 0; square < 9; square++)
            {
                SudokuValueSet values = board.GetValuesInSquare(square);
                if (values.Count == 9)
                {
                    // full
                    continue;
                }
                foreach (var num in values.EnumMissingValues())
                {
                    int freeCell = -1;
                    for (int cell = 0; cell < 9; cell++)
                    {
                        int index = SudokuBoard.GetIndexFromSquare(square, cell);
                        byte temp = board[index];
                        if (temp == num)
                        {
                            // should not happen since we enumerate missing values here
                            throw new InvalidOperationException();
                        }

                        if (temp != 0)
                            continue;

                        if (board.IsAllowedValue(index, num))
                        {
                            if (freeCell != -1)
                            {
                                // second free cell
                                freeCell = -2;
                                break;
                            }

                            freeCell = cell;
                        }
                    }

                    if (freeCell == -1)
                    {
                        // no free cell in square - fail the solution
                        return SudokuAlgorithmResultDetails.Failed("Single In Square: No free cell in square: {0}", square);
                    }

                    if (freeCell == -2)
                    {
                        // duplicate found
                        continue;
                    }

                    if (freeCell != -1)
                    {
                        int index = SudokuBoard.GetIndexFromSquare(square, freeCell);

                        // found free cell that is single in its square
                        board.SetValue(index, num, false);
                        int weight = (int)SudokuAlgorithmWeight.Easy;
                        return SudokuAlgorithmResultDetails.Succeeded(index, num, weight);
                    }
                }
            }

            return SudokuAlgorithmResultDetails.Unknown();
        }
    }

    #endregion

    #region single in row

    class SudokuAlgorithmSingleInRow : ISudokuAlgorithm
    {
        public SudokuAlgorithmSingleInRow()
        { }

        public SudokuAlgorithmResultDetails Run(SudokuSolver solver)
        {
            var board = solver.Board;
            // check rows
            for (int row = 0; row < 9; row++)
            {
                SudokuValueSet values = board.GetValuesInRow(row);
                if (values.Count == 9)
                {
                    // full
                    continue;
                }

                foreach (var num in values.EnumMissingValues())
                {
                    int freeCol = -1;
                    for (int col = 0; col < 9; col++)
                    {
                        int index = SudokuBoard.GetIndexFromCoordinates(row, col);
                        byte temp = board[index];
                        if (temp == num)
                        {
                            // should not happen since we enumerate missing values here
                            throw new InvalidOperationException();
                        }

                        if (temp != 0)
                            continue;

                        if (board.IsAllowedValue(index, num))
                        {
                            if (freeCol != -1)
                            {
                                // second free cell
                                freeCol = -2;
                                break;
                            }

                            freeCol = col;
                        }
                    }

                    if (freeCol == -1)
                    {
                        // no free cell in row - fail the solution
                        return SudokuAlgorithmResultDetails.Failed("Single In Row: No free cell in row {0}", row);
                    }

                    if (freeCol == -2)
                    {
                        // duplicate found
                        continue;
                    }

                    if (freeCol != -1)
                    {
                        int index = SudokuBoard.GetIndexFromCoordinates(row, freeCol);

                        // found free cell that is single in its row
                        board.SetValue(index, num, false);
                        int weight = (int)SudokuAlgorithmWeight.Easy;
                        return SudokuAlgorithmResultDetails.Succeeded(index, num, weight);
                    }
                }
            }
            return SudokuAlgorithmResultDetails.Unknown();
        }
    }

    #endregion

    #region single in column

    class SudokuAlgorithmSingleInColumn : ISudokuAlgorithm
    {
        public SudokuAlgorithmSingleInColumn()
        { }

        public SudokuAlgorithmResultDetails Run(SudokuSolver solver)
        {
            var board = solver.Board;

            // check columns
            for (int col = 0; col < 9; col++)
            {
                SudokuValueSet values = board.GetValuesInColumn(col);
                if (values.Count == 9)
                {
                    // full
                    continue;
                }

                foreach (var num in values.EnumMissingValues())
                {
                    int freeRow = -1;
                    for (int row = 0; row < 9; row++)
                    {
                        int index = SudokuBoard.GetIndexFromCoordinates(row, col);
                        byte temp = board[index];
                        if (temp == num)
                        {
                            // should not happen since we enumerate missing values here
                            throw new InvalidOperationException();
                        }

                        if (temp != 0)
                            continue;

                        if (board.IsAllowedValue(index, num))
                        {
                            if (freeRow != -1)
                            {
                                // second free cell
                                freeRow = -2;
                                break;
                            }

                            freeRow = row;
                        }
                    }

                    if (freeRow == -1)
                    {
                        // no free cell in column - fail the solution
                        return SudokuAlgorithmResultDetails.Failed("Single In Column: No free cell in column {0}", col);
                    }

                    if (freeRow == -2)
                    {
                        // duplicate found
                        continue;
                    }

                    if (freeRow != -1)
                    {
                        int index = SudokuBoard.GetIndexFromCoordinates(freeRow, col);

                        // found free cell that is single in its column
                        board.SetValue(index, num, false);
                        int weight = (int)SudokuAlgorithmWeight.Easy;
                        return SudokuAlgorithmResultDetails.Succeeded(index, num, weight);
                    }
                }
            }
            return SudokuAlgorithmResultDetails.Unknown();
        }
    }

    #endregion

    #region the only choice in cell

    class SudokuAlgorithmTheOnlyChoiceInCell : ISudokuAlgorithm
    {
        public SudokuAlgorithmTheOnlyChoiceInCell()
        { }

        public SudokuAlgorithmResultDetails Run(SudokuSolver solver)
        {
            var board = solver.Board;
            for (int index = 0; index < 81; index++)
            {
                byte val = board[index];
                if (val != 0)
                    continue; // already set

                var res = CheckCell(solver, index);
                if (res.Result != SudokuAlgorithmResult.Unknown)
                    return res;
            }

            return SudokuAlgorithmResultDetails.Unknown();
        }

        SudokuAlgorithmResultDetails CheckCell(SudokuSolver solver, int index)
        {
            var board = solver.Board;
            var allowedValues = board.GetAllowedValues(index);

            if (allowedValues.Count == 0)
            {
                // no value can be used for that specific cell
                return SudokuAlgorithmResultDetails.Failed("No value can be used for the given cell {0}", index);
            }

            if (allowedValues.Count > 1)
            {
                // more than one solution
                return SudokuAlgorithmResultDetails.Unknown();
            }

            int allowedValue = allowedValues.Combined;

            // exact one value can be used for this cell
            board.SetValue(index, (byte)allowedValue, false);
            int weight = (int)SudokuAlgorithmWeight.Medium;
            return SudokuAlgorithmResultDetails.Succeeded(index, (byte)allowedValue, weight);
        }
    }

    #endregion

    #region identify pairs algorithms

    class SudokuAlgorithmIdentifyPairs : ISudokuEliminationAlgorithm
    {
        public SudokuAlgorithmIdentifyPairs()
        { }

        public SudokuAlgorithmResultDetails Eliminate(SudokuSolver solver)
        {
            SudokuBoard board = solver.Board;
            for (int index = 0; index < 81; index++)
            {
                var allowedValues = board.GetAllowedValues(index);
                if (allowedValues.Count != 2)
                    continue;

                ushort mask = allowedValues.Mask;
                int peerIndex = FindPeerWithSimilarMask(board, mask, SudokuBoard.SudokuRelatedIndicies[index]);
                if(peerIndex < 0)
                {
                    continue;
                }

                // found peer
                int eliminationCount = 0;
                int row = index / 9;
                if (row == peerIndex/9)
                {
                    // same row
                    if (TryEliminate(board, index, peerIndex, mask, SudokuBoard.SudokuRowIndicies[row]))
                        eliminationCount++;
                }

                int col = index % 9;
                if (col == peerIndex % 9)
                {
                    // same column
                    if (TryEliminate(board, index, peerIndex, mask, SudokuBoard.SudokuColIndicies[col]))
                        eliminationCount++;
                }

                int square = SudokuBoard.GetSquareFromIndex(index);
                if (square == SudokuBoard.GetSquareFromIndex(peerIndex))
                {
                    // same square
                    if (TryEliminate(board, index, peerIndex, mask, SudokuBoard.SudokuSquareIndicies[square]))
                        eliminationCount++;
                }

                if (eliminationCount > 0)
                    return SudokuAlgorithmResultDetails.Succeeded(-1, 0, (int)SudokuAlgorithmWeight.Hard);
            }

            return SudokuAlgorithmResultDetails.Unknown();
        }

        bool TryEliminate(SudokuBoard board, int index, int peerIndex, ushort allowedValuesMask, int[] indicies)
        {
            bool found = false;
            for (int temp = 0; temp < 9; temp++)
            {
                int i = indicies[temp];
                if (i == index || i == peerIndex)
                    continue;

                if (board[i] != 0)
                    continue;

                ushort tempAllowedValuesMask = board.GetAllowedValues(i).Mask;
                if ((tempAllowedValuesMask & allowedValuesMask) > 0)
                {
                    // found a cell that we can remove values - turn them off
                    board.DisallowValueMask(i, allowedValuesMask);
                    found = true;
                }
            }

            return found;
        }

        int FindPeerWithSimilarMask(SudokuBoard board, ushort mask, int[] indicies)
        {
            for (int i = 0; i < indicies.Length; i++)
            {
                int peerIndex = indicies[i];
                ushort peerMask = board.GetAllowedValues(peerIndex).Mask;
                if (peerMask == mask)
                    return peerIndex;
            }
            return -1;
        }
    }

    #endregion

    #region trial and error algorithm

    class SudokuAlgorithmTrialAndError : ISudokuRecursiveAlgorithm
    {
        public SudokuAlgorithmTrialAndError()
        { }

        class CompareByAllowedValues : Comparer<int>
        {
            SudokuBoard board_;
            internal CompareByAllowedValues(SudokuBoard board)
            {
                board_ = board;
            }

            public override int Compare(int x, int y)
            {
                return board_.GetAllowedValues(x).Count - board_.GetAllowedValues(y).Count;
            }
        }

        public SudokuAlgorithmResultDetails Run(SudokuSolver solver)
        {
            if (solver.CurrentRecursionLevel >= solver.MaxRecursionLevel)
                return SudokuAlgorithmResultDetails.Unknown(); // this algrithm activates recursion

            int[] sortedIndicies = new int[81];
            int total = 0;
            for (int i = 0; i < sortedIndicies.Length; i++)
            {
                if (solver.Board[i] == 0)
                {
                    sortedIndicies[total++] = i;
                }
            }
            Array.Sort(sortedIndicies, 0, total, new CompareByAllowedValues(solver.Board));

            return TryAndCheckAlgorithm(solver, sortedIndicies, total);
        }

        int RecursionWeight(int currentRecursionLevel)
        {
            int weight;
            switch (currentRecursionLevel)
            {
                case 0:
                    weight = (int)SudokuAlgorithmWeight.VeryHard;
                    break;
                case 1:
                    weight = (int)SudokuAlgorithmWeight.Evil;
                    break;
                case 2:
                    weight = (int)SudokuAlgorithmWeight.DarkEvil;
                    break;
                case 3:
                default:
                    weight = (int)SudokuAlgorithmWeight.BlackHole;
                    for (int i = 4; i <= currentRecursionLevel; i++)
                    {
                        weight *= 3;
                    }
                    break;
            }
            return weight;
        }

        SudokuAlgorithmResultDetails TryAndCheckAlgorithm(SudokuSolver solver, int[] sortedIndicies, int total)
        {
            var board = solver.Board;
            for (int maxRecursionLevel = solver.CurrentRecursionLevel + 1; maxRecursionLevel <= solver.MaxRecursionLevel; maxRecursionLevel++)
            {
                for (int i_ = 0; i_ < total; i_++)
                {
                    int index = sortedIndicies[i_];
                    if (board[index] != 0)
                        continue;

                    ushort allowedValuesMask = board.GetAllowedValues(index).Mask;

                    for (byte testValue = 1; testValue <= 9; testValue++)
                    {
                        ushort testMask = SudokuValueSet.MaskOf(testValue);
                        if ((allowedValuesMask & testMask) == 0)
                        {
                            continue;
                        }

                        // copy the board before testing it
                        SudokuBoard testBoard = new SudokuBoard(solver.Board);
                        testBoard.SetValue(index, testValue, false);
                        var testSolution = SudokuSolver.RecursiveSolve(testBoard, solver.SourceBoard, solver.CurrentRecursionLevel + 1, maxRecursionLevel);

                        if (testSolution.Result.Result == SudokuAlgorithmResult.Unknown)
                        {
                            continue;
                        }

                        if (testSolution.Result.Result == SudokuAlgorithmResult.Failed)
                        {
                            // when settings this value, the solution fails
                            // thus, we can safely remove it
                            board.DisallowValue(index, testValue);

                            int currentWeight = RecursionWeight(solver.CurrentRecursionLevel);
                            solver.ReportStats(currentWeight);

                            solver.ReportStats(testSolution.Statistics);
                            return SudokuAlgorithmResultDetails.Succeeded(-1, 0, currentWeight + testSolution.TotalWeight);
                        }

                        if (testSolution.Result.Result == SudokuAlgorithmResult.TwoSolutions)
                        {
                            // no need to continue from this point if we detect more than one solution is available
                            return SudokuAlgorithmResultDetails.TwoSolutions();
                        }

                        if (testSolution.Result.Result == SudokuAlgorithmResult.Succeeded)
                        {
                            // solution succeeded
                            int currentWeight = RecursionWeight(solver.CurrentRecursionLevel);
                            int totalWeight = currentWeight + testSolution.TotalWeight;
                            if (solver.Mode == SudokuSolverMode.Hint)
                            {
                                // if asked to provide only a hint, set the hint value and stop
                                board.SetValue(index, testValue, false);
                                solver.ReportStats(currentWeight);
                                solver.ReportStats(testSolution.Statistics);
                                return SudokuAlgorithmResultDetails.Succeeded(index, testValue, totalWeight);
                            }

                            if (solver.Mode == SudokuSolverMode.Solve)
                            {
                                // if asked to solve, we can copy the values from the test board and stop
                                solver.ReportStats(currentWeight);
                                solver.ReportStats(testSolution.Statistics);
                                for (int ci = 0; ci < 81; ci++)
                                {
                                    if (board[ci] == 0)
                                    {
                                        board.SetValue(ci, testBoard[ci], false);
                                    }
                                    else if (board[ci] != testBoard[ci])
                                    {
                                        throw new InvalidOperationException();
                                    }
                                }
                                return SudokuAlgorithmResultDetails.Succeeded(index, testValue, totalWeight);
                            }

                            if (solver.Mode == SudokuSolverMode.Prove)
                            {
                                // Prove
                                solver.ReportStats(currentWeight);
                                solver.ReportStats(testSolution.Statistics);
                                if (solver.SourceBoard == null)
                                {
                                    // set the solution found as a possible solver's source board
                                    // if we find next solution which is different from this one, it will cause TwoSolutions failure
                                    solver.SourceBoard = testBoard;
                                }
                                else if (solver.SourceBoard[index] != testValue)
                                {
                                    // got solution that is different from previously found or provided source board => two solutions are available
                                    return SudokuAlgorithmResultDetails.TwoSolutions();
                                }

                                continue;
                            }
                        }
                        
                        throw new InvalidOperationException();
                    }
                }
            }

            return SudokuAlgorithmResultDetails.Unknown();
        }
    }

    #endregion
}
