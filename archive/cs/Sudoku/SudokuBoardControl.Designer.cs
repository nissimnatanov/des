namespace Sudoku
{
    partial class SudokuBoardControl
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

        #region Component Designer generated code

        /// <summary> 
        /// Required method for Designer support - do not modify 
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(SudokuBoardControl));
            this.toolStrip1 = new System.Windows.Forms.ToolStrip();
            this.ValueEmptyButton = new System.Windows.Forms.ToolStripButton();
            this.Value1Button = new System.Windows.Forms.ToolStripButton();
            this.Value2Button = new System.Windows.Forms.ToolStripButton();
            this.Value3Button = new System.Windows.Forms.ToolStripButton();
            this.Value4Button = new System.Windows.Forms.ToolStripButton();
            this.Value5Button = new System.Windows.Forms.ToolStripButton();
            this.Value6Button = new System.Windows.Forms.ToolStripButton();
            this.Value7Button = new System.Windows.Forms.ToolStripButton();
            this.Value8Button = new System.Windows.Forms.ToolStripButton();
            this.Value9Button = new System.Windows.Forms.ToolStripButton();
            this.toolStripSeparator1 = new System.Windows.Forms.ToolStripSeparator();
            this.hintButton = new System.Windows.Forms.ToolStripButton();
            this.solveButton = new System.Windows.Forms.ToolStripButton();
            this.toolStripSeparator2 = new System.Windows.Forms.ToolStripSeparator();
            this.restartButton = new System.Windows.Forms.ToolStripButton();
            this.toolStripSeparator3 = new System.Windows.Forms.ToolStripSeparator();
            this.aboutButton = new System.Windows.Forms.ToolStripButton();
            this.gamePanel = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel8 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel7 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel6 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel5 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel4 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel3 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel2 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel1 = new System.Windows.Forms.TableLayoutPanel();
            this.squarePanel0 = new System.Windows.Forms.TableLayoutPanel();
            this.evaluateButton = new System.Windows.Forms.ToolStripButton();
            this.toolStrip1.SuspendLayout();
            this.gamePanel.SuspendLayout();
            this.SuspendLayout();
            // 
            // toolStrip1
            // 
            this.toolStrip1.GripStyle = System.Windows.Forms.ToolStripGripStyle.Hidden;
            this.toolStrip1.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.ValueEmptyButton,
            this.Value1Button,
            this.Value2Button,
            this.Value3Button,
            this.Value4Button,
            this.Value5Button,
            this.Value6Button,
            this.Value7Button,
            this.Value8Button,
            this.Value9Button,
            this.toolStripSeparator1,
            this.hintButton,
            this.solveButton,
            this.evaluateButton,
            this.toolStripSeparator2,
            this.restartButton,
            this.toolStripSeparator3,
            this.aboutButton});
            this.toolStrip1.Location = new System.Drawing.Point(0, 0);
            this.toolStrip1.Name = "toolStrip1";
            this.toolStrip1.Size = new System.Drawing.Size(500, 25);
            this.toolStrip1.TabIndex = 6;
            this.toolStrip1.Text = "toolStrip1";
            // 
            // ValueEmptyButton
            // 
            this.ValueEmptyButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.ValueEmptyButton.Image = ((System.Drawing.Image)(resources.GetObject("ValueEmptyButton.Image")));
            this.ValueEmptyButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.ValueEmptyButton.Name = "ValueEmptyButton";
            this.ValueEmptyButton.Size = new System.Drawing.Size(45, 22);
            this.ValueEmptyButton.Text = "&Empty";
            this.ValueEmptyButton.Click += new System.EventHandler(this.ValueEmptyButton_Click);
            this.ValueEmptyButton.MouseDown += new System.Windows.Forms.MouseEventHandler(this.ValueEmptyButton_MouseDown);
            // 
            // Value1Button
            // 
            this.Value1Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value1Button.Image = ((System.Drawing.Image)(resources.GetObject("Value1Button.Image")));
            this.Value1Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value1Button.Name = "Value1Button";
            this.Value1Button.Size = new System.Drawing.Size(23, 22);
            this.Value1Button.Text = "&1";
            this.Value1Button.Click += new System.EventHandler(this.Value1Button_Click);
            this.Value1Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value1Button_MouseDown);
            // 
            // Value2Button
            // 
            this.Value2Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value2Button.Image = ((System.Drawing.Image)(resources.GetObject("Value2Button.Image")));
            this.Value2Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value2Button.Name = "Value2Button";
            this.Value2Button.Size = new System.Drawing.Size(23, 22);
            this.Value2Button.Text = "&2";
            this.Value2Button.Click += new System.EventHandler(this.Value2Button_Click);
            this.Value2Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value2Button_MouseDown);
            // 
            // Value3Button
            // 
            this.Value3Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value3Button.Image = ((System.Drawing.Image)(resources.GetObject("Value3Button.Image")));
            this.Value3Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value3Button.Name = "Value3Button";
            this.Value3Button.Size = new System.Drawing.Size(23, 22);
            this.Value3Button.Text = "3";
            this.Value3Button.Click += new System.EventHandler(this.Value3Button_Click);
            this.Value3Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value3Button_MouseDown);
            // 
            // Value4Button
            // 
            this.Value4Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value4Button.Image = ((System.Drawing.Image)(resources.GetObject("Value4Button.Image")));
            this.Value4Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value4Button.Name = "Value4Button";
            this.Value4Button.Size = new System.Drawing.Size(23, 22);
            this.Value4Button.Text = "4";
            this.Value4Button.Click += new System.EventHandler(this.Value4Button_Click);
            this.Value4Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value4Button_MouseDown);
            // 
            // Value5Button
            // 
            this.Value5Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value5Button.Image = ((System.Drawing.Image)(resources.GetObject("Value5Button.Image")));
            this.Value5Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value5Button.Name = "Value5Button";
            this.Value5Button.Size = new System.Drawing.Size(23, 22);
            this.Value5Button.Text = "5";
            this.Value5Button.Click += new System.EventHandler(this.Value5Button_Click);
            this.Value5Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value5Button_MouseDown);
            // 
            // Value6Button
            // 
            this.Value6Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value6Button.Image = ((System.Drawing.Image)(resources.GetObject("Value6Button.Image")));
            this.Value6Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value6Button.Name = "Value6Button";
            this.Value6Button.Size = new System.Drawing.Size(23, 22);
            this.Value6Button.Text = "6";
            this.Value6Button.Click += new System.EventHandler(this.Value6Button_Click);
            this.Value6Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value6Button_MouseDown);
            // 
            // Value7Button
            // 
            this.Value7Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value7Button.Image = ((System.Drawing.Image)(resources.GetObject("Value7Button.Image")));
            this.Value7Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value7Button.Name = "Value7Button";
            this.Value7Button.Size = new System.Drawing.Size(23, 22);
            this.Value7Button.Text = "7";
            this.Value7Button.Click += new System.EventHandler(this.Value7Button_Click);
            this.Value7Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value7Button_MouseDown);
            // 
            // Value8Button
            // 
            this.Value8Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value8Button.Image = ((System.Drawing.Image)(resources.GetObject("Value8Button.Image")));
            this.Value8Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value8Button.Name = "Value8Button";
            this.Value8Button.Size = new System.Drawing.Size(23, 22);
            this.Value8Button.Text = "8";
            this.Value8Button.Click += new System.EventHandler(this.Value8Button_Click);
            this.Value8Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value8Button_MouseDown);
            // 
            // Value9Button
            // 
            this.Value9Button.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.Value9Button.Image = ((System.Drawing.Image)(resources.GetObject("Value9Button.Image")));
            this.Value9Button.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.Value9Button.Name = "Value9Button";
            this.Value9Button.Size = new System.Drawing.Size(23, 22);
            this.Value9Button.Text = "9";
            this.Value9Button.Click += new System.EventHandler(this.Value9Button_Click);
            this.Value9Button.MouseDown += new System.Windows.Forms.MouseEventHandler(this.Value9Button_MouseDown);
            // 
            // toolStripSeparator1
            // 
            this.toolStripSeparator1.Name = "toolStripSeparator1";
            this.toolStripSeparator1.Size = new System.Drawing.Size(6, 25);
            // 
            // hintButton
            // 
            this.hintButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.hintButton.Image = ((System.Drawing.Image)(resources.GetObject("hintButton.Image")));
            this.hintButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.hintButton.Name = "hintButton";
            this.hintButton.Size = new System.Drawing.Size(34, 22);
            this.hintButton.Text = "&Hint";
            this.hintButton.Click += new System.EventHandler(this.hintButton_Click);
            // 
            // solveButton
            // 
            this.solveButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.solveButton.Image = ((System.Drawing.Image)(resources.GetObject("solveButton.Image")));
            this.solveButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.solveButton.Name = "solveButton";
            this.solveButton.Size = new System.Drawing.Size(39, 22);
            this.solveButton.Text = "&Solve";
            this.solveButton.Click += new System.EventHandler(this.solveButton_Click);
            // 
            // toolStripSeparator2
            // 
            this.toolStripSeparator2.Name = "toolStripSeparator2";
            this.toolStripSeparator2.Size = new System.Drawing.Size(6, 25);
            // 
            // restartButton
            // 
            this.restartButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.restartButton.Image = ((System.Drawing.Image)(resources.GetObject("restartButton.Image")));
            this.restartButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.restartButton.Name = "restartButton";
            this.restartButton.Size = new System.Drawing.Size(47, 22);
            this.restartButton.Text = "&Restart";
            this.restartButton.Click += new System.EventHandler(this.restartButton_Click);
            // 
            // toolStripSeparator3
            // 
            this.toolStripSeparator3.Name = "toolStripSeparator3";
            this.toolStripSeparator3.Size = new System.Drawing.Size(6, 25);
            // 
            // aboutButton
            // 
            this.aboutButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.aboutButton.Image = ((System.Drawing.Image)(resources.GetObject("aboutButton.Image")));
            this.aboutButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.aboutButton.Name = "aboutButton";
            this.aboutButton.Size = new System.Drawing.Size(44, 19);
            this.aboutButton.Text = "&About";
            this.aboutButton.Click += new System.EventHandler(this.aboutButton_Click);
            // 
            // gamePanel
            // 
            this.gamePanel.AllowDrop = true;
            this.gamePanel.AutoSize = true;
            this.gamePanel.ColumnCount = 3;
            this.gamePanel.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.gamePanel.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33334F));
            this.gamePanel.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33334F));
            this.gamePanel.Controls.Add(this.squarePanel8, 2, 2);
            this.gamePanel.Controls.Add(this.squarePanel7, 1, 2);
            this.gamePanel.Controls.Add(this.squarePanel6, 0, 2);
            this.gamePanel.Controls.Add(this.squarePanel5, 2, 1);
            this.gamePanel.Controls.Add(this.squarePanel4, 1, 1);
            this.gamePanel.Controls.Add(this.squarePanel3, 0, 1);
            this.gamePanel.Controls.Add(this.squarePanel2, 2, 0);
            this.gamePanel.Controls.Add(this.squarePanel1, 1, 0);
            this.gamePanel.Controls.Add(this.squarePanel0, 0, 0);
            this.gamePanel.Dock = System.Windows.Forms.DockStyle.Fill;
            this.gamePanel.GrowStyle = System.Windows.Forms.TableLayoutPanelGrowStyle.FixedSize;
            this.gamePanel.Location = new System.Drawing.Point(0, 25);
            this.gamePanel.Margin = new System.Windows.Forms.Padding(0);
            this.gamePanel.Name = "gamePanel";
            this.gamePanel.RowCount = 3;
            this.gamePanel.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.gamePanel.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.gamePanel.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.gamePanel.Size = new System.Drawing.Size(500, 395);
            this.gamePanel.TabIndex = 7;
            // 
            // squarePanel8
            // 
            this.squarePanel8.AllowDrop = true;
            this.squarePanel8.AutoSize = true;
            this.squarePanel8.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel8.ColumnCount = 3;
            this.squarePanel8.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel8.Location = new System.Drawing.Point(332, 262);
            this.squarePanel8.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel8.Name = "squarePanel8";
            this.squarePanel8.RowCount = 3;
            this.squarePanel8.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel8.Size = new System.Drawing.Size(168, 133);
            this.squarePanel8.TabIndex = 8;
            this.squarePanel8.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel8_CellPaint);
            this.squarePanel8.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel8.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel8.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel7
            // 
            this.squarePanel7.AllowDrop = true;
            this.squarePanel7.AutoSize = true;
            this.squarePanel7.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel7.ColumnCount = 3;
            this.squarePanel7.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel7.Location = new System.Drawing.Point(166, 262);
            this.squarePanel7.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel7.Name = "squarePanel7";
            this.squarePanel7.RowCount = 3;
            this.squarePanel7.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel7.Size = new System.Drawing.Size(166, 133);
            this.squarePanel7.TabIndex = 7;
            this.squarePanel7.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel7_CellPaint);
            this.squarePanel7.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel7.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel7.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel6
            // 
            this.squarePanel6.AllowDrop = true;
            this.squarePanel6.AutoSize = true;
            this.squarePanel6.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel6.ColumnCount = 3;
            this.squarePanel6.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel6.Location = new System.Drawing.Point(0, 262);
            this.squarePanel6.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel6.Name = "squarePanel6";
            this.squarePanel6.RowCount = 3;
            this.squarePanel6.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel6.Size = new System.Drawing.Size(166, 133);
            this.squarePanel6.TabIndex = 6;
            this.squarePanel6.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel6_CellPaint);
            this.squarePanel6.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel6.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel6.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel5
            // 
            this.squarePanel5.AllowDrop = true;
            this.squarePanel5.AutoSize = true;
            this.squarePanel5.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel5.ColumnCount = 3;
            this.squarePanel5.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel5.Location = new System.Drawing.Point(332, 131);
            this.squarePanel5.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel5.Name = "squarePanel5";
            this.squarePanel5.RowCount = 3;
            this.squarePanel5.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel5.Size = new System.Drawing.Size(168, 131);
            this.squarePanel5.TabIndex = 5;
            this.squarePanel5.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel5_CellPaint);
            this.squarePanel5.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel5.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel5.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel4
            // 
            this.squarePanel4.AllowDrop = true;
            this.squarePanel4.AutoSize = true;
            this.squarePanel4.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel4.ColumnCount = 3;
            this.squarePanel4.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel4.Location = new System.Drawing.Point(166, 131);
            this.squarePanel4.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel4.Name = "squarePanel4";
            this.squarePanel4.RowCount = 3;
            this.squarePanel4.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel4.Size = new System.Drawing.Size(166, 131);
            this.squarePanel4.TabIndex = 4;
            this.squarePanel4.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel4_CellPaint);
            this.squarePanel4.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel4.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel4.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel3
            // 
            this.squarePanel3.AllowDrop = true;
            this.squarePanel3.AutoSize = true;
            this.squarePanel3.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel3.ColumnCount = 3;
            this.squarePanel3.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel3.Location = new System.Drawing.Point(0, 131);
            this.squarePanel3.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel3.Name = "squarePanel3";
            this.squarePanel3.RowCount = 3;
            this.squarePanel3.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel3.Size = new System.Drawing.Size(166, 131);
            this.squarePanel3.TabIndex = 3;
            this.squarePanel3.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel3_CellPaint);
            this.squarePanel3.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel3.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel3.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel2
            // 
            this.squarePanel2.AllowDrop = true;
            this.squarePanel2.AutoSize = true;
            this.squarePanel2.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel2.ColumnCount = 3;
            this.squarePanel2.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel2.Location = new System.Drawing.Point(332, 0);
            this.squarePanel2.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel2.Name = "squarePanel2";
            this.squarePanel2.RowCount = 3;
            this.squarePanel2.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel2.Size = new System.Drawing.Size(168, 131);
            this.squarePanel2.TabIndex = 2;
            this.squarePanel2.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel2_CellPaint);
            this.squarePanel2.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel2.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel2.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel1
            // 
            this.squarePanel1.AllowDrop = true;
            this.squarePanel1.AutoSize = true;
            this.squarePanel1.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel1.ColumnCount = 3;
            this.squarePanel1.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel1.Location = new System.Drawing.Point(166, 0);
            this.squarePanel1.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel1.Name = "squarePanel1";
            this.squarePanel1.RowCount = 3;
            this.squarePanel1.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel1.Size = new System.Drawing.Size(166, 131);
            this.squarePanel1.TabIndex = 1;
            this.squarePanel1.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel1_CellPaint);
            this.squarePanel1.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel1.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel1.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // squarePanel0
            // 
            this.squarePanel0.AllowDrop = true;
            this.squarePanel0.AutoSize = true;
            this.squarePanel0.CellBorderStyle = System.Windows.Forms.TableLayoutPanelCellBorderStyle.InsetDouble;
            this.squarePanel0.ColumnCount = 3;
            this.squarePanel0.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.Dock = System.Windows.Forms.DockStyle.Fill;
            this.squarePanel0.Location = new System.Drawing.Point(0, 0);
            this.squarePanel0.Margin = new System.Windows.Forms.Padding(0);
            this.squarePanel0.Name = "squarePanel0";
            this.squarePanel0.RowCount = 3;
            this.squarePanel0.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 33.33333F));
            this.squarePanel0.Size = new System.Drawing.Size(166, 131);
            this.squarePanel0.TabIndex = 0;
            this.squarePanel0.CellPaint += new System.Windows.Forms.TableLayoutCellPaintEventHandler(this.squarePanel0_CellPaint);
            this.squarePanel0.DragDrop += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragDrop);
            this.squarePanel0.DragOver += new System.Windows.Forms.DragEventHandler(this.squarePanel_DragOver);
            this.squarePanel0.MouseClick += new System.Windows.Forms.MouseEventHandler(this.squarePanel_MouseClick);
            // 
            // evaluateButton
            // 
            this.evaluateButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Text;
            this.evaluateButton.Image = ((System.Drawing.Image)(resources.GetObject("evaluateButton.Image")));
            this.evaluateButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.evaluateButton.Name = "evaluateButton";
            this.evaluateButton.Size = new System.Drawing.Size(55, 22);
            this.evaluateButton.Text = "E&valuate";
            this.evaluateButton.Click += new System.EventHandler(this.evaluateButton_Click);
            // 
            // SudokuBoardControl
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.Controls.Add(this.gamePanel);
            this.Controls.Add(this.toolStrip1);
            this.Name = "SudokuBoardControl";
            this.Size = new System.Drawing.Size(500, 420);
            this.toolStrip1.ResumeLayout(false);
            this.toolStrip1.PerformLayout();
            this.gamePanel.ResumeLayout(false);
            this.gamePanel.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.ToolStrip toolStrip1;
        private System.Windows.Forms.ToolStripButton ValueEmptyButton;
        private System.Windows.Forms.ToolStripButton Value1Button;
        private System.Windows.Forms.ToolStripButton Value2Button;
        private System.Windows.Forms.ToolStripButton Value3Button;
        private System.Windows.Forms.ToolStripButton Value4Button;
        private System.Windows.Forms.ToolStripButton Value5Button;
        private System.Windows.Forms.ToolStripButton Value6Button;
        private System.Windows.Forms.ToolStripButton Value7Button;
        private System.Windows.Forms.ToolStripButton Value8Button;
        private System.Windows.Forms.ToolStripButton Value9Button;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator1;
        private System.Windows.Forms.ToolStripButton hintButton;
        private System.Windows.Forms.TableLayoutPanel gamePanel;
        private System.Windows.Forms.TableLayoutPanel squarePanel8;
        private System.Windows.Forms.TableLayoutPanel squarePanel7;
        private System.Windows.Forms.TableLayoutPanel squarePanel6;
        private System.Windows.Forms.TableLayoutPanel squarePanel5;
        private System.Windows.Forms.TableLayoutPanel squarePanel4;
        private System.Windows.Forms.TableLayoutPanel squarePanel3;
        private System.Windows.Forms.TableLayoutPanel squarePanel2;
        private System.Windows.Forms.TableLayoutPanel squarePanel1;
        private System.Windows.Forms.TableLayoutPanel squarePanel0;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator2;
        private System.Windows.Forms.ToolStripButton restartButton;
        private System.Windows.Forms.ToolStripButton aboutButton;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator3;
        private System.Windows.Forms.ToolStripButton solveButton;
        private System.Windows.Forms.ToolStripButton evaluateButton;
    }
}
