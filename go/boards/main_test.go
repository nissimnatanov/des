package boards

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetIntegrityChecks(true)
	os.Exit(m.Run())
}
