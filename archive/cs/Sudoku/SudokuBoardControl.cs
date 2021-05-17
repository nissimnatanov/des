using System;
using System.Diagnostics;
using System.Collections.Generic;
using System.ComponentModel;
using System.Drawing;
using System.Data;
using System.Text;
using System.Windows.Forms;
using System.Threading;

namespace Sudoku
{
    public partial class SudokuBoardControl : UserControl
    {
        #region SquareInfo
        class SquareInfo
        {
            TableLayoutPanel m_square;

            public SquareInfo(TableLayoutPanel square, int squareIndex)
            {
                SquareIndex = squareIndex;
                m_square = square;
            }

            public readonly Rectangle[] Cells = new Rectangle[9];
            public readonly int SquareIndex;

            public int GetCellFromScreen(int screenX, int screenY)
            {
                for (int i = 0; i < Cells.Length; i++)
                {
                    Rectangle cellScreen = m_square.RectangleToScreen(Cells[i]);
                    if (cellScreen.Contains(screenX, screenY))
                    {
                        return i;
                    }
                }

                return -1;
            }

            public int GetCellFromTable(int x, int y)
            {
                for (int i = 0; i < Cells.Length; i++)
                {
                    Rectangle cellCoords = Cells[i];
                    if (cellCoords.Contains(x, y))
                    {
                        return i;
                    }
                }

                return -1;
            }

        }
        #endregion
        
        TableLayoutPanel[] squarePanels;
        SudokuBoard m_currentGame;
        int m_selectedIndex = -1;

        public SudokuBoardControl()
        {
            InitializeComponent();

            squarePanel0.Tag = new SquareInfo(squarePanel0, 0);
            squarePanel1.Tag = new SquareInfo(squarePanel1, 1);
            squarePanel2.Tag = new SquareInfo(squarePanel2, 2);
            squarePanel3.Tag = new SquareInfo(squarePanel3, 3);
            squarePanel4.Tag = new SquareInfo(squarePanel4, 4);
            squarePanel5.Tag = new SquareInfo(squarePanel5, 5);
            squarePanel6.Tag = new SquareInfo(squarePanel6, 6);
            squarePanel7.Tag = new SquareInfo(squarePanel7, 7);
            squarePanel8.Tag = new SquareInfo(squarePanel8, 8);

            squarePanels = new TableLayoutPanel[] {
                squarePanel0, squarePanel1, squarePanel2, 
                squarePanel3, squarePanel4, squarePanel5, 
                squarePanel6, squarePanel7, squarePanel8, 
            };

            string valueToolTipFormat = "right-click the button and drag-and-drop on the cell to set its value";
            ValueEmptyButton.ToolTipText = valueToolTipFormat;
            Value1Button.ToolTipText = valueToolTipFormat;
            Value2Button.ToolTipText = valueToolTipFormat;
            Value3Button.ToolTipText = valueToolTipFormat;
            Value4Button.ToolTipText = valueToolTipFormat;
            Value5Button.ToolTipText = valueToolTipFormat;
            Value6Button.ToolTipText = valueToolTipFormat;
            Value7Button.ToolTipText = valueToolTipFormat;
            Value8Button.ToolTipText = valueToolTipFormat;
            Value9Button.ToolTipText = valueToolTipFormat;
        }

        #region Game info
        [Browsable(false)]
        public SudokuBoard CurrentGame
        {
            get
            {
                return m_currentGame;
            }
            set
            {
                if (m_currentGame != null)
                    m_currentGame.OnChange -= new SudokuBoard.ValueChangeHandler(m_currentGame_OnChange);

                m_currentGame = value;
                if (m_currentGame == null)
                {
                    hintButton.Enabled = solveButton.Enabled = restartButton.Enabled = evaluateButton.Enabled = false;
                    return;
                }

                m_currentGame.OnChange += new SudokuBoard.ValueChangeHandler(m_currentGame_OnChange);
                if (CurrentGame.FreeCellsCount > 0)
                {
                    hintButton.Enabled = true;
                    solveButton.Enabled =
                        CurrentGame.GameLevel == SudokuLevel.Easy ||
                        CurrentGame.GameLevel == SudokuLevel.Unclassified;
                }
                else
                {
                    hintButton.Enabled = solveButton.Enabled = false;
                }
                evaluateButton.Enabled = true;
                GamePanelRefresh();
            }
        }

