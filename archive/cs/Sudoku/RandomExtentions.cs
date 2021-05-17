using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Sudoku
{
    static class RandomExtentions
    {        
        public static void Shuffle<T>(this Random r, IList<T> a, int first = 0, int count = -1)
        {
            // for i from n - 1 downto 1 do
            // j = random integer with 0 <= j <= i
            // swap a[j] and a[i]

            int last = (count == -1) ? a.Count - 1 : first + count - 1;
            if (last <= first)
                throw new ArgumentOutOfRangeException();

            for (int i = last; i > first; --i)
            {
                int j = first + r.Next(i + 1 - first);
                if (i != j)
                {
                    T temp = a[i];
                    a[i] = a[j];
                    a[j] = temp;
                }
            }
        }
    }
}
