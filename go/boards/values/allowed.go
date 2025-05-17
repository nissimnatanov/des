package values

import (
	"github.com/nissimnatanov/des/go/boards/indexes"
)

type Allowed struct {
	byRelated  [indexes.BoardSize]Set
	byUser     [indexes.BoardSize]Set
	emptyCells indexes.BitSet81
	hints01    indexes.BitSet81
}

func (a *Allowed) Get(index int) Set {
	return Union(a.byRelated[index], a.byUser[index]).Complement()
}

func (a *Allowed) GetByRelated(index int) Set {
	return a.byRelated[index].Complement()
}

// ReportPresent is used when board cell has a value set
func (a *Allowed) ReportPresent(index int) {
	a.byRelated[index] = FullSet()
	a.byUser[index] = EmptySet()
	a.emptyCells.Set(index, false)
	a.hints01.Set(index, false)
}

func (a *Allowed) ReportEmpty(index int, related Set) {
	a.byRelated[index] = related
	a.emptyCells.Set(index, true)
	a.updateHint(index)
}

func (a *Allowed) DisallowRelated(index int, v Value) {
	if !a.emptyCells.Get(index) {
		panic("disallowing a value in a cell that has a value")
	}
	a.byRelated[index] = a.byRelated[index].With(v.AsSet())
	a.updateHint(index)
}

func (a *Allowed) Hint01() (int, bool) {
	return a.hints01.First()
}

func (a *Allowed) AllowAll() {
	clear(a.byRelated[:])
	clear(a.byUser[:])
	a.emptyCells.SetAll(true)
	a.hints01.SetAll(false)
}

func (a *Allowed) Clone() Allowed {
	// all fields are fixed size arrays by value
	return *a
}

func (a *Allowed) GetDisallowedByUser(index int) Set {
	return a.byUser[index]
}

func (a *Allowed) DisallowByUser(index int, vs Set) {
	a.byUser[index] = a.byUser[index].With(vs)
	a.updateHint(index)
}

func (a *Allowed) ResetDisallowedByUser(index int) {
	a.byUser[index] = EmptySet()
	a.updateHint(index)
}

func (a *Allowed) updateHint(index int) {
	isHint01 := a.Get(index).Size() <= 1
	a.hints01.Set(index, isHint01)
}

// IndexesByAllowedSize returns a BitSet81 of indexes that have the same size as the given size,
// use for troubleshooting only
func (a *Allowed) IndexesByAllowedSize(size int) indexes.BitSet81 {
	// calculate manually
	emptyCells := a.emptyCells
	var allowedBySize indexes.BitSet81
	for i := range emptyCells.Indexes {
		if a.Get(i).Size() == size {
			allowedBySize.Set(i, true)
		}
	}
	return allowedBySize
}
