using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Sudoku
{
    public enum SudokuLevel
    {
        Unclassified,
        Easy,
        Medium,
        Hard,
        VeryHard,
        Evil,
        DarkEvil,
        BlackHole,

        // keep last
        Max = 8
    }
}
