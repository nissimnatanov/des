using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.Text;
using System.Xml;

namespace Sudoku
{
    public class SudokuBoard
    {
        #region static info

        public static readonly int[][] SudokuRowIndicies;
        public static readonly int[][] SudokuColIndicies;
        public static readonly int[][] SudokuSquareIndicies;
        public static readonly int[][] SudokuRelatedIndicies;
        public static readonly int[] SudokuIndexToSquare;
        public static readonly int[] SudokuIndexToCellIndexSquare;

        static SudokuBoard()
        {
            SudokuRowIndicies = new int[9][];
            SudokuColIndicies = new int[9][];
            SudokuSquareIndicies = new int[9][];
            SudokuRelatedIndicies = new int[81][];
            SudokuIndexToSquare = new int[81];
            SudokuIndexToCellIndexSquare = new int[81];

            for (int i = 0; i < 9; i++)
            {
                SudokuRowIndicies[i] = new int[9];
                SudokuColIndicies[i] = new int[9];
                SudokuSquareIndicies[i] = new int[9];

                for (int j = 0; j < 9; j++)
                {
                    // set row indicies
                    // 0 .. 8; 9 .. 17; ...
                    SudokuRowIndicies[i][j] = i * 9 + j;

                    // set column indicies
                    // { 0, 9, 18 ...} {1, 10, 19 ...} ...
                    SudokuColIndicies[i][j] = i + j * 9;

                    // set square indicies
                    // { 0, 1, 2, 9, 10, 11, ...}
                    int index = 0;
                    index += (i / 3) * 27; // note that i / 3 gives int value so the expression cannot be replaced with i * 9
                    index += (j / 3) * 9;
                    index += ((i % 3) * 3 + (j % 3));

                    SudokuSquareIndicies[i][j] = index;

                    // set map between index to square
                    SudokuIndexToSquare[index] = i;
                    SudokuIndexToCellIndexSquare[index] = j;
                }
            }

            // get the related indicies for each cell
            for (int i = 0; i < 81; i++)
            {
                List<int> related = new List<int>();
                foreach (var relatedIndex in SudokuRowIndicies[i / 9])
                {
                    if (relatedIndex == i)
                        continue;
                    if (!related.Contains(relatedIndex))
                        related.Add(relatedIndex);
                }
                foreach (var relatedIndex in SudokuColIndicies[i % 9])
                {
                    if (relatedIndex == i)
                        continue;
                    if (!related.Contains(relatedIndex))
                        related.Add(relatedIndex);
                }
                foreach (var relatedIndex in SudokuSquareIndicies[GetSquareFromIndex(i)])
                {
                    if (relatedIndex == i)
                        continue;
                    if (!related.Contains(relatedIndex))
                        related.Add(relatedIndex);
                }

                if (related.Count != 20)
                    throw new InvalidOperationException();
                SudokuRelatedIndicies[i] = related.ToArray();
            }
        }

        public static int GetSquareFromIndex(int index)
        {
            return SudokuIndexToSquare[index];
        }

        public static void GetSquareFromIndex(int index, out int square, out int cellIndexInSquare)
        {
            square = SudokuIndexToSquare[index];
            cellIndexInSquare = SudokuIndexToCellIndexSquare[index];
        }

        public static int GetIndexFromCoordinates(int row, int col)
        {
            return row * 9 + col;
        }

        public static int GetIndexFromSquare(int squareIndex, int cellIndex)
        {
            return SudokuSquareIndicies[squareIndex][cellIndex];
        }

        #endregion

        public class ValueChangeEventArgs : EventArgs
        {
            public ValueChangeEventArgs()
            {
            }
        }

        public delegate void ValueChangeHandler(object sender, ValueChangeEventArgs e);

        [Flags]
        public enum CellState : byte
        {
            None = 0,

            ReadOnly = 0x1,
            Invalid = 0x2,
        }

        [Serializable]
        struct Cell
        {
            private byte Value_;
            private CellState State_;
            SudokuValueSet DisallowedValues_;

