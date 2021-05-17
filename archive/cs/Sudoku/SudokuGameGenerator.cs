using System;
using System.Collections.Generic;
using System.Text;
using System.Threading;
using System.Diagnostics;
using System.Windows.Forms;

namespace Sudoku
{
    class SudokuGameGenerator
    {
        struct RemovedValue
        {
            public RemovedValue(int i, byte val)
            {
                Index = i;
                Value = val;
            }
            public int Index;
            public byte Value;
        }

        [Serializable]
        public class ProgressEventArgs : EventArgs
        {
            public ProgressEventArgs(int progress, int gamesGenerated)
            {
                this.Progress = progress;
                this.GamesGenerated = gamesGenerated;
            }

            public readonly int Progress;
            public readonly int GamesGenerated;
        }

        public delegate void ProgressEventHandler(object sender, ProgressEventArgs e);

        bool m_stop;
        SudokuBoard m_generatedBoard;
        SudokuBoard m_sourceBoard = null;
        SudokuLevel m_targetLevel;
        public event EventHandler OnCompleted;
        public event ProgressEventHandler OnProgress;

        int m_gamesRequested;

        int m_gamesInProcess;
        int m_activeThreads;
        Stopwatch m_watch;

        SpinLock m_statsLock = new SpinLock(true);
        int m_gamesGenerated;
        int[] m_totalGames = new int[(int)SudokuLevel.Max];
        long[] m_totalWeights = new long[(int)SudokuLevel.Max];

        public SudokuGameGenerator(SudokuLevel targetLevel, SudokuBoard sourceBoard = null)
        {
            m_targetLevel = targetLevel;
            m_sourceBoard = sourceBoard;
        }

        internal string StatsToString()
        {
            bool lockTaken = false;
            int[] games;
            long[] weights;
            int gamesGenerated;

            try
            {
                m_statsLock.Enter(ref lockTaken);

                games = (int[])m_totalGames.Clone();
                weights = (long[])m_totalWeights.Clone();
                gamesGenerated = m_gamesGenerated;
            }
            finally
            {
                if (lockTaken)
                    m_statsLock.Exit();
            }

            StringBuilder sb = new StringBuilder();
            for (int level = 0; level < games.Length; level++)
            {
                if (games[level] > 0)
                {
                    double percentage = (gamesGenerated == 0) ? 0.0 : ((double)games[level] / gamesGenerated);
                    sb.AppendLine($"{(SudokuLevel)level}: {games[level]} ({percentage:P})");
                }
            }

            sb.AppendLine();

            {
                double percentage = (double)gamesGenerated / m_gamesRequested;
                sb.AppendLine($"Total: {gamesGenerated} ({percentage:P})");
            }

            sb.AppendLine();

            TimeSpan elapsed = m_watch?.Elapsed ?? TimeSpan.Zero;
            if (elapsed > TimeSpan.Zero)
            {
                sb.AppendLine($"Elapsed time: {elapsed.ToString(@"d\.hh\:mm\:ss")}");
                int gamesPerMin = (int)(gamesGenerated / elapsed.TotalMinutes);
                sb.AppendLine($"Games per minute: {gamesPerMin}");

                if (games[(int)SudokuLevel.Evil] > 0)
                {
                    int evilGamesPerMin = (int)(games[(int)SudokuLevel.Evil] / elapsed.TotalMinutes);
                    sb.AppendLine($"Evil games per minute: {evilGamesPerMin}");

                    int averageEvilWeight = (int)(weights[(int)SudokuLevel.Evil] / games[(int)SudokuLevel.Evil]);
                    sb.AppendLine($"Evil average weight: {averageEvilWeight}");
                }

                if (games[(int)SudokuLevel.DarkEvil] > 0)
                {
                    int darkEvilGamesPerHour = (int)(games[(int)SudokuLevel.DarkEvil] / elapsed.TotalHours);
                    sb.AppendLine($"Dark Evil games per hour: {darkEvilGamesPerHour}");

                    int averageDarkEvilWeight = (int)(weights[(int)SudokuLevel.DarkEvil] / games[(int)SudokuLevel.DarkEvil]);
                    sb.AppendLine($"Dark Evil average weight: {averageDarkEvilWeight}");
                }

                if (gamesPerMin > 0)
                {
                    TimeSpan timeLeft = TimeSpan.FromMinutes((double)(m_gamesRequested - gamesGenerated) / gamesPerMin);
                    if (timeLeft > TimeSpan.Zero)
                        sb.AppendLine($"Estimated time left: {timeLeft.ToString(@"d\.hh\:mm\:ss")}");
                }                        
            }

            return sb.ToString();
        }

        public void Stop()
        {
            m_stop = true;
        }

        public bool Done
        {
            get
            {
                return m_generatedBoard != null;
            }
        }

