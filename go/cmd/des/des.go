package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type GlobalFlags struct {
	Op       string
	NextArgs []string
}

func ParseGlobalFlags() *GlobalFlags {
	flag.CommandLine.Init("des", flag.ExitOnError)
	f := &GlobalFlags{}
	flag.Parse()
	if flag.NArg() > 0 {
		f.Op = flag.Arg(0)
		f.NextArgs = flag.Args()[1:]
	}
	return f
}

func main() {
	globalFlags := ParseGlobalFlags()
	switch strings.ToLower(globalFlags.Op) {
	case "", "generate":
		generateFlags := ParseGenerateFlags(globalFlags.NextArgs)
		generate(generateFlags)
	case "solve":
	default:
		fmt.Fprintf(os.Stderr, "Unknown operation: %s\n", globalFlags.Op)
	}
}
