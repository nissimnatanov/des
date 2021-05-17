#include <iostream>
#include "SudokuIndexesTest.h"
#include "asserts.h"

using namespace std;

void rowFromIndexTest()
{
    assertThat(rowFromIndex(45)).isEqualTo(5);
}

void columnFromIndexTest()
{
    assertThat(columnFromIndex(46)).isEqualTo(1);
}

void indexFromCoordinatesTest()
{
    assertThat(indexFromCoordinates(5, 2)).isEqualTo(47);
}

void indexFromSquareTest()
{
    assertThat(indexFromSquare(0, 0)).isEqualTo(0);
    assertThat(indexFromSquare(5, 2)).isEqualTo(35);
    assertThat(indexFromSquare(8, 8)).isEqualTo(80);
}

void getRowIndexesTest()
{
    assertThat(getRowIndexes(7).size()).isEqualTo(9);
    assertThat(getRowIndexes(3)[7]).isEqualTo(34);
}

void getColumnIndexesTest()
{
    assertThat(getColumnIndexes(8).size()).isEqualTo(9);
    assertThat(getColumnIndexes(8)[8]).isEqualTo(80);
}

void getSquareIndexesTest()
{
    assertThat(getSquareIndexes(5).size()).isEqualTo(9);
    assertThat(getSquareIndexes(0)[0]).isEqualTo(0);
    assertThat(getSquareIndexes(0)[8]).isEqualTo(20);
    assertThat(getSquareIndexes(1)[6]).isEqualTo(21);
    assertThat(getSquareIndexes(2)[5]).isEqualTo(17);
    assertThat(getSquareIndexes(3)[4]).isEqualTo(37);
    assertThat(getSquareIndexes(4)[3]).isEqualTo(39);
    assertThat(getSquareIndexes(5)[2]).isEqualTo(35);
    assertThat(getSquareIndexes(6)[0]).isEqualTo(54);
    assertThat(getSquareIndexes(7)[6]).isEqualTo(75);
    assertThat(getSquareIndexes(8)[8]).isEqualTo(80);
}

void getBoardIndexesTest()
{
    assertThat(getBoardIndexes().size()).isEqualTo(81);
    assertThat(getBoardIndexes()[12]).isEqualTo(12);
}

void getSequenceIndexesTest()
{
    assertThat(getSequenceIndexes().size()).isEqualTo(9);
    assertThat(getSequenceIndexes()[7]).isEqualTo(7);
}

void squareFromIndexTest()
{
    assertThat(squareFromIndex(2)).isEqualTo(0);
    assertThat(squareFromIndex(32)).isEqualTo(4);
    assertThat(squareFromIndex(80)).isEqualTo(8);
}

void cellIndexInSquareFromIndexTest()
{
    assertThat(cellIndexInSquareFromIndex(2)).isEqualTo(2);
    assertThat(cellIndexInSquareFromIndex(41)).isEqualTo(5);
    assertThat(cellIndexInSquareFromIndex(80)).isEqualTo(8);
}

void getRelatedTest()
{
    assertThat(getRelated(0)[0]).isEqualTo(1);
    assertThat(getRelated(0)[7]).isEqualTo(8);
    assertThat(getRelated(0)[8]).isEqualTo(9);
    assertThat(getRelated(0)[15]).isEqualTo(72);
    assertThat(getRelated(0)[19]).isEqualTo(20);

    assertThat(getRelated(40)[0]).isEqualTo(36);
    assertThat(getRelated(40)[7]).isEqualTo(44);
    assertThat(getRelated(40)[8]).isEqualTo(4);
    assertThat(getRelated(40)[15]).isEqualTo(76);
    assertThat(getRelated(40)[19]).isEqualTo(50);

    assertThat(getRelated(80)[0]).isEqualTo(72);
    assertThat(getRelated(80)[7]).isEqualTo(79);
    assertThat(getRelated(80)[8]).isEqualTo(8);
    assertThat(getRelated(80)[15]).isEqualTo(71);
    assertThat(getRelated(80)[19]).isEqualTo(70);
}

void runSudokuIndexesTests()
{
    cerr << "Running SudokuIndexes tests..." << endl;
    rowFromIndexTest();
    columnFromIndexTest();
    indexFromCoordinatesTest();
    indexFromSquareTest();
    getRowIndexesTest();
    getColumnIndexesTest();
    getSquareIndexesTest();
    getBoardIndexesTest();
    getSequenceIndexesTest();
    squareFromIndexTest();
    cellIndexInSquareFromIndexTest();
    getRelatedTest();
}