        public SudokuBoard GeneratedBoard
        {
            get
            {
                if (m_generatedBoard == null)
                    throw new InvalidOperationException("No board was generated");
                return m_generatedBoard;
            }
        }

        public void StartSingle()
        {
            Thread t = new Thread(new ThreadStart(GenerateGameThread));
            t.IsBackground = true;
            t.Start();
        }

        public void StartMany(int requestedGames, int threads = -1)
        {
            if (threads == -1)
                threads = Environment.ProcessorCount;

            m_gamesRequested = requestedGames;
            m_gamesGenerated = 0;
            m_gamesInProcess = 0;
            m_activeThreads = threads;
            for (int i = 0; i < m_totalGames.Length; i++)
            {
                m_totalGames[i] = 0;
                m_totalWeights[i] = 0;
            }

            m_watch = Stopwatch.StartNew();
            for (int i = 0; i < threads; i++)
            {
                Thread t = new Thread(new ThreadStart(GenerateGamesThread));
                t.IsBackground = true;
                t.Start();
            }
        }

        void SelectValuesToRemove(SudokuBoard board, int numberOfValuesToRemove, RemovedValue[] removedValues, int[] indicies, ref int remainedIndicies)
        {
            m_generatedBoard = null;

            // select the values to remove and remember them
            for (int i = 0; i < numberOfValuesToRemove && !m_stop; i++)
            {
                int idxToRemove = indicies[remainedIndicies - 1];
                remainedIndicies--;

                if (board[idxToRemove] == 0)
                    throw new InvalidOperationException();

                removedValues[i] = new RemovedValue(idxToRemove, board[idxToRemove]);
                board.SetValue(idxToRemove, 0, isReadOnly: false);
            }
        }

        bool Prove(SudokuBoard board, SudokuBoard sourceBoard)
        {
            if (m_stop)
                return false;

            SudokuSolver.SolutionInfo info = SudokuSolver.Prove(board, sourceBoard);
            return (info.Result.Result == SudokuAlgorithmResult.Succeeded);
        }

        static int s_randomSeedMixer = Process.GetCurrentProcess().Id ^ Thread.CurrentThread.ManagedThreadId;

        Random CreateRandom()
        {
            int randomSeed = Environment.TickCount ^ Interlocked.Increment(ref s_randomSeedMixer);
            return new Random(randomSeed);
        }

        void GenerateGameThread()
        {
            Random r = CreateRandom();
            SudokuBoard board;

            do
            {
                board = GenerateBoardOld2(r, m_sourceBoard);
            } while (!m_stop && board != null && board.GameLevel != m_targetLevel);

            if (board != null)
            {
                board.RestartGame();
                board.EditMode = false;
            }

            // done
            m_generatedBoard = board;

            RaiseOnProgressEvent(100, 1);
            RaiseOnCompletedEvent();
        }

        void GenerateGamesThread()
        {
            Random r = CreateRandom();
            SudokuBoard board;
            int gamesGenerated = 0;
            bool firstThread = false;
            do
            {
                int gamesInProgress = Interlocked.Increment(ref m_gamesInProcess);
                if (gamesInProgress == 1)
                {
                    firstThread = true;
                    RaiseOnProgressEvent(0, 0);
                }

                if (gamesInProgress > m_gamesRequested)
                    break;

                board = GenerateHardestBoard(r, m_sourceBoard);
                if (board == null || m_stop)
                    continue;

                // estimate the level
                SudokuSolver.Solve(board);
                board.RestartGame();

                if (board.GameLevel >= SudokuLevel.Evil)
                {
                    Guid gameGuid = Guid.NewGuid();
                    string gameFile = $"{board.GameLevel}_{board.SolutionWeight}_{gameGuid.ToString("N")}.sudoku";
                    gameFile = System.IO.Path.Combine(@"D:\Old\Projects\Sudoku\games\autosaved", gameFile);
                    board.SaveToFile(gameFile);
                }

                int statIndex = (int)board.GameLevel;

                bool lockTaken = false;
                try
                {
                    m_statsLock.Enter(ref lockTaken);

                    ++m_totalGames[statIndex];
                    m_totalWeights[statIndex] += board.SolutionWeight;
                    gamesGenerated = ++m_gamesGenerated;
                }
                finally
                {
                    if (lockTaken)
                        m_statsLock.Exit();
                }

                if ((gamesGenerated % 10) == 0)
                    RaiseOnProgressEvent((gamesGenerated * 100) / m_gamesRequested, gamesGenerated);
            } while (!m_stop);

            Interlocked.Decrement(ref m_activeThreads);

            if (firstThread)
            {
                // wait for others
                while (m_activeThreads > 0)
                    Thread.Sleep(100);

                m_watch?.Stop();
                RaiseOnCompletedEvent();
            }
        }

