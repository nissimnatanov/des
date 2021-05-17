using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;

namespace Sudoku
{
    public partial class HintMessage : Form
    {
        TimeSpan m_showTime;
        DateTime m_shownAt;
        bool m_forceShow;
        bool m_shown;

        public HintMessage(string msg, Font font, TimeSpan showTime, Point formLocation, Size formSize)
        {
            InitializeComponent();

            if (font != null)
                notificationControl.Font = font;

            notificationControl.Text = msg;
            // the label is resized
            System.Drawing.Size minFormSize = notificationControl.Size + this.Size - this.ClientSize;

            if (!formSize.IsEmpty)
            {
                Point loc = new Point(0, 0);
                // place the message location in the middle
                if (formSize.Height < minFormSize.Height)
                {
                    formSize.Height = minFormSize.Height;
                }
                else
                {
                    loc.Y = (formSize.Height - minFormSize.Height) / 2;
                }
                if (formSize.Width < minFormSize.Width)
                {
                    formSize.Width = minFormSize.Width;
                }
                else
                {
                    loc.X = (formSize.Width - minFormSize.Width) / 2;
                }

                this.Size = formSize;
                notificationControl.Location = loc;
            }
            else
            {
                notificationControl.Location = new Point(0, 0);
                this.Size = minFormSize;
            }

            if (!formLocation.IsEmpty)
            {
                this.StartPosition = FormStartPosition.Manual;
                this.Location = formLocation;
            }
            else
            {
                this.StartPosition = FormStartPosition.CenterParent;
            }

            m_showTime = showTime;
        }

        private void label1_Click(object sender, EventArgs e)
        {
            this.Close();
        }

        private void HintMessage_Load(object sender, EventArgs e)
        {
            timer1.Enabled = true;
        }

        private void timer1_Tick(object sender, EventArgs e)
        {
            if (this.IsDisposed || !this.Visible)
                return;

            timer1.Enabled = false;
            bool reenableTimer = true;

            if (!m_shown)
            {
                // show the message by increasing the form's opacity
                double opacity = this.Opacity + 0.2;
                if (opacity >= 1.0)
                {
                    opacity = 1.0;
                    m_shown = true;
                    m_shownAt = DateTime.Now;
                }
                this.Opacity = opacity;
            }
            else if (m_showTime != TimeSpan.Zero && !m_forceShow)
            {
                TimeSpan passedTime = DateTime.Now - m_shownAt;
                if (passedTime > m_showTime)
                {
                    // hide the message
                    double opacity = this.Opacity - 0.1;
                    if (opacity <= 0)
                    {
                        opacity = 0;
                        reenableTimer = false;
                    }
                    this.Opacity = opacity;
                }
            }

            if (reenableTimer)
                timer1.Enabled = true;
            else
                this.Close();
        }

        private void HintMessage_FormClosed(object sender, FormClosedEventArgs e)
        {
            timer1.Enabled = false;

        }

        private void notificationControl_MouseEnter(object sender, EventArgs e)
        {
            m_forceShow = true;
            this.Opacity = 1.0;
        }

        private void notificationControl_MouseLeave(object sender, EventArgs e)
        {
            m_forceShow = false;
        }
    }
}
