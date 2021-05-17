using System;
using System.Reflection;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;
using System.Diagnostics;

using Sudoku.Properties;

namespace Sudoku
{
    public partial class SudokuForm : Form
    {
        public SudokuForm()
        {
            InitializeComponent();

            CurrentGame = null;
        }

        SudokuBoard CurrentGame
        {
            get
            {
                return sudokuControl.CurrentGame;
            }
            set
            {
                if (sudokuControl.CurrentGame != null)
                    sudokuControl.CurrentGame.OnChange -= new SudokuBoard.ValueChangeHandler(m_currentGame_OnChange);

                sudokuControl.CurrentGame = value;
                if (sudokuControl.CurrentGame != null)
                {
                    value.OnChange += new SudokuBoard.ValueChangeHandler(m_currentGame_OnChange);
                    gameIdLabel.Text = "Game ID: " + value.GameId;
                    levelLabel.Text = "Level: " + value.GameLevel;
                }
                else
                {
                    gameIdLabel.Text = "";
                    levelLabel.Text = "";
                }

                restartMenuItem.Enabled = saveGameMenuItem.Enabled = (value != null);

                UpdateGameBoard();
            }
        }

        void m_currentGame_OnChange(object sender, SudokuBoard.ValueChangeEventArgs e)
        {
            UpdateGameBoard();
        }

        void UpdateGameBoard()
        {
            if (CurrentGame != null)
            {
                if (!CurrentGame.IsValid)
                    hintLabel.Text = "Wrong cells!";
                else
                    hintLabel.Text = "";

                int emptyCells = CurrentGame.FreeCellsCount;
                if (emptyCells > 0 || !CurrentGame.IsValid)
                    cellsLeftLabel.Text = "Cells left: " + emptyCells;
                else
                    cellsLeftLabel.Text = "Done!";

                restartMenuItem.Enabled = CurrentGame.CanRestart;
            }
            else
            {
                hintLabel.Text = "No game";
                cellsLeftLabel.Text = "";
            }
        }

        private void exitToolStripMenuItem_Click(object sender, EventArgs e)
        {
            this.Close();
        }

        #region SetLevel

        SudokuLevel GameLevel
        {
            get { return (SudokuLevel)Settings.Default.Level; }
            set
            {
                easyLevelStripMenuItem.Checked = (value == SudokuLevel.Easy);
                mediumLevelStripMenuItem.Checked = (value == SudokuLevel.Medium);
                hardLevelStripMenuItem.Checked = (value == SudokuLevel.Hard);
                veryHardLevelStripMenuItem.Checked = (value == SudokuLevel.VeryHard);
                evilLevelStripMenuItem.Checked = (value == SudokuLevel.Evil);
                darkEvilLevelStripMenuItem.Checked = (value == SudokuLevel.DarkEvil);
                blackHoleToolStripMenuItem.Checked = (value == SudokuLevel.BlackHole);
                Settings.Default.Level = (int)value;
            }
        }

        #endregion

        private void SudokuForm_Load(object sender, EventArgs e)
        {
            GameLevel = (SudokuLevel)Settings.Default.Level;
            if (GameLevel < SudokuLevel.Easy || GameLevel > SudokuLevel.BlackHole)
                GameLevel = SudokuLevel.Easy;
        }

        void LoadRandomGame()
        {
            SudokuBoard board = SudokuGameGenerator.CreateRandomGame(this, GameLevel);
            if (board != null)
                CurrentGame = board;
        }
        
        private void SudokuForm_FormClosed(object sender, FormClosedEventArgs e)
        {
            Settings.Default.Save();
        }

        private void restartMenuItem_Click(object sender, EventArgs e)
        {
            sudokuControl.Restart();
        }

        private void newGameMenuItem_Click(object sender, EventArgs e)
        {
            StartNewGame();
        }
        
        private void aboutMenuItem_Click(object sender, EventArgs e)
        {
            sudokuControl.ShowAboutDialog();
        }

