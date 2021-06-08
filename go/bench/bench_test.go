package bench

import (
	"fmt"
	"os"
	"testing"

	"github.com/nissimnatanov/des/go/board"
)

func TestMain(m *testing.M) {
	board.SetIntegrityChecks(false)
	fmt.Println("Runing Benchmark Tests")
	os.Exit(m.Run())
}