            /// <summary>
            /// first 10 - disallowedvalues
            /// </summary>

            public bool IsReadOnly
            {
                get
                {
                    return (State_ & CellState.ReadOnly) != 0;
                }
                internal set
                {
                    if (value)
                    {
                        State_ |= CellState.ReadOnly;
                    }
                    else
                    {
                        State_ &= ~CellState.ReadOnly;
                    }
                }
            }

            public bool IsValid
            {
                get
                {
                    if (IsEmpty)
                        return AllowedValues.Count > 0;
                    else
                        return (State_ & CellState.Invalid) == 0;
                }
                internal set
                {
                    if (!value)
                    {
                        if (IsEmpty)
                            throw new InvalidOperationException("cannot set empty as invalid");
                        State_ |= CellState.Invalid;
                    }
                    else
                    {
                        State_ &= ~CellState.Invalid;
                    }
                }
            }

            public bool IsEmpty => Value_ == 0;

            public SudokuValueSet DisallowedValues
            {
                get
                {
                    if (!IsEmpty)
                        throw new InvalidOperationException("do not use for populated cells");
                    return DisallowedValues_;
                }
            }

            public SudokuValueSet AllowedValues => !DisallowedValues_;

            public byte Value
            {
                get
                {
                    return Value_;
                }
            }

            internal void SetValue(byte value, bool asReadOnly)
            {
                if (value < 1 || value > 9)
                    throw new ArgumentOutOfRangeException();

                IsReadOnly = asReadOnly;

                Value_ = value;
            }

            internal void DisallowValue(byte value)
            {
                if (!IsEmpty)
                    throw new InvalidOperationException("already has a value");

                DisallowedValues_.AddValue(value);
            }

            internal void DisallowValueMask(ushort valueMask)
            {
                if (!IsEmpty)
                    throw new InvalidOperationException("already has a value");

                DisallowedValues_.AddMask(valueMask);
            }

            internal bool IsAllowedValue(byte value)
            {
                if (!IsEmpty)
                    throw new InvalidOperationException("already has a value");

                return !DisallowedValues_.ContainsValue(value);
            }

            internal void Reset()
            {
                Value_ = 0;
                State_ = CellState.None;
                DisallowedValues_ = SudokuValueSet.Empty;
            }

            internal void ResetAllowedValues()
            {
                if (!IsEmpty)
                    throw new InvalidOperationException("already has a value");

                DisallowedValues_ = SudokuValueSet.Empty;
            }
        }

        Cell[] m_values;
        SudokuValueSet[] m_columnValues;
        SudokuValueSet[] m_rowValues;
        SudokuValueSet[] m_squareValues;
        ValueChangeHandler m_onChange;
        int m_onChangeSuspendCount;
        bool m_editMode;
        bool m_readOnly;
        bool m_isValid;
        SudokuLevel m_level;
        int m_solutionWeight;
        SudokuValueSet[] m_userDisallowed;
        int[] m_valueStats;

        public SudokuBoard(bool editMode)
        {
            m_values = new Cell[81];
            m_editMode = editMode;
            m_level = SudokuLevel.Unclassified;
            m_isValid = true;
            m_columnValues = new SudokuValueSet[9];
            m_rowValues = new SudokuValueSet[9];
            m_squareValues = new SudokuValueSet[9];
            m_userDisallowed = new SudokuValueSet[81];
            m_valueStats = new int[10]; // first is reserved to count free cells
            m_valueStats[0] = 81;

            CheckIntegrity();
        }

        public SudokuBoard(SudokuBoard boardToCopy)
        {
            m_values = (Cell[])boardToCopy.m_values.Clone();
            m_columnValues = (SudokuValueSet[])boardToCopy.m_columnValues.Clone();
            m_rowValues = (SudokuValueSet[])boardToCopy.m_rowValues.Clone();
            m_squareValues = (SudokuValueSet[])boardToCopy.m_squareValues.Clone();
            m_editMode = boardToCopy.m_editMode;
            m_level = boardToCopy.m_level;
            m_readOnly = boardToCopy.m_readOnly;
            m_isValid = boardToCopy.m_isValid;
            m_userDisallowed = (SudokuValueSet[])boardToCopy.m_userDisallowed.Clone();
            m_valueStats = (int[])boardToCopy.m_valueStats.Clone();

            CheckIntegrity();
        }

