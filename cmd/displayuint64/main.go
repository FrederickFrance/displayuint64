// Command displayuint64 displays a uint64, the hard way.
//
//	displayuint64 [-take-your-time 2s] [-timeout 1m] <value>
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"

	"github.com/FrederickFrance/displayuint64"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	os.Exit(run(ctx, os.Args[1:], os.Stderr))
}

func run(ctx context.Context, args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("displayuint64", flag.ContinueOnError)
	fs.SetOutput(stderr)
	takeYourTime := fs.Duration("take-your-time", 0, "pause before displaying, for dramatic effect")
	timeout := fs.Duration("timeout", 0, "give up after this duration (0 = never, as intended)")
	fs.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "usage: displayuint64 [flags] <value>")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return 2
	}

	input, err := strconv.ParseUint(fs.Arg(0), 10, 64)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "not a uint64: %q\n", fs.Arg(0))
		return 2
	}

	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	reached, err := displayuint64.DisplayContext(ctx, input, *takeYourTime)
	if err != nil {
		reason := "interrupted"
		if errors.Is(err, context.DeadlineExceeded) {
			reason = "timed out"
		}
		_, _ = fmt.Fprintf(stderr, "%s after reaching %d (%.6f%% done)\n", reason, reached, progress(reached, input))
		return 1
	}
	return 0
}

func progress(reached, input uint64) float64 {
	if input == 0 {
		return 100
	}
	return float64(reached) / float64(input) * 100
}