        void m_currentGame_OnChange(object sender, SudokuBoard.ValueChangeEventArgs e)
        {
            GamePanelRefresh();
        }

        int[] CellsCache = new int[81];
        void GamePanelRefresh()
        {
            bool isSolved = CurrentGame.IsSolved();

            if (CurrentGame == null)
            {
                hintButton.Enabled = solveButton.Enabled = evaluateButton.Enabled = false;
            }
            else
            {
                hintButton.Enabled = solveButton.Enabled = !isSolved;
                evaluateButton.Enabled = true;
            }

            for (int square = 0; square < 9; square++)
            {
                bool refreshNeeded = false;
                for (int cell = 0; cell < 9; cell++)
                {
                    int i = SudokuBoard.GetIndexFromSquare(square, cell);
                    int val = CurrentGame[i];
                    if (CurrentGame.IsDefault(i))
                        val |= 0x100;
                    else if (CurrentGame.IsValidCell(i))
                        val |= 0x200;

                    if (m_selectedIndex == i)
                        val |= 0x400;

                    if (CellsCache[i] != val)
                    {
                        CellsCache[i] = val;
                        refreshNeeded = true;
                    }
                }
                if (refreshNeeded)
                    squarePanels[square].Refresh();
            }

            if (isSolved)
            {
                // done, make the board read-only
                if (!CurrentGame.ReadOnly)
                {
                    OnSolved(new SolvedEventArgs(m_autoSolveInProgress));
                    CurrentGame.ReadOnly = true;
                }
            }

            restartButton.Enabled = CurrentGame.CanRestart;
        }

        public class SolvedEventArgs : EventArgs
        {
            public SolvedEventArgs(bool autoSolved)
            {
                this.AutoSolve = autoSolved;
            }

            public readonly bool AutoSolve;
        }

        public delegate void SolvedEventHandler(object sender, SolvedEventArgs e);

        public event SolvedEventHandler Solved;

        protected virtual void OnSolved(SolvedEventArgs args)
        {
            Solved?.Invoke(this, args);
        }

        #endregion

        #region CellPaint events

        Pen redPen = Pens.Red;
        Brush selectedSolidBrush = Brushes.Honeydew;
        SolidBrush notSelectedSolidBrush = new SolidBrush(SystemColors.Control);

        Font drawFontBold;
        Font drawFontNonBold;
        float drawFontSize = -1;

