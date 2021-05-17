using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;

namespace Sudoku
{
    partial class WaitForm : Form
    {
        public WaitForm(string msg)
        {
            InitializeComponent();

            this.Font = new Font(messageLabel.Font.FontFamily, messageLabel.Font.Size + 1, FontStyle.Bold);
            messageLabel.Text = msg;
        }

        public void StopEventHandler(object sender, EventArgs notUsed)
        {
            if (this.IsDisposed)
                return;

            try
            {
                if (Application.MessageLoop)
                {
                    if (Visible)
                    {
                        this.DialogResult = DialogResult.OK;
                        Close();
                    }
                }
                else
                {
                    Invoke(new EventHandler(StopEventHandler), sender, notUsed);
                }
            }
            catch { }
        }

        public void ProgressEventHandler(object sender, SudokuGameGenerator.ProgressEventArgs e)
        {
            if (this.IsDisposed)
                return;

            try
            {
                if (Application.MessageLoop)
                {
                    if (Visible)
                    {
                        progressBar1.Value = e.Progress;
                    }
                }
                else
                {
                    Invoke(new SudokuGameGenerator.ProgressEventHandler(ProgressEventHandler), sender, e);
                }
            }
            catch { }
        }
    }
}
