#pragma once

#include <memory>
#include <array>
#include "Random.h"
#include "SudokuBoard.h"

using namespace std;

class SudokuBoardIndexManager
{
private:
  // indexes are reserved/removed from the back!
  array<int, SudokuBoard::BOARD_SIZE> _indexes;
  int _remained;
  int _reserved;

public:
  SudokuBoardIndexManager &operator=(const SudokuBoardIndexManager &) = delete;

  SudokuBoardIndexManager(const SudokuBoardIndexManager &) = default;
  SudokuBoardIndexManager()
  {
    for (int i = 0; i < SudokuBoard::BOARD_SIZE; i++)
    {
      _indexes[i] = i;
    }
    _remained = SudokuBoard::BOARD_SIZE;
    _reserved = 0;
  }

  void shuffleRemained(Random &r)
  {
    r.shuffle(_indexes, 0, _remained);
  }

  void shuffleRemoved(Random &r)
  {
    int start = _remained + _reserved;
    int count = _indexes.size() - start;
    r.shuffle(_indexes, start, count);
  }

  int remainedSize() const
  {
    return _remained;
  }

  void swapRemained(int ri1, int ri2)
  {
    if (ri1 < 0 || ri1 >= _remained || ri2 < 0 || ri2 >= _remained)
    {
      throw logic_error("swapRemained received bad index");
    }
    if (ri1 != ri2)
    {
      swap(_indexes[ri1], _indexes[ri2]);
    }
  }

  void reserve(int n)
  {
    if (n <= 0 || n > _remained)
    {
      throw out_of_range("n is out of range");
    }
    if (_reserved != 0)
    {
      throw logic_error("cannot reserve twice");
    }

    _reserved = n;
    _remained -= n;
  }

  void removeReserved()
  {
    _reserved = 0;
  }

  void revertReserved()
  {
    _remained += _reserved;
    _reserved = 0;
  }

  class Iterable
  {
  private:
    friend class SudokuBoardIndexManager;
    const SudokuBoardIndexManager &_source;
    const int _start;
    const int _count;

    Iterable(const SudokuBoardIndexManager &source, int start, int count)
        : _source(source), _start(start), _count(count) {}

  public:
    auto begin() const
    {
      return _source._indexes.cbegin() + _start;
    }
    auto end() const
    {
      return _source._indexes.cbegin() + _start + _count;
    }
  };

  Iterable reserved() const
  {
    return Iterable(*this, _remained, _reserved);
  }

  Iterable remained() const
  {
    return Iterable(*this, 0, _remained);
  }

  Iterable removed() const
  {
    int start = _remained + _reserved;
    int count = _indexes.size() - start;
    return Iterable(*this, start, count);
  }

  bool tryPrioritizeIndex(int index)
  {
    for (int i = 0; i < _remained; i++)
    {
      if (_indexes.at(i) == index)
      {
        swapRemained(i, _remained - 1);
        return true;
      }
    }

    return false;
  }

  void prioritizeIndex(int index)
  {
    if (!tryPrioritizeIndex(index))
    {
      throw logic_error("Index to prioritize was not found");
    }
  }

  void restoreRemoved(int boardIndex)
  {
    for (int i = _remained + _reserved; i < _indexes.size(); i++)
    {
      if (_indexes.at(i) == boardIndex)
      {
        if (_remained != i)
        {
          swap(_indexes.at(_remained), _indexes.at(i));
        }
        if (_reserved > 0 && i != (_remained + _reserved))
        {
          swap(_indexes.at(_remained + _reserved), _indexes.at(i));
        }
        _remained++;
        return;
      }
    }

    throw logic_error("Index to restore from removed back to remained was not found");
  }

  void removeIndex(int index)
  {
    prioritizeIndex(index);
    reserve(1);
    removeReserved();
  }

  bool tryRemoveIndex(int index)
  {
    if (!tryPrioritizeIndex(index))
    {
      return false;
    }
    reserve(1);
    removeReserved();
    return true;
  }
};