        void CellPaint(TableLayoutPanel square, Graphics graphics, Rectangle bounds, int index, int cellIndex, bool isFirst)
        {
            SquareInfo info = (SquareInfo)square.Tag;
            info.Cells[cellIndex] = bounds;

            byte cellValue;

            if (CurrentGame == null)
                cellValue = 0;
            else
                cellValue = CurrentGame[index];

            Brush fillBrush;
            if (m_selectedIndex == index)
            {
                fillBrush = selectedSolidBrush;
                graphics.FillEllipse(fillBrush, bounds);
            }
            else
            {
                fillBrush = notSelectedSolidBrush;
                graphics.FillRectangle(fillBrush, bounds);
            }

            if (cellValue == 0)
                return;

            String drawString = "" + (int)cellValue;
            SolidBrush drawBrush;
            Font drawFont;

            if (bounds.Width > 0 && bounds.Height > 0)
            {

                float fontSizeByWidth = (float)bounds.Width / 3;
                float fontSizeByHeight = ((float)bounds.Height / 2) * 0.64F;

                float fontSize = Math.Min(fontSizeByWidth, fontSizeByHeight);
                if (drawFontSize != fontSize)
                {
                    drawFontBold = new Font("Arial", fontSize, FontStyle.Bold);
                    drawFontNonBold = new Font("Arial", fontSize, (FontStyle)0);
                    drawFontSize = fontSize;
                }
            }

            if (CurrentGame.IsDefault(index))
            {
                drawFont = drawFontBold;
                drawBrush = new SolidBrush(Color.Blue);
            }
            else
            {
                drawFont = drawFontNonBold;
                drawBrush = new SolidBrush(Color.Black);
            }


            // calculate the relative position within the cell
            float x = (float)bounds.X + (bounds.Width - drawFont.Size) / 2;
            float y = (float)bounds.Y + (bounds.Height - drawFont.Height) / 2;

            // Create point for upper-left corner of drawing.
            PointF drawPoint = new PointF(x, y);

            // Draw string to screen.
            graphics.DrawString(drawString, drawFont, drawBrush, drawPoint);

            if (!CurrentGame.IsValidCell(index))
            {
                int temp = (bounds.Right - bounds.Left) / 6;
                int rLeft = bounds.Left + temp;
                int rRight = bounds.Right - temp;

                temp = (bounds.Bottom - bounds.Top) / 6;
                int rTop = bounds.Top + temp;
                int rBottom = bounds.Bottom - temp;

                graphics.DrawLine(redPen, rRight, rTop, rLeft, rBottom);
                graphics.DrawLine(redPen, rLeft, rTop, rRight, rBottom);
            }
        }
        private void squarePanel0_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel0, e.Graphics, e.CellBounds, e.Row * 9 + e.Column, e.Row * 3 + e.Column, true);
        }

        private void squarePanel1_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel1, e.Graphics, e.CellBounds, e.Row * 9 + 3 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel2_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel2, e.Graphics, e.CellBounds, e.Row * 9 + 6 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel3_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel3, e.Graphics, e.CellBounds, (3 + e.Row) * 9 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel4_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel4, e.Graphics, e.CellBounds, (3 + e.Row) * 9 + 3 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel5_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel5, e.Graphics, e.CellBounds, (3 + e.Row) * 9 + 6 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel6_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel6, e.Graphics, e.CellBounds, (6 + e.Row) * 9 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel7_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel7, e.Graphics, e.CellBounds, (6 + e.Row) * 9 + 3 + e.Column, e.Row * 3 + e.Column, false);
        }

        private void squarePanel8_CellPaint(object sender, TableLayoutCellPaintEventArgs e)
        {
            CellPaint(squarePanel8, e.Graphics, e.CellBounds, (6 + e.Row) * 9 + 6 + e.Column, e.Row * 3 + e.Column, false);
        }
        #endregion

        #region Value Drag & Drop events

        void SquareDragOver(TableLayoutPanel square, DragEventArgs e)
        {
            if (CurrentGame == null)
            {
                e.Effect = DragDropEffects.None;
                return;
            }

            SquareInfo info = (SquareInfo)square.Tag;
            int cell = info.GetCellFromScreen(e.X, e.Y);

            if (cell == -1)
            {
                e.Effect = DragDropEffects.None;
                Trace.WriteLine("drag over - 1");
                return;
            }

            int index = SudokuBoard.GetIndexFromSquare(info.SquareIndex, cell);
            if (CurrentGame.IsReadOnly(index))
                e.Effect = DragDropEffects.None;
            else
                e.Effect = DragDropEffects.Move;

            Trace.WriteLine("drag over " + e.Effect);
        }

        void SquareDragDrop(TableLayoutPanel square, DragEventArgs e)
        {
            if (CurrentGame == null)
                return;

            SquareInfo info = (SquareInfo)square.Tag;
            int cell = info.GetCellFromScreen(e.X, e.Y);

            if (cell != -1)
            {
                byte number = (byte)e.Data.GetData(typeof(byte));
                int index = SudokuBoard.GetIndexFromSquare(info.SquareIndex, cell);
                m_selectedIndex = index;
                CurrentGame.SetValue(index, number, isReadOnly: false);
            }
        }

        void ValueButtonDrag(MouseEventArgs mouseEvent, byte value)
        {
            if (CurrentGame == null)
                return;

            if (mouseEvent.Button != MouseButtons.Right)
                return;

            gamePanel.DoDragDrop(value, DragDropEffects.Move);
        }

        private void squarePanel_DragDrop(object sender, DragEventArgs e)
        {
            SquareDragDrop((TableLayoutPanel)sender, e);
        }

        private void squarePanel_DragOver(object sender, DragEventArgs e)
        {
            SquareDragOver((TableLayoutPanel)sender, e);
        }

        private void ValueEmptyButton_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 0);
        }

        private void Value1Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 1);
        }

        private void Value2Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 2);
        }

        private void Value3Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 3);
        }

        private void Value4Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 4);
        }

        private void Value5Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 5);
        }

        private void Value6Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 6);
        }

        private void Value7Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 7);
        }

        private void Value8Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 8);
        }

        private void Value9Button_MouseDown(object sender, MouseEventArgs e)
        {
            ValueButtonDrag(e, 9);
        }

        #endregion

        #region input processing

        private void squarePanel_MouseClick(object sender, MouseEventArgs e)
        {
            SquareMouseClick((TableLayoutPanel)sender, e);
        }

        public bool ProcessGameBoardChar(char keyChar)
        {
            if (CurrentGame == null)
                return false;

            if (m_selectedIndex < 0 || CurrentGame.IsReadOnly(m_selectedIndex))
                return false;

            byte val = 255;

            if (keyChar >= '0' && keyChar <= '9')
                val = (byte)(keyChar - '0');
            else if (keyChar == ' ')
                val = 0;
            else
                return false;

            CurrentGame.SetValue(m_selectedIndex, val, isReadOnly: false);
            return true;
        }

        public bool ProcessGameBoardKey(Keys keyData)
        {
            if (CurrentGame == null)
                return false;

            bool handled = true;

            switch (keyData)
            {
                case Keys.Right:
                    if (m_selectedIndex < 0)
                        m_selectedIndex = 0;
                    else
                        m_selectedIndex = (m_selectedIndex + 1) % 81;
                    break;
                case Keys.Left:
                    if (m_selectedIndex < 0)
                        m_selectedIndex = 0;
                    else if (m_selectedIndex == 0)
                        m_selectedIndex = 80;
                    else
                        m_selectedIndex = m_selectedIndex - 1;
                    break;
                case Keys.Up:
                    if (m_selectedIndex < 0)
                        m_selectedIndex = 0;
                    else if (m_selectedIndex == 0)
                        m_selectedIndex = 80;
                    else if (m_selectedIndex < 9)
                        m_selectedIndex = 71 + m_selectedIndex;
                    else
                        m_selectedIndex = m_selectedIndex - 9;
                    break;
                case Keys.Down:
                    if (m_selectedIndex < 0)
                        m_selectedIndex = 0;
                    else if (m_selectedIndex == 80)
                        m_selectedIndex = 0;
                    else if (m_selectedIndex > 71)
                        m_selectedIndex = m_selectedIndex - 71;
                    else
                        m_selectedIndex = m_selectedIndex + 9;
                    break;
                case Keys.Delete:
                    if (m_selectedIndex < 0 || CurrentGame.IsReadOnly(m_selectedIndex))
                        return false;

                    CurrentGame.SetValue(m_selectedIndex, 0, isReadOnly: false);
                    break;
                default:
                    handled = false;
                    break;
            }

            if (handled)
                GamePanelRefresh();

            return handled;
        }

        void SquareMouseClick(TableLayoutPanel square, MouseEventArgs e)
        {
            if (CurrentGame == null)
                return;

            SquareInfo info = (SquareInfo)square.Tag;
            int cell = info.GetCellFromTable(e.X, e.Y);

            if (cell < 0)
                return;

            m_selectedIndex = SudokuBoard.GetIndexFromSquare(info.SquareIndex, cell);
            GamePanelRefresh();
        }

        void SetValueToSelectedIndex(byte val)
        {
            if (CurrentGame == null)
                return;
            
            if (m_selectedIndex < 0 || CurrentGame.IsReadOnly(m_selectedIndex))
                return;

            CurrentGame.SetValue(m_selectedIndex, val, isReadOnly: false);
        }

        private void ValueEmptyButton_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(0);
        }

        private void Value1Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(1);
        }

        private void Value2Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(2);
        }

        private void Value3Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(3);
        }

        private void Value4Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(4);
        }

        private void Value5Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(5);
        }

        private void Value6Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(6);
        }

        private void Value7Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(7);
        }

        private void Value8Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(8);
        }

        private void Value9Button_Click(object sender, EventArgs e)
        {
            SetValueToSelectedIndex(9);
        }

        protected override bool ProcessDialogKey(Keys keyData)
        {
            bool handled = ProcessGameBoardKey(keyData);
            if (handled)
                return true;
            else
                return base.ProcessDialogKey(keyData);
        }

        protected override bool ProcessDialogChar(char charCode)
        {
            bool handled = ProcessGameBoardChar(charCode);
            if (handled)
                return true;
            else
                return base.ProcessDialogChar(charCode);
        }

        #endregion

        public void Restart()
        {
            MessageForm askForm = new MessageForm("Are you sure you want to restart the current game?");

            if (askForm.ShowDialog(this) != DialogResult.Yes)
                return;

            CurrentGame.ReadOnly = false;
            CurrentGame.RestartGame();
        }

        public void ShowAboutDialog()
        {
            new SudokuAboutForm().ShowDialog(this);
        }

        private void restartButton_Click(object sender, EventArgs e)
        {
            Restart();
        }

        private void aboutButton_Click(object sender, EventArgs e)
        {
            ShowAboutDialog();
        }

        bool m_autoSolveInProgress;
        private void solveButton_Click(object sender, EventArgs e)
        {
            if (CurrentGame == null)
                return;

            m_autoSolveInProgress = true;
            SudokuSolver.SolutionInfo solution;

            solution = SudokuSolver.Solve(CurrentGame);
            if (solution.Result.Result == SudokuAlgorithmResult.Unknown)
            {
                new HintMessage("No solution has been found!", drawFontBold, TimeSpan.FromSeconds(3), Point.Empty,
                    System.Drawing.Size.Empty).ShowDialog(this);
            }

            m_autoSolveInProgress = false;
        }

        private void hintButton_Click(object sender, EventArgs e)
        {
            if (CurrentGame == null)
                return;

            var hintGame = new SudokuBoard(CurrentGame);
            SudokuSolver.SolutionInfo res = SudokuSolver.Hint(hintGame);

            if (res.Result.Result == SudokuAlgorithmResult.Unknown)
            {
                // cannot solve this board - ask the user
                MessageForm ask = new MessageForm("Sudoku Solver did not find any algorithm to solve this board. Do you want it to search harder?");
                ask.ShowDialog(this);
                if (ask.DialogResult != DialogResult.Yes)
                    return;

                res = SudokuSolver.Hint(CurrentGame, 4 /* increase recursion level from 2 to 4 */);
            }

            if (res.Result.Result == SudokuAlgorithmResult.Succeeded)
            {
                // show second hint with the index and the value
                string msg = res.Result.Value.ToString();
                int square, cellIndex;
                SudokuBoard.GetSquareFromIndex(res.Result.ValueIndex, out square, out cellIndex);
                SquareInfo squareInfo = (SquareInfo)squarePanels[square].Tag;
                new HintMessage(msg, drawFontBold, TimeSpan.FromSeconds(0.2),
                    squarePanels[square].PointToScreen(squareInfo.Cells[cellIndex].Location),
                    squareInfo.Cells[cellIndex].Size).ShowDialog(this);

                CurrentGame.SetValue(res.Result.ValueIndex, res.Result.Value, false);

                m_selectedIndex = res.Result.ValueIndex;
                GamePanelRefresh();
            }
        }

        private void evaluateButton_Click(object sender, EventArgs e)
        {
            if (CurrentGame == null)
                return;

            m_autoSolveInProgress = true;
            
            Stopwatch watch = Stopwatch.StartNew();

            var evaluateGame = new SudokuBoard(CurrentGame);
            
            int proveTimeMs, solveTimeMs = 0, solveWeight = 0;
            IReadOnlyDictionary<int, int> solveStats = null;

            // prove first!
            var proveSolution = SudokuSolver.Prove(evaluateGame);
            watch.Stop();

            proveTimeMs = (int)watch.Elapsed.TotalMilliseconds;
            
            if (proveSolution.Result.Result != SudokuAlgorithmResult.Succeeded)
            {
                new HintMessage(proveSolution.Result.FailureDetails, drawFontBold, TimeSpan.FromSeconds(3), Point.Empty,
                    System.Drawing.Size.Empty).ShowDialog(this);
            }
            else
            {
                // restart
                evaluateGame = new SudokuBoard(CurrentGame);

                watch = Stopwatch.StartNew();
                var solveSolution = SudokuSolver.Solve(evaluateGame);
                watch.Stop();
                solveTimeMs = (int)watch.Elapsed.TotalMilliseconds;
                solveWeight = solveSolution.TotalWeight;

                // get the stats
                solveStats = solveSolution.Statistics;
            }


            StringBuilder evaluationMsg = new StringBuilder();
            evaluationMsg.AppendLine($"Solve: {solveTimeMs} ms, weight = {solveWeight}");
            evaluationMsg.AppendLine($"Prove: {(int)proveTimeMs} ms, weight = {proveSolution.TotalWeight}");

            if (solveStats != null && solveStats.Count > 0)
            {
                evaluationMsg.AppendLine();

                foreach (var entry in solveStats)
                {
                    evaluationMsg.AppendLine($"{entry.Value}\tX [ {(SudokuAlgorithmWeight)entry.Key} ]");
                }
            }
            
            MessageBox.Show(this.Parent, evaluationMsg.ToString());
            m_autoSolveInProgress = false;
        }
    }
}
