package values

import (
	"github.com/nissimnatanov/des/go/boards/indexes"
)

type Allowed struct {
	disallowedByRelated [indexes.BoardSize]Set
	disallowedByUser    [indexes.BoardSize]Set
	emptyCells          indexes.BitSet81
	hints01             indexes.BitSet81
}

func (a *Allowed) Get(index int) Set {
	return Union(a.disallowedByRelated[index], a.disallowedByUser[index]).Complement()
}

func (a *Allowed) GetByRelated(index int) Set {
	return a.disallowedByRelated[index].Complement()
}

// ReportPresent is used when board cell has a value set
func (a *Allowed) ReportPresent(index int) {
	a.disallowedByRelated[index] = FullSet
	a.disallowedByUser[index] = EmptySet
	a.emptyCells.Set(index, false)
	a.hints01.Set(index, false)
}

func (a *Allowed) ReportEmpty(index int, related Set) {
	a.disallowedByRelated[index] = related
	a.emptyCells.Set(index, true)
	a.updateHint(index)
}

func (a *Allowed) DisallowRelated(index int, v Value) {
	if !a.emptyCells.Get(index) {
		panic("disallowing a value in a cell that has a value")
	}
	a.disallowedByRelated[index] = a.disallowedByRelated[index].With(v.AsSet())
	a.updateHint(index)
}

func (a *Allowed) Hint01() int {
	return a.hints01.First()
}

func (a *Allowed) AllowAll() {
	clear(a.disallowedByRelated[:])
	clear(a.disallowedByUser[:])
	a.emptyCells = indexes.MaxBitSet81
	a.hints01 = indexes.MinBitSet81
}

func (a *Allowed) Clone() Allowed {
	// all fields are fixed size arrays by value
	return *a
}

func (a *Allowed) GetDisallowedByUser(index int) Set {
	return a.disallowedByUser[index]
}

func (a *Allowed) DisallowByUser(index int, vs Set) {
	a.disallowedByUser[index] = a.disallowedByUser[index].With(vs)
	a.updateHint(index)
}

func (a *Allowed) ResetDisallowedByUser(index int) {
	a.disallowedByUser[index] = EmptySet
	a.updateHint(index)
}

func (a *Allowed) updateHint(index int) {
	isHint01 := a.Get(index).Size() <= 1
	a.hints01.Set(index, isHint01)
}
