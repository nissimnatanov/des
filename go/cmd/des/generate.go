package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/solver"
)

type GenerateFlags struct {
	MinLevel solver.Level
	MaxLevel solver.Level
	NextArgs []string
	Count    int

	Timeout time.Duration
}

func ParseGenerateFlags(args []string) *GenerateFlags {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	f := &GenerateFlags{}
	var min, max string
	var timeout string

	fs.IntVar(&f.Count, "count", 1, "number of boards to generate")
	fs.StringVar(&min, "min", "Easy", "minimum level of the generated board (Easy, Medium, Hard, VeryHard, Evil, DarkEvil, Nightmare, BlackHole)")
	fs.StringVar(&max, "max", "", "maximum level of the generated board (optional), defaults to the min level")
	fs.StringVar(&timeout, "timeout", "", "timeout to wait for the generation process, e.g., 5s, 1m, 2h. If not set, it will wait indefinitely.")
	fs.Parse(args)
	f.MinLevel = solver.LevelFromString(min)
	if f.MinLevel == solver.LevelUnknown {
		fmt.Fprintf(os.Stderr, "Unknown min level: %s\n", min)
		os.Exit(2)
	}
	if max == "" {
		f.MaxLevel = f.MinLevel
	} else {
		f.MaxLevel = solver.LevelFromString(max)
		if f.MaxLevel == solver.LevelUnknown {
			fmt.Fprintf(os.Stderr, "Unknown max level: %s\n", min)
			os.Exit(2)
		}
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected flags, starting with %s\n", fs.Args()[0])
		os.Exit(2)
	}
	if timeout == "" || timeout == "0" {
		f.Timeout = 0
	} else {
		t, err := time.ParseDuration(timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unexpected timeout value %s: %v\n", timeout, err)
		}
		f.Timeout = t
	}
	return f
}

func generate(f *GenerateFlags) {
	if f.MinLevel > f.MaxLevel {
		fmt.Fprintf(os.Stderr, "Min level %s is greater than max level %s\n", f.MinLevel, f.MaxLevel)
		os.Exit(2)
	}
	if f.Count <= 0 {
		f.Count = 0
	}
	g := generators.New(&generators.Options{
		MinLevel: f.MinLevel,
		MaxLevel: f.MaxLevel,
		Count:    f.Count,
		OnNewResult: func(res *solver.Result) {
			fmt.Printf("Generated %s board: %s\n%s\n", res.Level, boards.Serialize(res.Input), res.Input.String())
		},
	})
	fmt.Printf("Generating boards with levels from %s to %s...\n", f.MinLevel, f.MaxLevel)
	ctx := context.Background()
	if f.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, f.Timeout)
		defer cancel()
	}

	start := time.Now()
	res := g.Generate(ctx)
	switch len(res) {
	case 0:
		fmt.Fprintln(os.Stderr, "Failed to generate any result")
		os.Exit(1)
	case 1:
		fmt.Print("Best generated board:")
	default:
		fmt.Printf("Best generated %d boards:\n", len(res))
	}

	for _, r := range res {
		fmt.Println(" ", boards.Serialize(r.Input))
	}
	fmt.Printf("Generation completed successfully in %s.\n", time.Since(start).Round(time.Millisecond))
}
