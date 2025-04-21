package bench

import (
	"testing"

	"github.com/nissimnatanov/des/go/board"
)

func BenchmarkIndexFromSquare(b *testing.B) {
	for range b.N {
		for square := range board.SequenceSize {
			for cell := range board.SequenceSize {
				if board.IndexFromSquare(square, cell) < 0 {
					b.FailNow()
				}
			}
		}
	}
}
