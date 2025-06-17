package solver

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nissimnatanov/des/go/boards"
)

const nightmareSerialized = "nightmare.serialized"
const nightmareLog = "nightmare.log"

var disableNLog bool

var logFolder = func() string {
	var found string
	{
		wd, err := os.Getwd()
		if err == nil {
			for wd != "" && err == nil {
				ns := path.Join(wd, nightmareSerialized)
				st, err := os.Stat(ns)
				if err == nil && !st.IsDir() {
					found = wd
					break
				}
				parent := path.Dir(wd)
				if parent == wd {
					break // reached root directory
				}
				wd = parent
			}
		}
	}
	if found == "" {
		found = os.TempDir()
	}
	return found
}()

var registeredCache = map[string]bool{}

func logResult(res *Result) {
	if disableNLog {
		return
	}
	if res.Status != StatusSucceeded {
		return
	}
	if res.Level < LevelNightmare || res.Action != ActionSolve || res.Status != StatusSucceeded {
		return
	}
	// already have enough logged below 80K
	if res.Complexity < 80000 {
		return
	}
	serialized := boards.Serialize(res.Input)
	if registered, ok := registeredCache[serialized]; ok && registered {
		return // already registered
	}
	registered := false
	nightmareSerialized := filepath.Join(logFolder, nightmareSerialized)
	nightmareLog := filepath.Join(logFolder, nightmareLog)
	b, err := os.ReadFile(nightmareSerialized)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to read the nightmare serialized file %s: %v\n",
			nightmareSerialized, err)
		return
	}

	lines := strings.Split(string(b), "\n")
	registered = slices.Contains(lines, serialized)
	if registered {
		registeredCache[serialized] = true
		return
	}

	ns, err := os.OpenFile(nightmareSerialized, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		_, err = fmt.Fprintln(ns, serialized)
		ns.Close()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to append to the nightmare serialized file %s: %v\n",
			nightmareSerialized, err)
		return
	}
	registeredCache[serialized] = true

	nl, err := os.OpenFile(nightmareLog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		_, err = fmt.Fprintf(nl,
			"Solver reached level %s with complexity %d: %s\n%s\n",
			res.Level,
			res.Complexity,
			boards.Serialize(res.Input),
			res.Input.String())
		nl.Close()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to append to the nightmare log file %s: %v\n",
			nightmareLog, err)
	}
}