        public SudokuLevel GameLevel
        {
            get
            {
                return m_level;
            }
        }


        public int SolutionWeight
        {
            get
            {
                return m_solutionWeight;
            }
        }

        internal void SetLevel(SudokuLevel level, int weight)
        {
            if (!EditMode)
                throw new InvalidOperationException("Level details can be set in Edit mode only");

            m_level = level;
            m_solutionWeight = weight;
        }

        public int GameId
        {
            get
            {
                return 0;
            }
        }

        public bool EditMode
        {
            get
            {
                return m_editMode;
            }
            set
            {
                m_editMode = value;
            }
        }

        public bool ReadOnly
        {
            get
            {
                return m_readOnly;
            }
            set
            {
                m_readOnly = value;
            }
        }

        public event ValueChangeHandler OnChange
        {
            add
            {
                m_onChange += value;
            }
            remove
            {
                m_onChange -= value;
            }
        }

        public bool IsValid => m_isValid;

        public int FreeCellsCount => m_valueStats[0];

        public int GetValueCount(byte value) => m_valueStats[value];

        public void OnChanged()
        {
            if (m_onChangeSuspendCount == 0)
                m_onChange?.Invoke(this, new ValueChangeEventArgs());
        }

        public byte this[int index]
        {
            get
            {
                CheckIndex(index);

                return m_values[index].Value;
            }
        }

        public void SuspendEvents()
        {
            ++m_onChangeSuspendCount;
        }

        public void ResumeEvents()
        {
            --m_onChangeSuspendCount;
            if (m_onChangeSuspendCount < 0)
                throw new InvalidOperationException();

            OnChanged();
        }

        public void SetValue(int index, byte value, bool isReadOnly)
        {
            CheckIndex(index);
            CheckValue(value);

            if (m_readOnly)
                throw new ArgumentException("This game is read-only");

            if (!m_editMode && (m_values[index].IsReadOnly || isReadOnly))
                throw new ArgumentException("Edit mode is not allowed");

            byte previousValue = m_values[index].Value;
            if (previousValue == value)
            {
                m_values[index].IsReadOnly = isReadOnly;
                return;
            }
            --m_valueStats[previousValue];
            ++m_valueStats[value];

            if (value == 0)
            {
                if (isReadOnly)
                    throw new ArgumentException("Empty cell cannot be read-only");

                m_values[index].Reset();
            }
            else
            {
                m_values[index].SetValue(value, isReadOnly);
            }

            // recalculate state of related row/column/square
            int row = index / 9;
            int column = index % 9;
            int square = GetSquareFromIndex(index);

            if (!m_isValid)
            {
                // fully recalculate the state
                m_isValid = true;

                RecalculateAggregateValues(m_rowValues, row, SudokuRowIndicies[row]);
                RecalculateAggregateValues(m_columnValues, column, SudokuColIndicies[column]);
                RecalculateAggregateValues(m_squareValues, square, SudokuSquareIndicies[square]);

                m_isValid = RecalculateAllowedValuesFromAggregates();
            }
            else
            {
                UpdateAggregateValues(m_rowValues, SudokuRowIndicies[row], index, row, previousValue);
                UpdateAggregateValues(m_columnValues, SudokuColIndicies[column], index, column, previousValue);
                UpdateAggregateValues(m_squareValues, SudokuSquareIndicies[square], index, square, previousValue);

                m_isValid = UpdateAllowedValuesFromAggregates(index, previousValue);
            }
            OnChanged();

            CheckIntegrity();
        }

