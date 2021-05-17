namespace Sudoku
{
    partial class HintMessage
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            this.notificationControl = new System.Windows.Forms.Label();
            this.timer1 = new System.Windows.Forms.Timer(this.components);
            this.SuspendLayout();
            // 
            // notificationControl
            // 
            this.notificationControl.AutoSize = true;
            this.notificationControl.BackColor = System.Drawing.SystemColors.ControlDark;
            this.notificationControl.BorderStyle = System.Windows.Forms.BorderStyle.Fixed3D;
            this.notificationControl.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.notificationControl.Location = new System.Drawing.Point(193, 107);
            this.notificationControl.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.notificationControl.Name = "notificationControl";
            this.notificationControl.Size = new System.Drawing.Size(83, 22);
            this.notificationControl.TabIndex = 0;
            this.notificationControl.Text = "message";
            this.notificationControl.TextAlign = System.Drawing.ContentAlignment.MiddleCenter;
            this.notificationControl.MouseLeave += new System.EventHandler(this.notificationControl_MouseLeave);
            this.notificationControl.Click += new System.EventHandler(this.label1_Click);
            this.notificationControl.MouseEnter += new System.EventHandler(this.notificationControl_MouseEnter);
            // 
            // timer1
            // 
            this.timer1.Interval = 50;
            this.timer1.Tick += new System.EventHandler(this.timer1_Tick);
            // 
            // HintMessage
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(10F, 20F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.Color.Magenta;
            this.ClientSize = new System.Drawing.Size(468, 237);
            this.Controls.Add(this.notificationControl);
            this.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(204)));
            this.ForeColor = System.Drawing.Color.GreenYellow;
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.None;
            this.Margin = new System.Windows.Forms.Padding(4);
            this.Name = "HintMessage";
            this.Opacity = 0.02;
            this.ShowInTaskbar = false;
            this.StartPosition = System.Windows.Forms.FormStartPosition.Manual;
            this.Text = "HintMessage";
            this.TransparencyKey = System.Drawing.Color.Magenta;
            this.FormClosed += new System.Windows.Forms.FormClosedEventHandler(this.HintMessage_FormClosed);
            this.Load += new System.EventHandler(this.HintMessage_Load);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Label notificationControl;
        private System.Windows.Forms.Timer timer1;

    }
}