        SudokuBoard GenerateBoardOld(Random r, SudokuBoard sourceBoard = null)
        {
            r = r ?? CreateRandom();
            SudokuBoard lastGoodBoard = null;

            RaiseOnSingleGameProgressEvent(0);

            SudokuBoard startingBoard;
            if (sourceBoard == null)
            {
                sourceBoard = SudokuFullBoardGenerator.GetRandomBoard(r);
                startingBoard = sourceBoard;
            }
            else
            {
                // ensure user-provided source board is fully solved
                // if not - solve it
                if (!sourceBoard.IsValid)
                {
                    throw new InvalidOperationException();
                }

                if (sourceBoard.FreeCellsCount > 0)
                {
                    startingBoard = sourceBoard;
                    sourceBoard = new SudokuBoard(sourceBoard);
                    var info = SudokuSolver.Solve(sourceBoard, null);
                    if (info.Result.Result != SudokuAlgorithmResult.Succeeded)
                        throw new InvalidOperationException(info.Result.FailureDetails);
                }
                else
                {
                    startingBoard = sourceBoard;
                }
            }
            SudokuBoard board = new SudokuBoard(true);

            // get and shuffle remained indicies from the board
            int[] indicies = new int[81 - startingBoard.FreeCellsCount];
            int remainedIndicies = 0;
            for (int i = 0; i < 81; i++)
            {
                if (startingBoard[i] != 0)
                {
                    indicies[remainedIndicies++] = i;
                    board.SetValue(i, startingBoard[i], true);
                }
            }
            
            r.Shuffle(indicies);

            // start removing the values randomly
            // remember each one for undo operation
            RemovedValue[] removedValues = new RemovedValue[16];
            Guid gameGuid = Guid.Empty;
            string gameFile = null;

            do
            {
                int maxToRemove;
                if (remainedIndicies > 80)
                {
                    // just started
                    maxToRemove = 12;
                }
                else if (remainedIndicies > 74)
                {
                    // in progress
                    maxToRemove = 6;
                }
                else if (remainedIndicies > 66)
                {
                    // in progress
                    maxToRemove = 2;
                }
                else
                {
                    // about the end
                    maxToRemove = 1;
                }
                int numberOfValuesToRemove = 1;
                if (maxToRemove > 1)
                {
                    int rTemp = maxToRemove / 2;
                    numberOfValuesToRemove = (maxToRemove + 1 - rTemp) + r.Next(rTemp);
                }

                numberOfValuesToRemove = Math.Min(numberOfValuesToRemove, remainedIndicies);

#if INTEGRITYCHECK
                // before removing, assert proved
                Debug.Assert(Prove(board, sourceBoard));
                board.RestartGame();
#endif

                SelectValuesToRemove(board, numberOfValuesToRemove, removedValues, indicies, ref remainedIndicies);

                if (m_stop)
                    break;

                // now, ensure the board created has a valid solution
                // if not, revert the changes using binary-search like algorithm
                bool proved = Prove(board, sourceBoard);
                board.RestartGame();

                if (m_stop)
                    break;
                
                int estimateProgress = (int)(1.3 * (double)(81 - remainedIndicies));
                if (estimateProgress >= 100)
                    estimateProgress = 99;
                RaiseOnSingleGameProgressEvent(estimateProgress);

                if (!proved)
                {
                    for (int i = 0; i < numberOfValuesToRemove; i++)
                    {
                        board.SetValue(removedValues[i].Index, removedValues[i].Value, true);

                        proved = Prove(board, sourceBoard);
                        board.RestartGame();
                        if (proved)
                            break;
                    }
                }

                if (!proved)
                {
                    if (remainedIndicies <= 0)
                    {
                        // TODO smart backtracking, for now - stop forcibly
                        //Trace.WriteLine("No more indicies to remove, max level reached: " + board.GameLevel);
                        break;
                    }

                    // none removed, continue with different set of random values
                    // we reduced the number of remained indicies, so it will stop!
                    continue;
                }

                // estimate the level
                Stopwatch watch = Stopwatch.StartNew();

                // re-estimate the game using Solve
                var solution = SudokuSolver.Solve(board);
                board.RestartGame();

                watch.Stop();
                if (watch.Elapsed >= TimeSpan.FromMilliseconds(50))
                {
                    if (gameGuid == Guid.Empty)
                        gameGuid = Guid.NewGuid();

                    string fileName = $"{board.GameLevel}_{gameGuid.ToString("N")}.sudoku";
                    board.SaveToFile(System.IO.Path.Combine(@"D:\Old\Projects\Sudoku\games\autosaved\long", fileName));
                }
                
                if (board.GameLevel >= SudokuLevel.Evil)
                {
                    if (gameGuid == Guid.Empty)
                        gameGuid = Guid.NewGuid();
                    if (gameFile != null)
                        System.IO.File.Delete(gameFile);

                    // auto-save as long
                    gameFile = $"{board.GameLevel}_{solution.TotalWeight}_{gameGuid.ToString("N")}.sudoku";
                    gameFile = System.IO.Path.Combine(@"D:\Old\Projects\Sudoku\games\autosaved", gameFile);
                    board.SaveToFile(gameFile);
                }

                if (board.GameLevel < m_targetLevel)
                {
                    // the board did not reach the target level
                    if (remainedIndicies > 0)
                    {
                        continue;
                    }
                    else
                    {
                        //Trace.WriteLine("No more indicies to remove, max level reached: " + board.GameLevel);
                        break;
                    }
                }
                else if (board.GameLevel > m_targetLevel)
                {
                    // no board was generated, break with bigger level
                    // if lastGoodBoard does not exist, current board is returned and the caller will restart!
                    //Trace.WriteLine("Level reached: " + board.GameLevel + " while asked for " + m_targetLevel);
                    break;
                }

                // the generated game has the correct level
                if (remainedIndicies == 0)
                {
                    //Trace.WriteLine("No more left to remove, target level reached: " + board.GameLevel);
                    break;
                }

                double chanceToContinue;
                if (m_targetLevel == SudokuLevel.Easy)
                {
                    // in easy level make sure at least 40 cells where cleared
                    const int MinimumCleared = 40;
                    if (remainedIndicies >= (81 - MinimumCleared))
                        continue;

                    chanceToContinue = 0.95;
                }
                else if (m_targetLevel < SudokuLevel.Hard)
                {
                    chanceToContinue = 0.9;
                }
                else if (m_targetLevel < SudokuLevel.Evil)
                {
                    chanceToContinue = 0.85;
                }
                else
                {
                    chanceToContinue = 1; // for evil and dark-evil - keep making it even more evil
                }

                if (r.NextDouble() > chanceToContinue)
                {
                    //Trace.WriteLine("Decided to stop now, target level reached: " + board.GameLevel);
                    break;
                }

                // decided to continue despite the fact that the game is at the target level
                // remember it - it will be used if the generated level will 'overflow'
                lastGoodBoard = new SudokuBoard(board);

            } while (!m_stop);

            if ((board == null || board.GameLevel != m_targetLevel) && (lastGoodBoard != null))
            {
                board = lastGoodBoard;
            }

            return board;
        }