        internal void DisallowValue(int index, byte value)
        {
            this.m_values[index].DisallowValue(value);
            m_userDisallowed[index].AddValue(value);

            CheckIntegrity();
        }

        internal void DisallowValueMask(int index, ushort valueMask)
        {
            m_values[index].DisallowValueMask(valueMask);
            m_userDisallowed[index].AddMask(valueMask);

            CheckIntegrity();
        }

        internal bool IsAllowedValue(int index, byte value)
        {
            return m_values[index].IsAllowedValue(value);
        }
        internal SudokuValueSet GetAllowedValues(int index)
        {
            return m_values[index].AllowedValues;
        }

        internal SudokuValueSet GetValuesInRow(int row)
        {
            return m_rowValues[row];
        }
        internal SudokuValueSet GetValuesInColumn(int column)
        {
            return m_columnValues[column];
        }
        internal SudokuValueSet GetValuesInSquare(int square)
        {
            return m_squareValues[square];
        }

        public bool IsDefault(int index)
        {
            CheckIndex(index);
            return m_values[index].IsReadOnly;
        }

        public bool IsReadOnly(int index)
        {
            return IsDefault(index) || m_readOnly;
        }

        public bool IsValidCell(int index)
        {
            CheckIndex(index);
            return m_values[index].IsValid;
        }

        #region validation

        void CheckIndex(int index)
        {
            if (index < 0 || index > 80)
                throw new ArgumentOutOfRangeException("index");
        }

        void CheckValue(byte value)
        {
            // 0 means empty
            // 1 - 9 is the value
            if (value < 0 || value > 9)
                throw new ArgumentOutOfRangeException("value");
        }

        private bool DisallowValue(int[] indices, int modifiedIndex, byte value)
        {
            bool valid = true;
            for (int i = 0; i < indices.Length; i++)
            {
                int target = indices[i];
                if (target == modifiedIndex)
                    throw new InvalidOperationException("should not have modified index in related indicies");

                if (m_values[target].IsEmpty)
                {
                    m_values[target].DisallowValue(value);
                }
                else if (m_values[target].Value == value)
                {
                    // collision detected with new value
                    valid = false;

                    if (m_values[target].IsReadOnly)
                    {
                        m_values[modifiedIndex].IsValid = false;
                        if (m_values[modifiedIndex].IsReadOnly)
                        {
                            // both read-only? mark both as invalid
                            m_values[target].IsValid = false;
                        }
                    }
                    else
                    {
                        m_values[target].IsValid = false;
                        if (!m_values[modifiedIndex].IsReadOnly)
                        {
                            m_values[modifiedIndex].IsValid = false;
                        }
                    }
                }
            }

            return valid;
        }

        bool Validate(int[] indices)
        {
            bool valid = true;
            if (indices.Length != 9)
                throw new InvalidOperationException();

            ushort seenValuesMask = 0;
            ushort dupeValuesMask = 0;

            // pass 1 - get mask of existing values
            for (int i = 0; i < 9; i++)
            {
                int target = indices[i];
                byte value = m_values[target].Value;
                if (value == 0)
                    continue;

                ushort mask = SudokuValueSet.MaskOf(value);
                if ((mask & seenValuesMask) == 0)
                {
                    seenValuesMask |= mask;
                }
                else
                {
                    dupeValuesMask |= mask;
                }
            }

            if (dupeValuesMask != 0)
            {
                valid = false;

                // pass 2 - mark dupe cells as invalid
                bool atLeastOneInvalid = false;
                bool checkReadOnly = true;
                do
                {
                    for (int i = 0; i < 9; i++)
                    {
                        int target = indices[i];
                        if (m_values[target].Value == 0 || (checkReadOnly && m_values[target].IsReadOnly))
                            continue;

                        ushort mask = SudokuValueSet.MaskOf(m_values[target].Value);
                        if ((mask & dupeValuesMask) != 0)
                        {
                            m_values[target].IsValid = false;
                            atLeastOneInvalid = true;
                        }
                    }
                    if (!checkReadOnly)
                        break;// no third retry
                    else
                        checkReadOnly = false; // if we have to go second time
                }
                while (!atLeastOneInvalid);

                if (checkReadOnly && !atLeastOneInvalid)
                    throw new InvalidOperationException("should never happen");
            }

            return valid;
        }

