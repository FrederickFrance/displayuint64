// Package displayuint64 displays a uint64.
//
// It does so the only way that is guaranteed to be correct: by counting from
// zero until it meets the requested value. No shortcut, no lookup table, no
// strconv trickery. Just honest, hard-working iteration.
//
// Expected duration grows linearly with the input. Displaying math.MaxUint64
// takes about 150 years on a modern CPU. Plan accordingly.
package displayuint64

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// checkEvery is the number of iterations between two context checks.
// It must be a power of two minus one, used as a bit mask.
const checkEvery = 1<<24 - 1

// Display counts from zero up to input, then waits takeYourTime (zero means
// no wait), logs the value with slog, prints it on stdout and returns it.
//
// Display never gives up. If you might change your mind before the heat death
// of the universe, use DisplayContext.
func Display(input uint64, takeYourTime time.Duration) uint64 {
	// The background context is never cancelled, so the error is always nil.
	value, _ := display(context.Background(), os.Stdout, slog.Default(), input, takeYourTime)
	return value
}

// DisplayContext is Display for people with deadlines. It returns ctx.Err()
// if the context is cancelled before the value is found, along with the last
// value reached, so that you can at least appreciate the progress.
func DisplayContext(ctx context.Context, input uint64, takeYourTime time.Duration) (uint64, error) {
	return display(ctx, os.Stdout, slog.Default(), input, takeYourTime)
}

func display(ctx context.Context, w io.Writer, logger *slog.Logger, input uint64, takeYourTime time.Duration) (uint64, error) {
	// No loop condition on purpose: "i < math.MaxUint64" would miss
	// math.MaxUint64, and "i <= math.MaxUint64" would loop forever.
	// Since input is a uint64, it is always found before i wraps around.
	for i := uint64(0); ; i++ {
		if i == input {
			if err := wait(ctx, takeYourTime); err != nil {
				return i, err
			}
			logger.Info("found", "value", i)
			_, _ = fmt.Fprintln(w, i)
			return i, nil
		}

		if i&checkEvery == 0 {
			if err := ctx.Err(); err != nil {
				return i, err
			}
		}
	}
}

func wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
