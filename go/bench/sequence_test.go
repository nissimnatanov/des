package bench

import (
	"testing"

	"github.com/nissimnatanov/des/go/board"
)

func BenchmarkIndexFromSquare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for square := 0; square < board.SequenceSize; square++ {
			for cell := 0; cell < board.SequenceSize; cell++ {
				if board.IndexFromSquare(square, cell) < 0 {
					b.FailNow()
				}
			}
		}
	}
}