        void UpdateAllowedValuesFromAggregates(int index)
        {
            if (m_values[index].IsEmpty)
            {
                int row = index / 9;
                int col = index % 9;
                int square = GetSquareFromIndex(index);

                ushort seenValuesMask = (ushort)(m_rowValues[row].Mask | m_columnValues[col].Mask | m_squareValues[square].Mask);
                seenValuesMask |= m_userDisallowed[index].Mask;

                m_values[index].ResetAllowedValues();
                m_values[index].DisallowValueMask(seenValuesMask);
            }
        }

        void UpdateAllowedValuesFromAggregates(int[] indices)
        {
            for (int i = 0; i < indices.Length; i++)
            {
                UpdateAllowedValuesFromAggregates(indices[i]);
            }
        }

        private bool RecalculateAllowedValuesFromAggregates()
        {
            // first, update all allowed values for empty cells and reset Valid flag for non-empty ones
            for (int i = 0; i < 81; i++)
            {
                m_values[i].IsValid = true;
                UpdateAllowedValuesFromAggregates(i);
            }

            bool valid = true;
            // revalidate all
            for (int i = 0; i < 9; i++)
            {
                valid &= Validate(SudokuRowIndicies[i]);
                valid &= Validate(SudokuColIndicies[i]);
                valid &= Validate(SudokuSquareIndicies[i]);
            }

            return valid;
        }

        bool UpdateAllowedValuesFromAggregates(int modifiedIndex, byte previousValue)
        {
            if (!m_isValid)
                throw new InvalidOperationException("Use RecalculateAllowedValuesFromAggregates instead");

            bool valid = true;
            if (previousValue == 0)
            {
                // previous state is healthy and new value is added in empty spot - safe to perform delta update only
                byte newValue = m_values[modifiedIndex].Value;
                if (newValue != 0)
                {
                    valid &= DisallowValue(SudokuRelatedIndicies[modifiedIndex], modifiedIndex, newValue);
                }
            }
            else
            {
                UpdateAllowedValuesFromAggregates(SudokuRelatedIndicies[modifiedIndex]);

                // update for self as well (if empty)
                UpdateAllowedValuesFromAggregates(modifiedIndex);

                int row = modifiedIndex / 9;
                int col = modifiedIndex % 9;
                int square = GetSquareFromIndex(modifiedIndex);

                valid &= Validate(SudokuRowIndicies[row]);
                valid &= Validate(SudokuColIndicies[col]);
                valid &= Validate(SudokuSquareIndicies[square]);
            }

            return valid;
        }

        void UpdateAggregateValues(SudokuValueSet[] stats, int[] indicies, int modifiedIndex, int sequenceIndex, byte previousValue)
        {
            if (!m_isValid)
                throw new InvalidOperationException("Use RecalculateAggregateValues instead");

            // previous state is healthy  - safe to perform delta update only
            if (previousValue != 0)
            {
                stats[sequenceIndex].RemoveValue(previousValue);
            }
            byte newValue = m_values[modifiedIndex].Value;
            if (newValue != 0)
            {
                stats[sequenceIndex].AddValue(newValue);
            }
        }

        void RecalculateAggregateValues()
        {
            for (int sequenceIndex = 0; sequenceIndex < 9; sequenceIndex++)
            {
                RecalculateAggregateValues(m_squareValues, sequenceIndex, SudokuBoard.SudokuSquareIndicies[sequenceIndex]);
                RecalculateAggregateValues(m_rowValues, sequenceIndex, SudokuBoard.SudokuRowIndicies[sequenceIndex]);
                RecalculateAggregateValues(m_columnValues, sequenceIndex, SudokuBoard.SudokuColIndicies[sequenceIndex]);
            }
        }

