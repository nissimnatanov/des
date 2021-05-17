#include <iostream>
#include "SudokuValueTest.h"
#include "asserts.h"

using namespace std;

void emptySet()
{
    SudokuValueSet set;

    assertThat(set.size()).isEqualTo(0);
}

void setWithOneValue()
{
    SudokuValueSet set(SudokuValue::THREE);

    assertThat(set.size()).isEqualTo(1);
    assertThat(set[0]).isEqualTo(SudokuValue::THREE);
}

void set_add()
{
    SudokuValueSet set1(SudokuValue::THREE, SudokuValue::SEVEN);
    SudokuValueSet set2(SudokuValue::FIVE);

    auto set3 = set1 + set2;

    assertThat(set3.size()).isEqualTo(3);
    assertThat(set3[0]).isEqualTo(SudokuValue::THREE);
    assertThat(set3[1]).isEqualTo(SudokuValue::FIVE);
    assertThat(set3[2]).isEqualTo(SudokuValue::SEVEN);
}

void set_complement()
{
    SudokuValueSet set1(SudokuValue::THREE);

    auto set2 = ~set1;

    assertThat(set2.size()).isEqualTo(8);
    assertThat(set2[0]).isEqualTo(SudokuValue::ONE);
    assertThat(set2[2]).isEqualTo(SudokuValue::FOUR);
    assertThat(set2[7]).isEqualTo(SudokuValue::NINE);
}

void set_minus()
{
    SudokuValueSet set(SudokuValue::THREE, SudokuValue::FIVE);

    set = ~set - SudokuValueSet(SudokuValue::ONE, SudokuValue::FIVE, SudokuValue::NINE);

    assertThat(set.size()).isEqualTo(5);
    assertThat(set[0]).isEqualTo(SudokuValue::TWO);
    assertThat(set[2]).isEqualTo(SudokuValue::SIX);
    assertThat(set[4]).isEqualTo(SudokuValue::EIGHT);
}

void equals()
{
    SudokuValueSet set1(SudokuValue::TWO);
    SudokuValueSet set2(SudokuValue::SEVEN);
    SudokuValueSet set3(SudokuValue::TWO, SudokuValue::FIVE);

    assertThat(set1 == set1).isEqualTo(true);
    assertThat(set1 != set1).isEqualTo(false);
    assertThat(set1 == set2).isEqualTo(false);
    assertThat(set1 != set2).isEqualTo(true);
    assertThat(set1 == set3).isEqualTo(false);
}

void contains()
{
    SudokuValueSet set1(SudokuValue::TWO, SudokuValue::FIVE, SudokuValue::NINE);
    SudokuValueSet set2(SudokuValue::SEVEN);
    SudokuValueSet set3(SudokuValue::TWO, SudokuValue::NINE);
    SudokuValueSet set4(SudokuValue::TWO, SudokuValue::EIGHT);

    assertThat(set1.contains(SudokuValue::NINE)).isEqualTo(true);
    assertThat(set1.contains(SudokuValue::ONE)).isEqualTo(false);
    assertThat(set1.containsAll(SudokuValueSet())).isEqualTo(true);
    assertThat(set1.containsAll(set2)).isEqualTo(false);
    assertThat(set1.containsAll(set3)).isEqualTo(true);
    assertThat(set1.containsAll(set4)).isEqualTo(false);
}

void runSudokuValueTests()
{
    cerr << "Running SudokuValue tests..." << endl;
    emptySet();
    setWithOneValue();
    set_add();
    set_complement();
    set_minus();
    equals();
    contains();
}
