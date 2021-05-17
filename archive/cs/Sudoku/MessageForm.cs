using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;

namespace Sudoku
{
    public partial class MessageForm : Form
    {
        public MessageForm(string msg)
        {
            InitializeComponent();

            this.Font = new Font(messageLabel.Font.FontFamily, messageLabel.Font.Size + 1, FontStyle.Bold);
            messageLabel.Text = msg;
        }
    }
}