        SudokuBoard GenerateBoardOld2(Random r, SudokuBoard sourceBoard = null)
        {
            r = r ?? CreateRandom();
            SudokuBoard lastGoodBoard = null;

            if (m_gamesRequested <= 1)
                RaiseOnSingleGameProgressEvent(0);

            SudokuBoard startingBoard;
            if (sourceBoard == null)
            {
                sourceBoard = SudokuFullBoardGenerator.GetRandomBoard(r);
                startingBoard = sourceBoard;
            }
            else
            {
                // ensure user-provided source board is fully solved
                // if not - solve it
                if (!sourceBoard.IsValid)
                {
                    throw new InvalidOperationException();
                }

                if (sourceBoard.FreeCellsCount > 0)
                {
                    startingBoard = sourceBoard;
                    sourceBoard = new SudokuBoard(sourceBoard);
                    var info = SudokuSolver.Solve(sourceBoard, null);
                    if (info.Result.Result != SudokuAlgorithmResult.Succeeded)
                        throw new InvalidOperationException(info.Result.FailureDetails);
                }
                else
                {
                    startingBoard = sourceBoard;
                }
            }

            SudokuBoard board = new SudokuBoard(true);

            // get and shuffle remained indicies from the board
            int[] indicies = new int[81 - startingBoard.FreeCellsCount];
            int remainedIndicies = 0;
            for (int i = 0; i < 81; i++)
            {
                if (startingBoard[i] != 0)
                {
                    indicies[remainedIndicies++] = i;
                    board.SetValue(i, startingBoard[i], true);
                }
            }

            r.Shuffle(indicies);

            // start removing the values randomly
            // remember each one for undo operation
            RemovedValue[] removedValues = new RemovedValue[3];
            Guid gameGuid = Guid.Empty;
            string gameFile = null;
            int highestScore = 0;
            // remove first set of random cells
            while (board.FreeCellsCount < 33)
            {
                const int numberOfValuesToRemove = 3;
                SelectValuesToRemove(board, numberOfValuesToRemove, removedValues, indicies, ref remainedIndicies);

                if (m_stop)
                    break;

                // now, ensure the board created has a valid solution
                // if not, revert the changes using binary-search like algorithm
                var proveRes = SudokuSolver.Prove(board, sourceBoard);
                board.RestartGame();
                bool proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;

                if (m_stop)
                    break;

                if (m_gamesRequested <= 1)
                {
                    int estimateProgress = (int)(1.3 * (double)(81 - remainedIndicies));
                    if (estimateProgress >= 100)
                        estimateProgress = 99;
                    RaiseOnSingleGameProgressEvent(estimateProgress);
                }

                if (!proved)
                {
                    for (int i = 0; i < numberOfValuesToRemove; i++)
                    {
                        board.SetValue(removedValues[i].Index, removedValues[i].Value, true);

                        proveRes = SudokuSolver.Prove(board, sourceBoard);
                        proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;
                        board.RestartGame();
                        if (proved)
                            break;
                    }
                }

                if (proved)
                {
                    highestScore = proveRes.TotalWeight;
                }
            }

            // start the hardcore algorithm to remove the cell that leads to highest raise in the score
            List<int> selectedIndices = new List<int>();
            do
            {
                int newHighestScore = -1;
                selectedIndices.Clear();

                for (int i = 0; i < remainedIndicies; i++)
                {
                    if (m_stop)
                        break;

                    int indexToRemove = indicies[i];
                    byte valueToRemove = board[indexToRemove];
                    if (valueToRemove == 0)
                        throw new InvalidOperationException();

                    board.SetValue(indexToRemove, 0, false);

                    var proveRes = SudokuSolver.Prove(board, sourceBoard);
                    board.RestartGame();

                    board.SetValue(indexToRemove, valueToRemove, true);

                    bool proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;

                    if (proved)
                    {
                        if (proveRes.TotalWeight > newHighestScore)
                        {
                            newHighestScore = proveRes.TotalWeight;
                            selectedIndices.Clear();
                            selectedIndices.Add(i);
                        }
                        else if (proveRes.TotalWeight == newHighestScore)
                        {
                            selectedIndices.Add(i);
                        }
                    }
                    else
                    {
                        --remainedIndicies;
                        if (i != remainedIndicies)
                        {
                            // swap with last and stay on same value
                            indicies[i] = indicies[remainedIndicies];
                            indicies[remainedIndicies] = indexToRemove;
                            i--; // stay on the same index
                        }
                    }
                }

                if (m_stop)
                    break;

                if (selectedIndices.Count == 0)
                {
                    // found no indices, break
                    //Trace.WriteLine("No more indicies to remove, max level reached: " + board.GameLevel);
                    break;
                }

                // pick one of the best indicies and swap it with last in indicices list
                int selectedI = selectedIndices[r.Next(selectedIndices.Count)];
                int selectedIndex = indicies[selectedI];
                if (selectedI != (remainedIndicies - 1))
                {
                    indicies[selectedI] = indicies[remainedIndicies - 1];
                    indicies[remainedIndicies - 1] = selectedIndex;
                }

                remainedIndicies--;
                board.SetValue(selectedIndex, 0, false);

                if (m_gamesRequested <= 1)
                {
                    int estimateProgress = (int)(1.3 * (double)(81 - remainedIndicies));
                    if (estimateProgress >= 100)
                        estimateProgress = 99;
                    RaiseOnSingleGameProgressEvent(estimateProgress);
                }

                // estimate the level
                Stopwatch watch = Stopwatch.StartNew();

                // re-estimate the game using Solve
                var solution = SudokuSolver.Solve(board);
                board.RestartGame();

                watch.Stop();
                if (watch.Elapsed >= TimeSpan.FromMilliseconds(50))
                {
                    if (gameGuid == Guid.Empty)
                        gameGuid = Guid.NewGuid();

                    string fileName = $"{board.GameLevel}_{gameGuid.ToString("N")}.sudoku";
                    board.SaveToFile(System.IO.Path.Combine(@"D:\Old\Projects\Sudoku\games\autosaved\long", fileName));
                }

                if (board.GameLevel >= SudokuLevel.Evil)
                {
                    if (gameGuid == Guid.Empty)
                        gameGuid = Guid.NewGuid();
                    if (gameFile != null)
                        System.IO.File.Delete(gameFile);

                    // auto-save as long
                    gameFile = $"{board.GameLevel}_{solution.TotalWeight}_{gameGuid.ToString("N")}.sudoku";
                    gameFile = System.IO.Path.Combine(@"D:\Old\Projects\Sudoku\games\autosaved", gameFile);
                    board.SaveToFile(gameFile);
                }

                if (board.GameLevel < m_targetLevel)
                {
                    // the board did not reach the target level
                    if (remainedIndicies > 0)
                    {
                        continue;
                    }
                    else
                    {
                        //Trace.WriteLine("No more indicies to remove, max level reached: " + board.GameLevel);
                        break;
                    }
                }
                else if (board.GameLevel > m_targetLevel)
                {
                    // no board was generated, break with bigger level
                    // if lastGoodBoard does not exist, current board is returned and the caller will restart!
                    //Trace.WriteLine("Level reached: " + board.GameLevel + " while asked for " + m_targetLevel);
                    break;
                }

                // the generated game has the correct level
                if (remainedIndicies == 0)
                {
                    //Trace.WriteLine("No more left to remove, target level reached: " + board.GameLevel);
                    break;
                }

                double chanceToContinue;
                if (m_targetLevel == SudokuLevel.Easy)
                {
                    // in easy level make sure at least 40 cells where cleared
                    const int MinimumCleared = 40;
                    if (remainedIndicies >= (81 - MinimumCleared))
                        continue;

                    chanceToContinue = 0.95;
                }
                else if (m_targetLevel < SudokuLevel.Hard)
                {
                    chanceToContinue = 0.9;
                }
                else if (m_targetLevel < SudokuLevel.Evil)
                {
                    chanceToContinue = 0.85;
                }
                else
                {
                    chanceToContinue = 1; // for evil and dark-evil - keep making it even more evil
                }

                if (r.NextDouble() > chanceToContinue)
                {
                    //Trace.WriteLine("Decided to stop now, target level reached: " + board.GameLevel);
                    break;
                }

                // decided to continue despite the fact that the game is at the target level
                // remember it - it will be used if the generated level will 'overflow'
                lastGoodBoard = new SudokuBoard(board);

            } while (!m_stop);

            if ((board == null || board.GameLevel != m_targetLevel) && (lastGoodBoard != null))
            {
                board = lastGoodBoard;
            }

            return board;
        }