        void RecalculateAggregateValues(SudokuValueSet[] stats, int sequenceIndex, int[] indicies)
        {
            // need to recalc the stats from scratch for the given sequence as there might be dupes there before and/or now...
            SudokuValueSet values = new SudokuValueSet();
            for (int i = 0; i < 9; i++)
            {
                byte value = m_values[indicies[i]].Value;
                if (value != 0)
                    values.AddValue(value);
            }

            stats[sequenceIndex] = values;
        }

        #endregion

        public void RestartGame()
        {
            for (int i = 0; i < m_valueStats.Length; i++)
            {
                m_valueStats[i] = 0;
            }

            for (int i = 0; i < m_values.Length; i++)
            {
                if (!m_values[i].IsReadOnly)
                {
                    m_values[i].Reset();
                }

                ++m_valueStats[m_values[i].Value];
            }

            for (int i = 0; i < m_userDisallowed.Length; i++)
            {
                m_userDisallowed[i] = new SudokuValueSet();
            }

            RecalculateAggregateValues();
            m_isValid = RecalculateAllowedValuesFromAggregates();

            OnChanged();
            CheckIntegrity();
        }

        public bool CanRestart
        {
            get
            {
                for (int i = 0; i < m_values.Length; i++)
                {
                    if (!m_values[i].IsReadOnly && m_values[i].Value != 0)
                        return true;
                }

                return false;
            }
        }

        public bool IsSolved()
        {
            for (int i = 0; i < m_values.Length; i++)
                if (m_values[i].Value == 0 || !m_values[i].IsValid)
                    return false;

            return true;
        }

        const string GameLevelAttribute = "gameLevel";
        const string GameIdAttribute = "gameId";
        const string BoardValuesElement = "BoardValues";
        const string ValueStatesElement = "ValueStates";

        public void SaveToFile(string fileName)
        {
            XmlDocument doc = new XmlDocument();

            XmlElement boardElm = doc.CreateElement("SudokuBoard");
            doc.AppendChild(boardElm);

            XmlAttribute levelAttr = doc.CreateAttribute(GameLevelAttribute);
            levelAttr.Value = GameLevel.ToString();
            boardElm.Attributes.Append(levelAttr);

            XmlAttribute idAttr = doc.CreateAttribute(GameIdAttribute);
            idAttr.Value = GameId.ToString();
            boardElm.Attributes.Append(idAttr);

            XmlElement boardValuesElm = doc.CreateElement(BoardValuesElement);
            XmlElement valueStatesElm = doc.CreateElement(ValueStatesElement);

            StringBuilder boardValues = new StringBuilder(81);
            StringBuilder valueStates = new StringBuilder(81);

            for (int i = 0; i < 81; i++)
            {
                char val = (char)('0' + (int)m_values[i].Value);
                boardValues.Append(val);

                char state = m_values[i].IsReadOnly ? 'r' : 'w';
                valueStates.Append(state);
            }

            boardValuesElm.InnerText = boardValues.ToString();
            valueStatesElm.InnerText = valueStates.ToString();

            boardElm.AppendChild(boardValuesElm);
            boardElm.AppendChild(valueStatesElm);

            doc.Save(fileName);
        }

