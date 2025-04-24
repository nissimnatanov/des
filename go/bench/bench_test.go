package bench

import (
	"fmt"
	"os"
	"testing"

	"github.com/nissimnatanov/des/go/boards"
)

func TestMain(m *testing.M) {
	boards.SetIntegrityChecks(false)
	fmt.Println("Running Benchmark Tests")
	os.Exit(m.Run())
}
