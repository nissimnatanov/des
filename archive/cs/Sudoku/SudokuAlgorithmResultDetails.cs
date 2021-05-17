using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Sudoku
{
    public struct SudokuAlgorithmResultDetails
    {
        private SudokuAlgorithmResultDetails(SudokuAlgorithmResult result, int indexOrLocation, byte val, int weight, string algorithmOrFailureDetails)
        {
            this.Result = result;
            m_indexOrLocation = indexOrLocation;
            m_value = val;
            m_weight = weight;
            m_algorithmOrFailureDetails = algorithmOrFailureDetails;
        }

        public static SudokuAlgorithmResultDetails Succeeded(int index, byte val, int weight)
        {
            return new SudokuAlgorithmResultDetails(SudokuAlgorithmResult.Succeeded, index, val, weight, null);
        }

        public static SudokuAlgorithmResultDetails Failed(string details, int location = -1)
        {
            return new SudokuAlgorithmResultDetails(SudokuAlgorithmResult.Failed, location, 0, 0, details);
        }

        public static SudokuAlgorithmResultDetails TwoSolutions()
        {
            return new SudokuAlgorithmResultDetails(SudokuAlgorithmResult.TwoSolutions, -1, 0, 0, "Two or more solutions were found!");
        }

        static readonly SudokuAlgorithmResultDetails s_unknown = new SudokuAlgorithmResultDetails();
        public static SudokuAlgorithmResultDetails Unknown()
        {
            return s_unknown;
        }

        private int m_indexOrLocation;
        private byte m_value;
        private int m_weight;
        private string m_algorithmOrFailureDetails;

        public SudokuAlgorithmResult Result { get; private set; }

        public int ValueIndex
        {
            get
            {
                if (this.Result != SudokuAlgorithmResult.Succeeded)
                    return -1;
                return m_indexOrLocation;
            }
        }

        public byte Value
        {
            get
            {
                if (this.Result != SudokuAlgorithmResult.Succeeded)
                    return 0;
                return m_value;
            }
        }

        public int Weight
        {
            get
            {
                return m_weight;
            }
        }

        public string Algorithm
        {
            get
            {
                if (this.Result != SudokuAlgorithmResult.Succeeded)
                    return null;
                return m_algorithmOrFailureDetails;
            }
        }

        public string FailureDetails
        {
            get
            {
                switch(this.Result)
                {
                    case SudokuAlgorithmResult.Succeeded:
                        return null;
                    case SudokuAlgorithmResult.TwoSolutions:
                        if(m_algorithmOrFailureDetails == null)
                        {
                            return "More than one solution is available!";
                        }

                        goto case SudokuAlgorithmResult.Failed;
                    case SudokuAlgorithmResult.Failed:
                        if (m_algorithmOrFailureDetails == null)
                        {
                            throw new InvalidOperationException("missing clear error message");
                        }
                        if (m_indexOrLocation >= 0)
                        {
                            string formatted = string.Format(m_algorithmOrFailureDetails, m_indexOrLocation);
                            m_algorithmOrFailureDetails = formatted;
                            m_indexOrLocation = -1;
                        }

                        return m_algorithmOrFailureDetails;
                    case SudokuAlgorithmResult.Unknown:
                        return "Could not find any viable solution";
                    default:
                        throw new InvalidOperationException();
                }
            }
        }
    }
}
