package bench

import (
	"fmt"
	"os"
	"testing"

	"github.com/nissimnatanov/des/go/board"
)

func TestMain(m *testing.M) {
	board.SetIntegrityChecks(false)
	fmt.Println("Running Benchmark Tests")
	os.Exit(m.Run())
}
