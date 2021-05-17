namespace Sudoku
{
    partial class MessageForm
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
            this.yesButton = new System.Windows.Forms.Button();
            this.noButton = new System.Windows.Forms.Button();
            this.messageLabel = new System.Windows.Forms.Label();
            this.SuspendLayout();
            // 
            // yesButton
            // 
            this.yesButton.Anchor = System.Windows.Forms.AnchorStyles.Bottom;
            this.yesButton.DialogResult = System.Windows.Forms.DialogResult.Yes;
            this.yesButton.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.yesButton.ForeColor = System.Drawing.Color.GreenYellow;
            this.yesButton.Location = new System.Drawing.Point(58, 70);
            this.yesButton.Name = "yesButton";
            this.yesButton.Size = new System.Drawing.Size(75, 23);
            this.yesButton.TabIndex = 2;
            this.yesButton.Text = "&Yes";
            this.yesButton.UseVisualStyleBackColor = true;
            // 
            // noButton
            // 
            this.noButton.Anchor = System.Windows.Forms.AnchorStyles.Bottom;
            this.noButton.DialogResult = System.Windows.Forms.DialogResult.No;
            this.noButton.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.noButton.ForeColor = System.Drawing.Color.GreenYellow;
            this.noButton.Location = new System.Drawing.Point(139, 70);
            this.noButton.Name = "noButton";
            this.noButton.Size = new System.Drawing.Size(75, 23);
            this.noButton.TabIndex = 1;
            this.noButton.Text = "&No";
            this.noButton.UseVisualStyleBackColor = true;
            // 
            // messageLabel
            // 
            this.messageLabel.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom)
                        | System.Windows.Forms.AnchorStyles.Left)
                        | System.Windows.Forms.AnchorStyles.Right)));
            this.messageLabel.ForeColor = System.Drawing.Color.GreenYellow;
            this.messageLabel.Location = new System.Drawing.Point(12, 9);
            this.messageLabel.Name = "messageLabel";
            this.messageLabel.Size = new System.Drawing.Size(248, 53);
            this.messageLabel.TabIndex = 5;
            this.messageLabel.Text = "message";
            this.messageLabel.TextAlign = System.Drawing.ContentAlignment.MiddleCenter;
            // 
            // MessageForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.SystemColors.ControlDark;
            this.ClientSize = new System.Drawing.Size(272, 100);
            this.Controls.Add(this.messageLabel);
            this.Controls.Add(this.noButton);
            this.Controls.Add(this.yesButton);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.None;
            this.MinimumSize = new System.Drawing.Size(180, 100);
            this.Name = "MessageForm";
            this.ShowInTaskbar = false;
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterParent;
            this.Text = "NotificationForm";
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.Button yesButton;
        private System.Windows.Forms.Button noButton;
        private System.Windows.Forms.Label messageLabel;

    }
}