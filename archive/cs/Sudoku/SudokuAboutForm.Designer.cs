namespace Sudoku
{
    partial class SudokuAboutForm
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
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(SudokuAboutForm));
            this.label1 = new System.Windows.Forms.Label();
            this.sudokuPictureBox = new System.Windows.Forms.PictureBox();
            this.closeButton = new System.Windows.Forms.Button();
            this.authorLabel = new System.Windows.Forms.Label();
            ((System.ComponentModel.ISupportInitialize)(this.sudokuPictureBox)).BeginInit();
            this.SuspendLayout();
            // 
            // label1
            // 
            this.label1.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Left)
                        | System.Windows.Forms.AnchorStyles.Right)));
            this.label1.Font = new System.Drawing.Font("Microsoft Sans Serif", 14.25F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(204)));
            this.label1.ForeColor = System.Drawing.Color.GreenYellow;
            this.label1.Location = new System.Drawing.Point(12, 9);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(214, 31);
            this.label1.TabIndex = 0;
            this.label1.Text = "Sudoku Game";
            this.label1.TextAlign = System.Drawing.ContentAlignment.TopCenter;
            // 
            // sudokuPictureBox
            // 
            this.sudokuPictureBox.Anchor = System.Windows.Forms.AnchorStyles.Top;
            this.sudokuPictureBox.Image = ((System.Drawing.Image)(resources.GetObject("sudokuPictureBox.Image")));
            this.sudokuPictureBox.Location = new System.Drawing.Point(87, 43);
            this.sudokuPictureBox.Name = "sudokuPictureBox";
            this.sudokuPictureBox.Size = new System.Drawing.Size(64, 64);
            this.sudokuPictureBox.SizeMode = System.Windows.Forms.PictureBoxSizeMode.StretchImage;
            this.sudokuPictureBox.TabIndex = 1;
            this.sudokuPictureBox.TabStop = false;
            // 
            // closeButton
            // 
            this.closeButton.Anchor = System.Windows.Forms.AnchorStyles.Top;
            this.closeButton.DialogResult = System.Windows.Forms.DialogResult.OK;
            this.closeButton.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.closeButton.ForeColor = System.Drawing.Color.GreenYellow;
            this.closeButton.Location = new System.Drawing.Point(82, 161);
            this.closeButton.Name = "closeButton";
            this.closeButton.Size = new System.Drawing.Size(75, 23);
            this.closeButton.TabIndex = 2;
            this.closeButton.Text = "Close";
            this.closeButton.UseVisualStyleBackColor = true;
            // 
            // authorLabel
            // 
            this.authorLabel.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Left)
                        | System.Windows.Forms.AnchorStyles.Right)));
            this.authorLabel.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(204)));
            this.authorLabel.ForeColor = System.Drawing.Color.GreenYellow;
            this.authorLabel.Location = new System.Drawing.Point(15, 117);
            this.authorLabel.Name = "authorLabel";
            this.authorLabel.Size = new System.Drawing.Size(209, 38);
            this.authorLabel.TabIndex = 3;
            this.authorLabel.Text = "Author: Nissim Natanov";
            this.authorLabel.TextAlign = System.Drawing.ContentAlignment.TopCenter;
            // 
            // SudokuAboutForm
            // 
            this.AcceptButton = this.closeButton;
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.SystemColors.ControlDark;
            this.CancelButton = this.closeButton;
            this.ClientSize = new System.Drawing.Size(238, 196);
            this.Controls.Add(this.authorLabel);
            this.Controls.Add(this.closeButton);
            this.Controls.Add(this.sudokuPictureBox);
            this.Controls.Add(this.label1);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.None;
            this.Name = "SudokuAboutForm";
            this.ShowInTaskbar = false;
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterParent;
            this.Text = "SudokuAboutForm";
            this.Click += new System.EventHandler(this.SudokuAboutForm_Click);
            ((System.ComponentModel.ISupportInitialize)(this.sudokuPictureBox)).EndInit();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.PictureBox sudokuPictureBox;
        private System.Windows.Forms.Button closeButton;
        private System.Windows.Forms.Label authorLabel;
    }
}