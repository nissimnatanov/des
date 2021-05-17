namespace Sudoku
{
    partial class GameGenerationForm
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
            this.actionButton = new System.Windows.Forms.Button();
            this.label1 = new System.Windows.Forms.Label();
            this.numberOfGamesBox = new System.Windows.Forms.ComboBox();
            this.progressBar1 = new System.Windows.Forms.ProgressBar();
            this.statsTextBox = new System.Windows.Forms.RichTextBox();
            this.copyButton = new System.Windows.Forms.Button();
            this.SuspendLayout();
            // 
            // actionButton
            // 
            this.actionButton.Anchor = System.Windows.Forms.AnchorStyles.Bottom;
            this.actionButton.Location = new System.Drawing.Point(65, 326);
            this.actionButton.Name = "actionButton";
            this.actionButton.Size = new System.Drawing.Size(75, 23);
            this.actionButton.TabIndex = 0;
            this.actionButton.Text = "Start";
            this.actionButton.UseVisualStyleBackColor = true;
            this.actionButton.Click += new System.EventHandler(this.actionButton_Click);
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(13, 13);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(150, 13);
            this.label1.TabIndex = 1;
            this.label1.Text = "Number of games to generate:";
            // 
            // numberOfGamesBox
            // 
            this.numberOfGamesBox.DropDownStyle = System.Windows.Forms.ComboBoxStyle.DropDownList;
            this.numberOfGamesBox.FormattingEnabled = true;
            this.numberOfGamesBox.Location = new System.Drawing.Point(169, 10);
            this.numberOfGamesBox.Name = "numberOfGamesBox";
            this.numberOfGamesBox.Size = new System.Drawing.Size(103, 21);
            this.numberOfGamesBox.TabIndex = 2;
            // 
            // progressBar1
            // 
            this.progressBar1.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.progressBar1.Location = new System.Drawing.Point(16, 37);
            this.progressBar1.Name = "progressBar1";
            this.progressBar1.Size = new System.Drawing.Size(256, 23);
            this.progressBar1.TabIndex = 3;
            // 
            // statsTextBox
            // 
            this.statsTextBox.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.statsTextBox.Location = new System.Drawing.Point(16, 67);
            this.statsTextBox.Name = "statsTextBox";
            this.statsTextBox.ReadOnly = true;
            this.statsTextBox.Size = new System.Drawing.Size(256, 253);
            this.statsTextBox.TabIndex = 4;
            this.statsTextBox.Text = "";
            // 
            // copyButton
            // 
            this.copyButton.Anchor = System.Windows.Forms.AnchorStyles.Bottom;
            this.copyButton.Enabled = false;
            this.copyButton.Location = new System.Drawing.Point(146, 326);
            this.copyButton.Name = "copyButton";
            this.copyButton.Size = new System.Drawing.Size(75, 23);
            this.copyButton.TabIndex = 5;
            this.copyButton.Text = "Copy";
            this.copyButton.UseVisualStyleBackColor = true;
            this.copyButton.Click += new System.EventHandler(this.copyButton_Click);
            // 
            // GameGenerationForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(284, 361);
            this.Controls.Add(this.copyButton);
            this.Controls.Add(this.statsTextBox);
            this.Controls.Add(this.progressBar1);
            this.Controls.Add(this.numberOfGamesBox);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.actionButton);
            this.MaximizeBox = false;
            this.MaximumSize = new System.Drawing.Size(300, 400);
            this.MinimumSize = new System.Drawing.Size(300, 400);
            this.Name = "GameGenerationForm";
            this.Text = "Sudoku Game Generator";
            this.FormClosing += new System.Windows.Forms.FormClosingEventHandler(this.GameGenerationForm_FormClosing);
            this.Load += new System.EventHandler(this.GameGenerationForm_Load);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Button actionButton;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.ComboBox numberOfGamesBox;
        private System.Windows.Forms.ProgressBar progressBar1;
        private System.Windows.Forms.RichTextBox statsTextBox;
        private System.Windows.Forms.Button copyButton;
    }
}