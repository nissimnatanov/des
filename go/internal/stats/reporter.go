package stats

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type Reporter struct {
	SkipOnSilence int
	Duration      time.Duration
	OutputFile    string            // if empty, will use os.Stdout
	LogExtra      func(w io.Writer) // optional function to log extra stats at the same time
	start         time.Time
	cancel        context.CancelFunc
	outFile       *os.File
	done          chan struct{} // used to signal when the reporter is done
	extraOut      chan string
}

func (r *Reporter) LogNow(s string) {
	r.extraOut <- s
}

func (r *Reporter) Run() error {
	if r.extraOut != nil {
		panic("reporter already running")
	}
	r.start = time.Now()
	r.done = make(chan struct{})
	r.extraOut = make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	if r.OutputFile != "" {
		f, err := os.OpenFile(r.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open output file %s: %w", r.OutputFile, err)
		}
		r.outFile = f
	}

	go func() {
		ticker := time.NewTicker(r.Duration)
		defer ticker.Stop()
		skipped := 0
		prevGenerated := int64(0)
		for ctx.Err() == nil {
			select {
			case <-ctx.Done():
			case s := <-r.extraOut:
				fmt.Fprintln(r.writer(), s)
			case <-ticker.C:
				nowGenerated := Stats.Game().Count
				if nowGenerated == prevGenerated {
					skipped++
					if skipped < r.SkipOnSilence {
						continue
					}
					skipped = 0
				}
				prevGenerated = nowGenerated
				r.report()
			}
		}
		r.report()
		close(r.done)
	}()
	return nil
}
func (r *Reporter) writer() io.Writer {
	if r.outFile != nil {
		return r.outFile
	}
	return os.Stdout
}

func (r *Reporter) report() {
	w := r.writer()
	ReportStats(w)
	if r.LogExtra != nil {
		r.LogExtra(w)
	}
	fmt.Fprintf(w,
		"^^^ elapsed: %s, report time: %s ^^^\n",
		time.Since(r.start).Round(time.Second),
		time.Now().Format(time.RFC3339))
	if r.outFile != nil {
		r.outFile.Sync()
	}
}

func (r *Reporter) Stop() {
	if r.done == nil {
		panic("reporter not running")
	}
	r.cancel()
	<-r.done
	if r.outFile != nil {
		r.outFile.Sync()
		r.outFile.Close()
	}
	close(r.extraOut)
}

func ReportStats(w io.Writer) {
	gameStats := Stats.Game().String()
	solStats := Stats.Solution().String()
	cacheStats := Stats.Cache().String()
	fmt.Fprintln(w, gameStats)
	fmt.Fprintln(w, solStats)
	fmt.Fprintln(w, cacheStats)
}
