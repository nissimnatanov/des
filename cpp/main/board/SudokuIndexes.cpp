#include "SudokuIndexes.h"
#include <memory>
#include <string>
#include <stdexcept>

using namespace std;

array<sequence, 9> initRowIndexes()
{
    array<sequence, 9> rows;
    for (int row = 0; row < 9; row++)
    {
        sequence indexes;
        for (int column = 0; column < 9; column++)
        {
            indexes.at(column) = indexFromCoordinates(row, column);
        }
        rows.at(row) = std::move(indexes);
    }
    return std::move(rows);
}

array<sequence, 9> initColumnIndexes()
{
    array<sequence, 9> columns;
    for (int column = 0; column < 9; column++)
    {
        sequence indexes;
        for (int row = 0; row < 9; row++)
        {
            indexes.at(row) = indexFromCoordinates(row, column);
        }
        columns.at(column) = std::move(indexes);
    }
    return std::move(columns);
}

array<sequence, 9> initSquareIndexes()
{
    array<sequence, 9> squares;
    for (int square = 0; square < 9; square++)
    {
        sequence indexes;
        for (int cell = 0; cell < 9; cell++)
        {
            indexes.at(cell) = indexFromSquare(square, cell);
        }
        squares.at(square) = std::move(indexes);
    }
    return std::move(squares);
}

array<int, 81> initBoardIndexes()
{
    array<int, 81> indexes;
    for (int i = 0; i < 81; i++)
    {
        indexes.at(i) = i;
    }
    return std::move(indexes);
}

array<related, 81> initIndexToRelated()
{
    array<related, 81> all;
    for (int i = 0; i < 81; i++)
    {
        int relatedIndex = 0;
        int row = rowFromIndex(i);
        int column = columnFromIndex(i);
        int square = squareFromIndex(i);

        for (int rowIndex : getRowIndexes(row))
        {
            if (rowIndex == i)
            {
                continue;
            }
            all.at(i).at(relatedIndex++) = rowIndex;
        }
        for (int columnIndex : getColumnIndexes(column))
        {
            if (columnIndex == i)
            {
                continue;
            }
            all.at(i).at(relatedIndex++) = columnIndex;
        }
        for (int squareIndex : getSquareIndexes(square))
        {
            if (row == rowFromIndex(squareIndex) ||
                column == columnFromIndex(squareIndex))
            {
                continue;
            }
            all.at(i).at(relatedIndex++) = squareIndex;
        }
        if (relatedIndex != 20)
        {
            throw logic_error(string("expected # of related indexes to be 20, got ") + to_string(relatedIndex));
        }
    }

    return std::move(all);
}

array<sequence, 9> rowIndexes = initRowIndexes();
array<sequence, 9> columnIndexes = initColumnIndexes();
array<sequence, 9> squareIndexes = initSquareIndexes();
array<int, 81> boardIndexes = initBoardIndexes();
array<related, 81> indexToRelated = initIndexToRelated();

const sequence &getRowIndexes(int row)
{
    return rowIndexes.at(row);
}

const sequence &getColumnIndexes(int column)
{
    return columnIndexes.at(column);
}

const sequence &getSquareIndexes(int square)
{
    return squareIndexes.at(square);
}

const std::array<int, 81> &getBoardIndexes()
{
    return boardIndexes;
}

const sequence &getSequenceIndexes()
{
    return rowIndexes.at(0);
}

const related &getRelated(int index)
{
    return indexToRelated.at(index);
}