        void RemoveFirstSet(SudokuBoard board, SudokuBoard sourceBoard, int targetFreeCellCount, int[] indicies, ref int remainedIndicies)
        {
            RemovedValue[] removedValues = new RemovedValue[3];

            while (board.FreeCellsCount < targetFreeCellCount && remainedIndicies > 0)
            {
                const int numberOfValuesToRemove = 3;
                SelectValuesToRemove(board, Math.Min(numberOfValuesToRemove, remainedIndicies), removedValues, indicies, ref remainedIndicies);

                if (m_stop)
                    break;

                // now, ensure the board created has a valid solution
                // if not, revert the changes using binary-search like algorithm
                var proveRes = SudokuSolver.Prove(board, sourceBoard);
                board.RestartGame();
                bool proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;

                if (m_stop)
                    break;

                if (m_gamesRequested <= 1)
                {
                    int estimateProgress = (int)(1.3 * (double)(81 - remainedIndicies));
                    if (estimateProgress >= 100)
                        estimateProgress = 99;
                    RaiseOnSingleGameProgressEvent(estimateProgress);
                }

                if (!proved)
                {
                    for (int i = 0; i < numberOfValuesToRemove; i++)
                    {
                        board.SetValue(removedValues[i].Index, removedValues[i].Value, true);

                        proveRes = SudokuSolver.Prove(board, sourceBoard);
                        proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;
                        board.RestartGame();
                        if (proved)
                            break;
                    }
                }
            }
        }

