#include <iostream>
#include "SudokuBoardTest.h"
#include "asserts.h"

using namespace std;

void emptyBoard()
{
    SudokuBoardShared source = newEditBoard();
    assertThat(source->getAccessMode()).isEqualTo(SudokuBoard::AccessMode::EDIT);

    SudokuBoardConstShared empty = cloneAsImmutable(source);

    assertThat(empty->isValid()).isTrue();
    assertThat(empty->isSolved()).isFalse();

    assertThat(empty->getValue(40)).isEqualTo(SudokuValue::EMPTY);

    assertThat(empty->getRowValues(0)).isEqualTo(SudokuValueSet::none());
    assertThat(empty->getColumnValues(4)).isEqualTo(SudokuValueSet::none());
    assertThat(empty->getSquareValues(8)).isEqualTo(SudokuValueSet::none());

    assertThat(empty->getAccessMode()).isEqualTo(SudokuBoard::AccessMode::IMMUTABLE);
    assertThat(empty->getFreeCellCount()).isEqualTo(81);
    assertThat(empty->getValueCount(SudokuValue::FIVE)).isEqualTo(0);

    cerr << "--- empty board start ---" << endl;
    empty->print(cerr, "vrcst");
    cerr << "---  empty board end  ---" << endl;
}

void clone()
{
    SudokuBoardShared b1 = newEditBoard();
    SudokuBoardConstShared b2 = cloneAsImmutable(b1);
    SudokuBoardConstShared b3 = cloneAsImmutable(b2);
    SudokuPlayerBoardShared b4 = cloneToPlay(b2);
    SudokuEditBoardShared b5 = cloneToEdit(b2);

    assertThat(b1.use_count()).isEqualTo(1);
    assertThat(b2.use_count()).isEqualTo(2);
    assertThat(b2.get()).isEqualTo(b3.get());
    assertThat(b4.use_count()).isEqualTo(1);
    assertThat(b5.use_count()).isEqualTo(1);
}

void valid()
{
    SudokuEditBoardShared b1 = newEditBoard();

    int col0 = 0;
    int col1 = 3;
    int col2 = 6;
    for (SudokuValue v : SudokuValueSet::all())
    {
        b1->setValue(col0++, v);
        b1->setValue(9 + col1++, v);
        b1->setValue(18 + col2++, v);

        col0 %= 9;
        col1 %= 9;
        col2 %= 9;
    }

    assertThat(b1->isValid()).isTrue();

    b1->playValue(40, SudokuValue::FIVE);
    b1->print(cerr, "vrcst");
    assertThat(b1->isValid()).isFalse();
    assertThat(b1->isValidCell(40)).isFalse();
    assertThat(b1->isValidCell(4)).isTrue(); // single read-only

    b1->playValue(40, SudokuValue::EMPTY);
    b1->print(cerr, "vrcst");
    assertThat(b1->isValid()).isTrue();
}

void runSudokuBoardTests()
{
    cerr << "Running SudokuBoard tests..." << endl;
    emptyBoard();
    clone();
    valid();
}