        void StartNewGame()
        {
            if (CurrentGame != null && !CurrentGame.ReadOnly)
            {
                MessageForm yesNoForm = new MessageForm("Are you sure you want to stop the current game and start a new one?");
                if (yesNoForm.ShowDialog(this) != DialogResult.Yes)
                    return;
            }

            LoadRandomGame();
        }

        private void sudokuBoard_Solved(object sender, SudokuBoardControl.SolvedEventArgs e)
        {
            if (!e.AutoSolve)
            {
                MessageForm yesNoForm = new MessageForm("Well done!\nDo you want to start a new game?");
                if (yesNoForm.ShowDialog(this) != DialogResult.Yes)
                    return;

                StartNewGame();
            }
        }

        private void easyLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.Easy;
        }

        private void mediumLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.Medium;
        }

        private void hardLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.Hard;
        }

        private void veryHardLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.VeryHard;
        }

        private void evilLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.Evil;
        }

        private void darkEvilLevelStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.DarkEvil;
        }

        private void blackHoleToolStripMenuItem_Click(object sender, EventArgs e)
        {
            GameLevel = SudokuLevel.BlackHole;
        }

        private void saveGameMenuItem_Click(object sender, EventArgs e)
        {
            SaveFileDialog dlg = new SaveFileDialog();
            dlg.DefaultExt = "sudoku";
            dlg.Filter = "Sudoku Game files (*.sudoku)|*.sudoku|All files (*.*)|*.*";
            if (dlg.ShowDialog(this) != DialogResult.OK)
                return;

            try
            {
                CurrentGame.SaveToFile(dlg.FileName);
            }
            catch(Exception ex)
            {
                MessageBox.Show(this, "Failed to save the file: " + ex.Message, "Save File Error!", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }
        }

        private void loadGameMenuItem_Click(object sender, EventArgs e)
        {
            OpenFileDialog dlg = new OpenFileDialog();
            dlg.DefaultExt = "sudoku";
            dlg.Filter = "Sudoku Game files (*.sudoku)|*.sudoku|All files (*.*)|*.*";
            if (dlg.ShowDialog(this) != DialogResult.OK)
                return;

            try
            {
                SudokuBoard board = SudokuBoard.LoadFromFile(dlg.FileName);
                if (board.IsSolved())
                {
                    board.ReadOnly = true;
                }
                CurrentGame = board;
            }
            catch (Exception ex)
            {
                MessageBox.Show(this, "Failed to open the game file: " + ex.Message, "Open File Error!", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }
        }

        GameGenerationForm generatorForm;
        private void generatorGameMenuItem_Click(object sender, EventArgs e)
        {
            if (generatorForm == null)
            {
                generatorForm = new GameGenerationForm();
                generatorForm.Show();

                generatorForm.FormClosing += GeneratorForm_FormClosing;
            }
            else
            {
                generatorForm.Focus();
            }
        }

        private void GeneratorForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            generatorForm = null;
        }

        private void loadClipboardMenuItem_Click(object sender, EventArgs e)
        {
            string text = Clipboard.GetText(TextDataFormat.Text);
            try
            {
                SudokuBoard board = SudokuBoard.LoadFromText(text);
                if (board.IsSolved())
                {
                    board.ReadOnly = true;
                }
                CurrentGame = board;
            }
            catch (Exception ex)
            {
                MessageBox.Show(this, "Failed to convert text from clipboard to the game: " + ex.Message, "Open Clipboard Error!", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }
        }

        private void exitMenuItem_Click(object sender, EventArgs e)
        {
            // this.Close will cause FormClosing event with CloseReason == CloseReason.UserClosing
            // that is handled
            this.Close();
        }

        private void SudokuForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            if (e.CloseReason == CloseReason.UserClosing)
            {
                if (CurrentGame == null)
                    return;

                MessageForm yesNoForm = new MessageForm("Are you sure you want to stop the current game and exit?");
                if (yesNoForm.ShowDialog(this) != DialogResult.No)
                    return;

                e.Cancel = true;
            }
        }
    }
}