        SudokuBoard InitializeBoardsAndIndicies(Random r, ref SudokuBoard sourceBoard, out int[] indicies)
        {
            SudokuBoard startingBoard;
            if (sourceBoard == null)
            {
                sourceBoard = SudokuFullBoardGenerator.GetRandomBoard(r);
                startingBoard = sourceBoard;
            }
            else
            {
                // ensure user-provided source board is fully solved
                // if not - solve it
                if (!sourceBoard.IsValid)
                {
                    throw new InvalidOperationException();
                }

                if (sourceBoard.FreeCellsCount > 0)
                {
                    startingBoard = sourceBoard;
                    sourceBoard = new SudokuBoard(sourceBoard);
                    var info = SudokuSolver.Solve(sourceBoard, null);
                    if (info.Result.Result != SudokuAlgorithmResult.Succeeded)
                        throw new InvalidOperationException(info.Result.FailureDetails);
                }
                else
                {
                    startingBoard = sourceBoard;
                }
            }

            SudokuBoard board = new SudokuBoard(true);

            // get and shuffle remained indicies from the board
            indicies = new int[81 - startingBoard.FreeCellsCount];
            int remainedIndicies = 0;
            for (int i = 0; i < 81; i++)
            {
                if (startingBoard[i] != 0)
                {
                    indicies[remainedIndicies++] = i;
                    board.SetValue(i, startingBoard[i], true);
                }
            }

            r.Shuffle(indicies);

            return board;
        }