        public static SudokuBoard LoadFromFile(string fileName)
        {
            XmlDocument doc = new XmlDocument();
            doc.Load(fileName);

            XmlElement boardElm = doc.DocumentElement;

            string temp = boardElm.GetAttribute(GameLevelAttribute);
            SudokuLevel level = (SudokuLevel)Enum.Parse(typeof(SudokuLevel), temp, true);

            temp = boardElm.GetAttribute(GameIdAttribute);
            int id = int.Parse(temp);

            XmlElement boardValuesElm = null;
            XmlElement valueStatesElm = null;

            for (int i = 0; i < boardElm.ChildNodes.Count; i++)
            {
                XmlNode node = boardElm.ChildNodes[i];
                if (node.NodeType != XmlNodeType.Element)
                {
                    if (node.NodeType == XmlNodeType.Comment)
                        continue;

                    throw new ArgumentException("Unsupported XML node type: " + node.NodeType);
                }

                if (string.Compare(node.LocalName, BoardValuesElement, true) == 0)
                {
                    if (boardValuesElm != null)
                        throw new ArgumentException(BoardValuesElement + " element cannot appear twice");
                    boardValuesElm = (XmlElement)node;
                }
                else if (string.Compare(node.LocalName, ValueStatesElement, true) == 0)
                {
                    if (valueStatesElm != null)
                        throw new ArgumentException(ValueStatesElement + " element cannot appear twice");

                    valueStatesElm = (XmlElement)node;
                }
                else
                    throw new ArgumentException("Unsupported XML node: " + node.LocalName);
            }

            if (boardValuesElm == null)
                throw new ArgumentException(BoardValuesElement + " element is missing");
            string boardValues = boardValuesElm.InnerText.Trim();
            if (boardValues.Length != 81)
                throw new ArgumentException("There should be exactly 81 characters in " + BoardValuesElement + " element");

            string valueStates;
            if (valueStatesElm != null)
            {
                valueStates = valueStatesElm.InnerText.Trim();
                if (valueStates.Length != 81)
                    throw new ArgumentException("There should be exactly 81 characters in " + ValueStatesElement + " element");
            }
            else
            {
                valueStates = null;
            }

            SudokuBoard board = new SudokuBoard(true);

            for (int i = 0; i < 81; i++)
            {
                char c = boardValues[i];
                byte val = (byte)(c - '0');

                if (val > 9)
                    throw new ArgumentException("Invalid character in " + BoardValuesElement + " element");

                bool isReadOnly;
                if (valueStates != null)
                {
                    c = valueStates[i];
                    if (c == 'r' || c == 'R')
                        isReadOnly = true;
                    else
                        isReadOnly = false;
                }
                else
                {
                    isReadOnly = (val != 0);
                }

                board.SetValue(i, val, isReadOnly);
            }

            // solve to set the level and weight
            var solution = SudokuSolver.Solve(new SudokuBoard(board));
            board.SetLevel(solution.GameLevel, solution.TotalWeight);

            board.EditMode = false;

            board.CheckIntegrity();
            return board;
        }

        public static SudokuBoard LoadFromText(string text)
        {
            SudokuBoard board = new SudokuBoard(true);

            int index = 0;

            foreach (char c in text)
            {
                if (c < '0' || c > '9')
                    continue;

                byte val = (byte)(c - '0');

                if (index >= 81)
                    throw new ArgumentException("Too many 0-9 digits detected in the input");
                board.SetValue(index++, val, (c != '0'));
            }

            if (index < 81)
                throw new ArgumentException($"There are only {index} digits found, need 81");
            board.EditMode = false;

            board.CheckIntegrity();
            return board;
        }

        public override string ToString()
        {
            StringBuilder sb = new StringBuilder();
            AppendBoardValues(sb);
            AppendRowValues(sb);
            AppendColumnValues(sb);
            AppendSquareValues(sb);
            return sb.ToString();
        }

        public void AppendBoardValues(StringBuilder sb)
        {
            sb.AppendLine();
            sb.AppendLine("╔═══════╦═══════╦═══════╗");

            for (int row = 0; row < 9; row++)
            {
                if (row != 0 && row % 3 == 0)
                    sb.AppendLine("╠═══════╬═══════╬═══════╣");

                for (int col = 0; col < 9; col++)
                {
                    if (col % 3 == 0)
                        sb.Append("║ ");
                    int i = GetIndexFromCoordinates(row, col);
                    sb.AppendFormat("{0} ", m_values[i].Value);
                }
                sb.AppendLine("║");
            }
            sb.AppendLine("╚═══════╩═══════╩═══════╝");
        }

