#include "SudokuValue.h"

using namespace std;

class SudokuValueSetInfo
{
  private:
    vector<SudokuValue> _values;
    int _combined;

  public:
    SudokuValueSetInfo(int mask)
    {
        _combined = 0;
        for (int v = 1; v <= 9; v++)
        {
            int vmask = 1 << (v - 1);
            if (mask & vmask)
            {
                _values.push_back(static_cast<SudokuValue>(v));
                _combined = _combined * 10 + v;
            }
        }
    }

    int size()
    {
        return _values.size();
    }

    const vector<SudokuValue> &values()
    {
        return _values;
    }

    int combined()
    {
        return _combined;
    }
};

vector<SudokuValueSetInfo> initialize()
{
    vector<SudokuValueSetInfo> all;
    for (int mask = 0; mask <= 0x1FF; mask++)
    {
        all.emplace_back(SudokuValueSetInfo(mask));
    }
    return move(all);
}

static vector<SudokuValueSetInfo> allValues = initialize();

int SudokuValueSet::size() const noexcept
{
    return allValues.at(_mask).size();
}

const vector<SudokuValue> &SudokuValueSet::values() const noexcept
{
    return allValues.at(_mask).values();
}

int SudokuValueSet::combined() const noexcept
{
    return allValues.at(_mask).combined();
}