        void SelectBestIndexRec(SudokuBoard board, SudokuBoard sourceBoard, int[] indicies, int start, ref int remainedIndicies, int recursionLevel, int[] selectedIndicies, int[] bestBoardIndicies, ref int highestScore, int topN)
        {
            int[] topIndicies = null;
            int[] topIndiciesWeights = null ;
            byte[] topIndiciesValues = null;
            bool useTopNRec = (topN > 1) && ((recursionLevel + 1) < selectedIndicies.Length);
            if (useTopNRec)
            {
                topIndicies = new int[topN];
                for (int j = 0; j < topIndicies.Length; j++)
                    topIndicies[j] = -1;
                
                topIndiciesWeights = new int[topIndicies.Length];
                topIndiciesValues = new byte[topIndicies.Length];
            }

            for (int i = start; i < remainedIndicies; i++)
            {
                if (m_stop)
                    return;

                bool found = false;
                for (int j = 0; j < recursionLevel; j++)
                {
                    if (selectedIndicies[j] == i)
                    {
                        found = true;
                        break;
                    }
                }
                if (found)
                    continue;

                int indexToRemove = indicies[i];
                byte valueToRemove = board[indexToRemove];
                if (valueToRemove == 0)
                    throw new InvalidOperationException();

                board.SetValue(indexToRemove, 0, false);
                selectedIndicies[recursionLevel] = i;

                // top level of the recursion - check if selected values are good
                var proveRes = SudokuSolver.Prove(board, sourceBoard);
                board.RestartGame();

                board.SetValue(indexToRemove, valueToRemove, true);

                bool proved = proveRes.Result.Result == SudokuAlgorithmResult.Succeeded;
                if (!proved)
                {
                    if (recursionLevel == 0)
                    {
                        // we are at the top of the tree, current value was the first removed and it is not good
                        // remove the value by swapping it with the last (or ignore if last already)
                        --remainedIndicies;
                        if (i != remainedIndicies)
                        {
                            indicies[i] = indicies[remainedIndicies];
                            indicies[remainedIndicies] = indexToRemove;
                            i--; // stay on the same index
                        }
                    }
                    // else - we are in the middle of the tree, the value might be good if it comes without the first one
                    continue;
                }

                // did we get a better score?
                if (proveRes.TotalWeight > highestScore)
                {
                    // new high score
                    highestScore = proveRes.TotalWeight;
                    for (int j = 0; j <= recursionLevel; j++)
                    {
                        // remember the board indicies, not the shuffled ones
                        // later on we will remove those from the remained indicies
                        bestBoardIndicies[j] = indicies[selectedIndicies[j]];
                    }

                    // reset best indicies tail
                    for (int j = recursionLevel + 1; j < bestBoardIndicies.Length; j++)
                    {
                        bestBoardIndicies[j] = -1;
                    }
                }

                if (useTopNRec)
                {
                    int updateTop = -1;
                    for (int j = 0; j < topIndicies.Length; j++)
                    {
                        if (topIndiciesWeights[j] < proveRes.TotalWeight)
                        {
                            updateTop = j;
                            break;
                        }
                    }

                    if (updateTop >= 0)
                    {
                        if ((updateTop + 1) < topIndicies.Length)
                        {
                            int toCopy = (topIndicies.Length - updateTop - 1);
                            Array.Copy(topIndicies, updateTop, topIndicies, updateTop + 1, toCopy);
                            Array.Copy(topIndiciesWeights, updateTop, topIndiciesWeights, updateTop + 1, toCopy);
                            Array.Copy(topIndiciesValues, updateTop, topIndiciesValues, updateTop + 1, toCopy);
                        }

                        topIndicies[updateTop] = i;
                        topIndiciesWeights[updateTop] = highestScore;
                        topIndiciesValues[updateTop] = valueToRemove;
                    }
                }
            }
            
            if (useTopNRec)
            {
                // we are still in the middle of the tree, continue to the bottom for the TOPN values
                for (int j = 0; j < topIndicies.Length && topIndicies[j] != -1; j++)
                {
                    int i = topIndicies[j];
                    int indexToRemove = indicies[i];
                    byte removedValue = topIndiciesValues[j];

                    board.SetValue(indexToRemove, 0, false);
                    selectedIndicies[recursionLevel] = i;

                    SelectBestIndexRec(board, sourceBoard, indicies, i + 1, ref remainedIndicies, recursionLevel + 1, selectedIndicies, bestBoardIndicies, ref highestScore, topN);
                    board.SetValue(indexToRemove, removedValue, true);
                }
            }
        }
        class RemainedIndiciesComparer : Comparer<int>
        {
            SudokuBoard board_;

            public RemainedIndiciesComparer(SudokuBoard board)
            {
                board_ = board;
            }

