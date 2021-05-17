using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Sudoku
{
    public struct SudokuValueSet
    {
        private static ushort MaskShifted(byte value) => (ushort)(1 << (value + 3));

        // count:bits 0 to 3 (total value no more than 9 so 4 bits are enough)
        // mask: bits 4 to 12 (9 bits)
        private ushort CountAndMask_;

        public byte Count => (byte)(CountAndMask_ & 0xF);
        public ushort Mask => (ushort)(CountAndMask_ >> 4);

        public SudokuValueSet(SudokuValueSet value)
        {
            this.CountAndMask_ = value.CountAndMask_;
        }
        private SudokuValueSet(byte count, short mask)
        {
            this.CountAndMask_ = (ushort)(count | (mask << 4));
        }

        public bool AddValue(byte value)
        {
            if (value < 1 || value > 9)
                throw new ArgumentOutOfRangeException();

            ushort maskShifted = MaskShifted(value);
            if ((CountAndMask_ & maskShifted) == 0)
            {
                ++CountAndMask_;
                CountAndMask_ = (ushort)(maskShifted | CountAndMask_);
                return true;
            }
            else
            {
                return false;
            }
        }

        public bool ContainsValue(byte value)
        {
            if (value < 1 || value > 9)
                throw new ArgumentOutOfRangeException();

            ushort maskShifted = MaskShifted(value);
            return ((CountAndMask_ & maskShifted) != 0);
        }

        public bool RemoveValue(byte value)
        {
            if (value < 1 || value > 9)
                throw new ArgumentOutOfRangeException();

            ushort maskShifted = MaskShifted(value);
            if ((CountAndMask_ & maskShifted) == 0)
            {
                return false;
            }
            else
            {
                --CountAndMask_;
                CountAndMask_ = (ushort)(CountAndMask_ & ~maskShifted);
                return true;
            }
        }

        public IEnumerable<byte> EnumValues()
        {
            ushort maskShifted = MaskShifted(1);
            for(byte num = 1; num <= 9; num++)
            {
                if ((CountAndMask_ & maskShifted) != 0)
                {
                    yield return num;
                }
                maskShifted <<= 1;
            }
        }

        public IEnumerable<byte> EnumMissingValues()
        {
            ushort maskShifted = MaskShifted(1);
            for (byte num = 1; num <= 9; num++)
            {
                if ((CountAndMask_ & maskShifted) == 0)
                {
                    yield return num;
                }
                maskShifted <<= 1;
            }
        }

        public int Combined => CombineValueSet(this.Mask);
        
        public override string ToString()
        {
            return Combined.ToString();
        }

        public override int GetHashCode()
        {
            return CountAndMask_;
        }

        public override bool Equals(object obj)
        {
            if (obj == null || obj.GetType() != typeof(SudokuValueSet))
                return false;
            
            return this.CountAndMask_ == ((SudokuValueSet)obj).CountAndMask_;
        }

        public byte this[int valueIndex]
        {
            get
            {
                if (valueIndex < 0 || valueIndex >= Count)
                    throw new ArgumentOutOfRangeException(nameof(valueIndex));

                for (byte value = 1; value <= 9; value++)
                {
                    if (ContainsValue(value))
                    {
                        if (valueIndex == 0)
                            return value;
                        else
                            valueIndex--;
                    }
                }

                throw new InvalidOperationException();
            }
        }

        public static SudokuValueSet operator!(SudokuValueSet value)
        {
            byte count = (byte)(9 - value.Count);
            short mask = (short)(0x1FF & ~value.Mask);
            return new SudokuValueSet(count, mask);
        }

        /// <summary>
        /// create an int out of the allowed values combined together, from 1 to 9: example 146
        /// </summary>
        public static int CombineValueSet(ushort valuesMask)
        {
            if (valuesMask == 0)
                return 0;

            int combinedValue = 0;
            for (byte value = 1; value <= 9; value++)
            {
                int mask = MaskOf(value);
                if ((valuesMask & mask) != 0)
                {
                    combinedValue = combinedValue * 10 + value;
                }
            }

            return combinedValue;
        }

        public static ushort MaskOf(byte value) => (ushort)(1 << (value - 1));

        public void AddMask(ushort valueMask)
        {
            for (byte value = 1; value <= 9; value++)
            {
                ushort mask = MaskOf(value);
                if ((valueMask & mask) == 0)
                    continue;

                AddValue(value);
            }
        }

        public static SudokuValueSet Empty => new SudokuValueSet();
        public static SudokuValueSet All => new SudokuValueSet(9, 0x1FF);
    }
}