        public void AppendRowValues(StringBuilder sb)
        {
            sb.AppendLine();
            sb.AppendLine("Row values:");
            for (int i = 0; i < 9; i++)
            {
                sb.AppendFormat("  [r{0}] {1}", i, m_rowValues[i]);
                sb.AppendLine();
            }
        }

        public void AppendColumnValues(StringBuilder sb)
        {
            sb.AppendLine();
            sb.AppendLine("Column values:");
            for (int i = 0; i < 9; i++)
            {
                sb.AppendFormat("  [c{0}] {1}", i, m_columnValues[i]);
                sb.AppendLine();
            }
        }
        public void AppendSquareValues(StringBuilder sb)
        {
            sb.AppendLine();
            sb.AppendLine("Square values:");
            for (int i = 0; i < 9; i++)
            {
                sb.AppendFormat("  [s{0}] {1}", i, m_squareValues[i]);
                sb.AppendLine();
            }
        }

        [Conditional("INTEGRITYCHECK")]
        void CheckIntegrity()
        {
            int[] valueStats = new int[10];
            for (int i = 0; i < 81; i++)
            {
                byte value = this[i];
                ++valueStats[value];

                // ensure value does not appear
                if (value != 0)
                {
                    // check this value is disallowed in other places
                    int[] relatedIndicies = SudokuRelatedIndicies[i];
                    foreach (var related in relatedIndicies)
                    {
                        if (m_values[related].IsEmpty)
                        {
                            System.Diagnostics.Debug.Assert(!m_values[related].IsAllowedValue(value));
                        }
                        else if (m_values[related].Value == value)
                        {
                            // ensure one of them is marked as invalid
                            if (!m_values[related].IsReadOnly)
                                System.Diagnostics.Debug.Assert(!m_values[related].IsValid);
                            if (!m_values[i].IsReadOnly)
                                System.Diagnostics.Debug.Assert(!m_values[i].IsValid);
                            if (m_values[related].IsReadOnly && m_values[i].IsReadOnly)
                            {
                                System.Diagnostics.Debug.Assert(!m_values[related].IsValid);
                                System.Diagnostics.Debug.Assert(!m_values[i].IsValid);
                            }
                        }
                    }
                }
                else
                {
                    // check that disallowed values are a union of row/column/square
                    int row = i / 9;
                    int col = i % 9;
                    int square = GetSquareFromIndex(i);

                    SudokuValueSet disallowedValuesExpected = m_rowValues[row];
                    disallowedValuesExpected.AddMask(m_columnValues[col].Mask);
                    disallowedValuesExpected.AddMask(m_squareValues[square].Mask);

                    disallowedValuesExpected.AddMask(m_userDisallowed[i].Mask);

                    Debug.Assert(m_values[i].DisallowedValues.Mask == disallowedValuesExpected.Mask);
                }
            }

            for (int i = 0; i < valueStats.Length; i++)
            {
                Debug.Assert(valueStats[i] == m_valueStats[i]);
            }

            for (int sequenceIndex = 0; sequenceIndex < 9; sequenceIndex++)
            {
                CheckAggregateValues(m_squareValues, sequenceIndex, SudokuBoard.SudokuSquareIndicies[sequenceIndex]);
                CheckAggregateValues(m_rowValues, sequenceIndex, SudokuBoard.SudokuRowIndicies[sequenceIndex]);
                CheckAggregateValues(m_columnValues, sequenceIndex, SudokuBoard.SudokuColIndicies[sequenceIndex]);
            }
        }

        void CheckAggregateValues(SudokuValueSet[] stats, int sequenceIndex, int[] indicies)
        {
            // need to recalc the stats from scratch for the given sequence as there might be dupes there before and/or now...
            SudokuValueSet values = new SudokuValueSet();
            for (int i = 0; i < 9; i++)
            {
                byte value = m_values[indicies[i]].Value;
                if (value != 0)
                    values.AddValue(value);
            }

            System.Diagnostics.Debug.Assert(stats[sequenceIndex].Combined == values.Combined);
        }
    }
}