            public override int Compare(int x, int y)
            {
                byte xValue = board_[x];
                byte yValue = board_[y];
                
                int xCount = board_.GetValueCount(xValue);
                int yCount = board_.GetValueCount(yValue);

                if (xCount != yCount)
                {
                    // put more frequent in front
                    return yCount - xCount;
                }

                // stats around the cell itself
                int xRelatedCount = board_.GetValuesInRow(x / 9).Count + board_.GetValuesInColumn(x % 9).Count + board_.GetValuesInSquare(SudokuBoard.GetSquareFromIndex(x)).Count;
                int yRelatedCount = board_.GetValuesInRow(y / 9).Count + board_.GetValuesInColumn(y % 9).Count + board_.GetValuesInSquare(SudokuBoard.GetSquareFromIndex(y)).Count;
                return yRelatedCount - xRelatedCount;
            }
        }

        SudokuBoard GenerateHardestBoard(Random r, SudokuBoard sourceBoard = null)
        {
            r = r ?? CreateRandom();
            
            int[] indicies;
            SudokuBoard board = InitializeBoardsAndIndicies(r, ref sourceBoard, out indicies);
            
            int remainedIndicies = indicies.Length;

            // remove first set of random cells
            RemoveFirstSet(board, sourceBoard, 36, indicies, ref remainedIndicies);

            // start the hardcore algorithm to remove the cell that leads to highest raise in the score
            int highestScore = -1;
            int nextResortAfterRemained = 81;

            const int maxRecursionSize = 1;
            const int topN = 1;
            const int resortAfterInterval = 3;

            int[] selectedIndices = new int[maxRecursionSize];
            int[] bestBoardIndicies = new int[maxRecursionSize];
            do
            {
                if (nextResortAfterRemained >= remainedIndicies)
                {
                    // resort indicies
                    Array.Sort(indicies, 0, remainedIndicies, new RemainedIndiciesComparer(board));

                    nextResortAfterRemained = remainedIndicies - resortAfterInterval;
                }

                highestScore = -1;

                for (int i=0; i<maxRecursionSize; i++)
                {
                    selectedIndices[i] = -1;
                    bestBoardIndicies[i] = -1;
                }

                SelectBestIndexRec(board, sourceBoard, indicies, 0, ref remainedIndicies, 0, selectedIndices, bestBoardIndicies, ref highestScore, topN);
                if (m_stop)
                {
                    break;
                }

                // check if any selected
                if (bestBoardIndicies[0] == -1)
                {
                    // no better solution found, can safely break now
                    if (remainedIndicies != 0)
                        throw new InvalidOperationException();
                    break;
                }

                // found good indicies
                for (int i = 0; i < bestBoardIndicies.Length; i++)
                {
                    if (bestBoardIndicies[i] == -1)
                        break;

                    board.SetValue(bestBoardIndicies[i], 0, false);

                    // remove from indicies
                    --remainedIndicies;
                    for (int j = 0; j < remainedIndicies; j++)
                    {
                        if (indicies[j] == bestBoardIndicies[i])
                        {
                            indicies[j] = indicies[remainedIndicies];
                            indicies[remainedIndicies] = bestBoardIndicies[i];
                            break;
                        }
                    }
                }
#if DEBUG
                var proveRes = SudokuSolver.Prove(board, sourceBoard);
                board.RestartGame();
                if (proveRes.Result.Result != SudokuAlgorithmResult.Succeeded)
                {
                    throw new InvalidOperationException();
                }

                if (proveRes.TotalWeight != highestScore)
                {
                    throw new InvalidOperationException();
                }
#endif
            } while (!m_stop && remainedIndicies > 0);

            return board;
        }

        void RaiseOnCompletedEvent()
        {
            OnCompleted?.Invoke(this, EventArgs.Empty);
        }

        void RaiseOnSingleGameProgressEvent(int progress)
        {
            RaiseOnProgressEvent(progress, 0);
        }

        void RaiseOnProgressEvent(int progress, int gamesGenerated)
        {
            OnProgress?.Invoke(this, new ProgressEventArgs(progress, gamesGenerated));
        }

        /// <summary>
        /// creates a random game
        /// </summary>
        public static SudokuBoard CreateRandomGame(IWin32Window parent, SudokuLevel level, SudokuBoard sourceBoard = null)
        {
            SudokuGameGenerator generator = new SudokuGameGenerator(level, sourceBoard);
            SudokuBoard board;

            WaitForm waitDialog = new WaitForm("Please wait while the game is generated...");
            generator.OnProgress += new ProgressEventHandler(waitDialog.ProgressEventHandler);
            generator.OnCompleted += new EventHandler(waitDialog.StopEventHandler);

            generator.StartSingle();
            DialogResult res = waitDialog.ShowDialog(parent);
            generator.Stop();
            if (res != DialogResult.OK)
            {
                return null;
            }

            try
            {
                board = generator.GeneratedBoard;
            }
            catch
            {
                board = null;
            }

            return board;
        }
    }
}
