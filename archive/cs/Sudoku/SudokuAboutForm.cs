using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;

namespace Sudoku
{
    public partial class SudokuAboutForm : Form
    {
        public SudokuAboutForm()
        {
            InitializeComponent();
        }

        private void SudokuAboutForm_Click(object sender, EventArgs e)
        {
            this.Close();
        }
    }
}