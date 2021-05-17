using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace Sudoku
{
    public partial class GameGenerationForm : Form
    {
        SudokuGameGenerator m_generator;
        public GameGenerationForm()
        {
            InitializeComponent();
        }

        class ValueCount
        {
            public string name;
            public int value;

            public override string ToString() => name;
        }

        private void GameGenerationForm_Load(object sender, EventArgs e)
        {
            numberOfGamesBox.Items.AddRange(new object[]
                {
                    new ValueCount { name = "10", value = 10 },
                    new ValueCount { name = "100", value = 100 },
                    new ValueCount { name = "1K", value = 1000 },
                    new ValueCount { name = "10K", value = 10000 },
                    new ValueCount { name = "100K", value = 100000 },
                    new ValueCount { name = "1M", value = 1000000 },
                    new ValueCount { name = "10M", value = 10000000 },
                    new ValueCount { name = "100M", value = 100000000 },
                });
            numberOfGamesBox.SelectedIndex = 3;
        }

        private void actionButton_Click(object sender, EventArgs e)
        {
            if (m_generator == null)
            {
                actionButton.Text = "Stop";
                numberOfGamesBox.Enabled = false;
                m_generator = new SudokuGameGenerator(SudokuLevel.Max);
                m_generator.OnProgress += OnProgress;
                m_generator.OnCompleted += OnCompleted;

                int requestedGames = ((ValueCount)numberOfGamesBox.SelectedItem).value;

#if DEBUG
                m_generator.StartMany(requestedGames, 1);
#else
                m_generator.StartMany(requestedGames);
#endif
            }
            else
            {
                m_generator.Stop();
            }
        }

        private void OnProgress(object sender, SudokuGameGenerator.ProgressEventArgs e)
        {
            if (this.IsDisposed)
                return;

            try
            {
                if (Application.MessageLoop)
                {
                    if (Visible && m_generator != null)
                    {
                        progressBar1.Value = e.Progress;
                        SetStats(m_generator.StatsToString(), e.GamesGenerated > 0);
                    }
                }
                else
                {
                    Invoke(new SudokuGameGenerator.ProgressEventHandler(OnProgress), sender, e);
                }
            }
            catch { }
        }

        private void OnCompleted(object sender, EventArgs e)
        {

            if (this.IsDisposed)
                return;

            try
            {
                if (Application.MessageLoop)
                {
                    if (Visible && m_generator != null)
                    {
                        actionButton.Text = "Start";
                        numberOfGamesBox.Enabled = true;
                        SetStats(m_generator.StatsToString(), true);
                    }

                    m_generator = null;
                }
                else
                {
                    Invoke(new EventHandler(OnCompleted), sender, e);
                }
            }
            catch {

            }
        }

        private void GameGenerationForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            if (m_generator != null)
            {
                m_generator.OnProgress -= OnProgress;
                m_generator.OnCompleted -= OnCompleted;
                m_generator.Stop();
                m_generator = null;
            }
        }

        private void SetStats(string stats, bool madeProgress)
        {
            statsTextBox.Text = stats ?? "";
            copyButton.Enabled = !string.IsNullOrEmpty(stats) && madeProgress;
        }

        private void copyButton_Click(object sender, EventArgs e)
        {
            string text = statsTextBox.Text;
            if (string.IsNullOrWhiteSpace(text))
                return;
            if (text.Contains('\n') && !text.Contains("\r\n"))
                text = text.Replace("\n", "\r\n");
            Clipboard.SetText(text, TextDataFormat.UnicodeText);
        }
    }
}
