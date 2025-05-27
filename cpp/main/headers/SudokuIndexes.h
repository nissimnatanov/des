#pragma once

#include <array>

using sequence = std::array<int, 9>;
using related = std::array<int, 20>;

inline int rowFromIndex(int index) noexcept
{
    return index / 9;
}

inline int columnFromIndex(int index) noexcept
{
    return index % 9;
}

inline int indexFromCoordinates(int row, int col) noexcept
{
    return row * 9 + col;
}

inline int indexFromSquare(int squareIndex, int cellIndex) noexcept
{
    // index of the first cell
    int index = (squareIndex / 3) * 27 + (squareIndex % 3) * 3;
    // add cell location relative to first one
    index += (cellIndex / 3) * 9 + (cellIndex % 3);
    return index;
}

inline int squareFromIndex(int index) noexcept
{
    int square = index / 3;
    square = (square / 9) * 3 + square % 3;
    return square;
}

inline int cellIndexInSquareFromIndex(int index) noexcept
{
    // rows (3,4,5) and (6,7,8) are equivalent to (0,1,2)
    int row = (index / 9) % 3;
    return index % 3 + row * 3;
}

const sequence &getRowIndexes(int row);

const sequence &getColumnIndexes(int column);

const sequence &getSquareIndexes(int square);

const std::array<int, 81> &getBoardIndexes();

const sequence &getSequenceIndexes();

const related &getRelated(int index);
