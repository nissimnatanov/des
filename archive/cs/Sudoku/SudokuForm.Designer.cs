using System.Drawing;
using System.Text;
using System.Windows.Forms;
using System;

namespace Sudoku
{
    partial class SudokuForm
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
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(SudokuForm));
            this.menuStrip1 = new System.Windows.Forms.MenuStrip();
            this.gameMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.newGameMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.restartMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.levelMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.easyLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.mediumLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.hardLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.veryHardLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.evilLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.darkEvilLevelStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.blackHoleToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.generatorGameMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator2 = new System.Windows.Forms.ToolStripSeparator();
            this.loadGameMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.loadClipboardMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.saveGameMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripMenuItem1 = new System.Windows.Forms.ToolStripSeparator();
            this.exitMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.helpMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.aboutMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.statusStrip1 = new System.Windows.Forms.StatusStrip();
            this.levelLabel = new System.Windows.Forms.ToolStripStatusLabel();
            this.toolStripSeparator3 = new System.Windows.Forms.ToolStripSeparator();
            this.gameIdLabel = new System.Windows.Forms.ToolStripStatusLabel();
            this.toolStripSeparator5 = new System.Windows.Forms.ToolStripSeparator();
            this.cellsLeftLabel = new System.Windows.Forms.ToolStripStatusLabel();
            this.toolStripSeparator4 = new System.Windows.Forms.ToolStripSeparator();
            this.hintLabel = new System.Windows.Forms.ToolStripStatusLabel();
            this.sudokuControl = new Sudoku.SudokuBoardControl();
            this.menuStrip1.SuspendLayout();
            this.statusStrip1.SuspendLayout();
            this.SuspendLayout();
            // 
            // menuStrip1
            // 
            this.menuStrip1.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.gameMenuItem,
            this.helpMenuItem});
            this.menuStrip1.Location = new System.Drawing.Point(0, 0);
            this.menuStrip1.Name = "menuStrip1";
            this.menuStrip1.Size = new System.Drawing.Size(524, 24);
            this.menuStrip1.TabIndex = 4;
            this.menuStrip1.Text = "menuStrip1";
            // 
            // gameMenuItem
            // 
            this.gameMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.newGameMenuItem,
            this.restartMenuItem,
            this.levelMenuItem,
            this.generatorGameMenuItem,
            this.toolStripSeparator2,
            this.loadGameMenuItem,
            this.loadClipboardMenuItem,
            this.saveGameMenuItem,
            this.toolStripMenuItem1,
            this.exitMenuItem});
            this.gameMenuItem.Name = "gameMenuItem";
            this.gameMenuItem.Size = new System.Drawing.Size(50, 20);
            this.gameMenuItem.Text = "&Game";
            // 
            // newGameMenuItem
            // 
            this.newGameMenuItem.Name = "newGameMenuItem";
            this.newGameMenuItem.Size = new System.Drawing.Size(182, 22);
            this.newGameMenuItem.Text = "&New Game";
            this.newGameMenuItem.Click += new System.EventHandler(this.newGameMenuItem_Click);
            // 
            // restartMenuItem
            // 
            this.restartMenuItem.Name = "restartMenuItem";
            this.restartMenuItem.Size = new System.Drawing.Size(182, 22);
            this.restartMenuItem.Text = "&Restart Game";
            this.restartMenuItem.Click += new System.EventHandler(this.restartMenuItem_Click);
            // 
            // levelMenuItem
            // 
            this.levelMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.easyLevelStripMenuItem,
            this.mediumLevelStripMenuItem,
            this.hardLevelStripMenuItem,
            this.veryHardLevelStripMenuItem,
            this.evilLevelStripMenuItem,
            this.darkEvilLevelStripMenuItem,
            this.blackHoleToolStripMenuItem});
            this.levelMenuItem.Name = "levelMenuItem";
            this.levelMenuItem.Size = new System.Drawing.Size(182, 22);
            this.levelMenuItem.Text = "Game Le&vel";
            // 
            // easyLevelStripMenuItem
            // 
            this.easyLevelStripMenuItem.Checked = true;
            this.easyLevelStripMenuItem.CheckState = System.Windows.Forms.CheckState.Checked;
            this.easyLevelStripMenuItem.Name = "easyLevelStripMenuItem";
            this.easyLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.easyLevelStripMenuItem.Text = "&Easy";
            this.easyLevelStripMenuItem.Click += new System.EventHandler(this.easyLevelStripMenuItem_Click);
            // 
            // mediumLevelStripMenuItem
            // 
            this.mediumLevelStripMenuItem.Name = "mediumLevelStripMenuItem";
            this.mediumLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.mediumLevelStripMenuItem.Text = "&Medium";
            this.mediumLevelStripMenuItem.Click += new System.EventHandler(this.mediumLevelStripMenuItem_Click);
            // 
            // hardLevelStripMenuItem
            // 
            this.hardLevelStripMenuItem.Name = "hardLevelStripMenuItem";
            this.hardLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.hardLevelStripMenuItem.Text = "&Hard";
            this.hardLevelStripMenuItem.Click += new System.EventHandler(this.hardLevelStripMenuItem_Click);
            // 
            // veryHardLevelStripMenuItem
            // 
            this.veryHardLevelStripMenuItem.Name = "veryHardLevelStripMenuItem";
            this.veryHardLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.veryHardLevelStripMenuItem.Text = "&Very Hard";
            this.veryHardLevelStripMenuItem.Click += new System.EventHandler(this.veryHardLevelStripMenuItem_Click);
            // 
            // evilLevelStripMenuItem
            // 
            this.evilLevelStripMenuItem.Name = "evilLevelStripMenuItem";
            this.evilLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.evilLevelStripMenuItem.Text = "Evi&l";
            this.evilLevelStripMenuItem.Click += new System.EventHandler(this.evilLevelStripMenuItem_Click);
            // 
            // darkEvilLevelStripMenuItem
            // 
            this.darkEvilLevelStripMenuItem.Name = "darkEvilLevelStripMenuItem";
            this.darkEvilLevelStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.darkEvilLevelStripMenuItem.Text = "&Dark Evil";
            this.darkEvilLevelStripMenuItem.Click += new System.EventHandler(this.darkEvilLevelStripMenuItem_Click);
            // 
            // blackHoleToolStripMenuItem
            // 
            this.blackHoleToolStripMenuItem.Name = "blackHoleToolStripMenuItem";
            this.blackHoleToolStripMenuItem.Size = new System.Drawing.Size(130, 22);
            this.blackHoleToolStripMenuItem.Text = "&Black Hole";
            this.blackHoleToolStripMenuItem.Click += new System.EventHandler(this.blackHoleToolStripMenuItem_Click);
            // 
            // generatorGameMenuItem
            // 
            this.generatorGameMenuItem.Name = "generatorGameMenuItem";
            this.generatorGameMenuItem.Size = new System.Drawing.Size(182, 22);
            this.generatorGameMenuItem.Text = "Game Generator";
            this.generatorGameMenuItem.Click += new System.EventHandler(this.generatorGameMenuItem_Click);
            // 
            // toolStripSeparator2
            // 
            this.toolStripSeparator2.Name = "toolStripSeparator2";
            this.toolStripSeparator2.Size = new System.Drawing.Size(179, 6);
            // 
            // loadGameMenuItem
            // 
            this.loadGameMenuItem.Name = "loadGameMenuItem";
            this.loadGameMenuItem.Size = new System.Drawing.Size(182, 22);
            this.loadGameMenuItem.Text = "&Load from file";
            this.loadGameMenuItem.Click += new System.EventHandler(this.loadGameMenuItem_Click);
            // 
            // loadClipboardMenuItem
            // 
            this.loadClipboardMenuItem.Name = "loadClipboardMenuItem";
            this.loadClipboardMenuItem.Size = new System.Drawing.Size(182, 22);
            this.loadClipboardMenuItem.Text = "Load from &clipboard";
            this.loadClipboardMenuItem.Click += new System.EventHandler(this.loadClipboardMenuItem_Click);
            // 
            // saveGameMenuItem
            // 
            this.saveGameMenuItem.Name = "saveGameMenuItem";
            this.saveGameMenuItem.Size = new System.Drawing.Size(182, 22);
            this.saveGameMenuItem.Text = "&Save to file";
            this.saveGameMenuItem.Click += new System.EventHandler(this.saveGameMenuItem_Click);
            // 
            // toolStripMenuItem1
            // 
            this.toolStripMenuItem1.Name = "toolStripMenuItem1";
            this.toolStripMenuItem1.Size = new System.Drawing.Size(179, 6);
            // 
            // exitMenuItem
            // 
            this.exitMenuItem.Name = "exitMenuItem";
            this.exitMenuItem.Size = new System.Drawing.Size(182, 22);
            this.exitMenuItem.Text = "E&xit";
            this.exitMenuItem.Click += new System.EventHandler(this.exitMenuItem_Click);
            // 
            // helpMenuItem
            // 
            this.helpMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.aboutMenuItem});
            this.helpMenuItem.Name = "helpMenuItem";
            this.helpMenuItem.Size = new System.Drawing.Size(44, 20);
            this.helpMenuItem.Text = "&Help";
            // 
            // aboutMenuItem
            // 
            this.aboutMenuItem.Name = "aboutMenuItem";
            this.aboutMenuItem.Size = new System.Drawing.Size(107, 22);
            this.aboutMenuItem.Text = "&About";
            this.aboutMenuItem.Click += new System.EventHandler(this.aboutMenuItem_Click);
            // 
            // statusStrip1
            // 
            this.statusStrip1.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.levelLabel,
            this.toolStripSeparator3,
            this.gameIdLabel,
            this.toolStripSeparator5,
            this.cellsLeftLabel,
            this.toolStripSeparator4,
            this.hintLabel});
            this.statusStrip1.Location = new System.Drawing.Point(0, 438);
            this.statusStrip1.Name = "statusStrip1";
            this.statusStrip1.Size = new System.Drawing.Size(524, 23);
            this.statusStrip1.TabIndex = 0;
            this.statusStrip1.Text = "statusStrip1";
            // 
            // levelLabel
            // 
            this.levelLabel.Name = "levelLabel";
            this.levelLabel.Size = new System.Drawing.Size(37, 18);
            this.levelLabel.Text = "Level:";
            // 
            // toolStripSeparator3
            // 
            this.toolStripSeparator3.Name = "toolStripSeparator3";
            this.toolStripSeparator3.Size = new System.Drawing.Size(6, 23);
            // 
            // gameIdLabel
            // 
            this.gameIdLabel.Name = "gameIdLabel";
            this.gameIdLabel.Size = new System.Drawing.Size(55, 18);
            this.gameIdLabel.Text = "Game ID:";
            // 
            // toolStripSeparator5
            // 
            this.toolStripSeparator5.Name = "toolStripSeparator5";
            this.toolStripSeparator5.Size = new System.Drawing.Size(6, 23);
            // 
            // cellsLeftLabel
            // 
            this.cellsLeftLabel.Name = "cellsLeftLabel";
            this.cellsLeftLabel.Size = new System.Drawing.Size(55, 18);
            this.cellsLeftLabel.Text = "Cells left:";
            // 
            // toolStripSeparator4
            // 
            this.toolStripSeparator4.Name = "toolStripSeparator4";
            this.toolStripSeparator4.Size = new System.Drawing.Size(6, 23);
            // 
            // hintLabel
            // 
            this.hintLabel.Name = "hintLabel";
            this.hintLabel.Size = new System.Drawing.Size(0, 18);
            // 
            // sudokuControl
            // 
            this.sudokuControl.CurrentGame = null;
            this.sudokuControl.Dock = System.Windows.Forms.DockStyle.Fill;
            this.sudokuControl.Location = new System.Drawing.Point(0, 24);
            this.sudokuControl.Name = "sudokuControl";
            this.sudokuControl.Size = new System.Drawing.Size(524, 414);
            this.sudokuControl.TabIndex = 5;
            this.sudokuControl.Solved += new Sudoku.SudokuBoardControl.SolvedEventHandler(this.sudokuBoard_Solved);
            // 
            // SudokuForm
            // 
            this.ClientSize = new System.Drawing.Size(524, 461);
            this.Controls.Add(this.sudokuControl);
            this.Controls.Add(this.menuStrip1);
            this.Controls.Add(this.statusStrip1);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.MainMenuStrip = this.menuStrip1;
            this.MinimumSize = new System.Drawing.Size(400, 400);
            this.Name = "SudokuForm";
            this.Text = "Sudoku";
            this.FormClosing += new System.Windows.Forms.FormClosingEventHandler(this.SudokuForm_FormClosing);
            this.FormClosed += new System.Windows.Forms.FormClosedEventHandler(this.SudokuForm_FormClosed);
            this.Load += new System.EventHandler(this.SudokuForm_Load);
            this.menuStrip1.ResumeLayout(false);
            this.menuStrip1.PerformLayout();
            this.statusStrip1.ResumeLayout(false);
            this.statusStrip1.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.MenuStrip menuStrip1;
        private System.Windows.Forms.ToolStripMenuItem gameMenuItem;
        private System.Windows.Forms.ToolStripMenuItem newGameMenuItem;
        private System.Windows.Forms.ToolStripSeparator toolStripMenuItem1;
        private System.Windows.Forms.ToolStripMenuItem exitMenuItem;
        private System.Windows.Forms.ToolStripMenuItem helpMenuItem;
        private ToolStripMenuItem saveGameMenuItem;
        private ToolStripMenuItem levelMenuItem;
        private ToolStripMenuItem easyLevelStripMenuItem;
        private ToolStripMenuItem mediumLevelStripMenuItem;
        private ToolStripMenuItem hardLevelStripMenuItem;
        private ToolStripSeparator toolStripSeparator2;
        private ToolStripMenuItem loadGameMenuItem;
        private ToolStripMenuItem aboutMenuItem;
        private ToolStripMenuItem restartMenuItem;
        private StatusStrip statusStrip1;
        private ToolStripStatusLabel levelLabel;
        private ToolStripSeparator toolStripSeparator3;
        private ToolStripStatusLabel gameIdLabel;
        private ToolStripSeparator toolStripSeparator4;
        private ToolStripStatusLabel cellsLeftLabel;
        private ToolStripSeparator toolStripSeparator5;
        private ToolStripStatusLabel hintLabel;
        private SudokuBoardControl sudokuControl;
        private ToolStripMenuItem veryHardLevelStripMenuItem;
        private ToolStripMenuItem evilLevelStripMenuItem;
        private ToolStripMenuItem darkEvilLevelStripMenuItem;
        private ToolStripMenuItem loadClipboardMenuItem;
        private ToolStripMenuItem blackHoleToolStripMenuItem;
        private ToolStripMenuItem generatorGameMenuItem;
    }
}

