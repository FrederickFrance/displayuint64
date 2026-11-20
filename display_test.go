package displayuint64

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, nil))
}

func TestDisplayFindsValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input uint64
	}{
		{name: "zero", input: 0},
		{name: "one", input: 1},
		{name: "the answer", input: 42},
		{name: "past the first context check", input: checkEvery + 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var out, logs bytes.Buffer
			got, err := display(context.Background(), &out, newLogger(&logs), tt.input, 0)

			require.NoError(t, err)
			require.Equal(t, tt.input, got)
			require.Equal(t, formatLine(tt.input), out.String())
			require.Contains(t, logs.String(), "msg=found")
			require.Contains(t, logs.String(), "value="+formatValue(tt.input))
		})
	}
}

func TestDisplayTakesItsTime(t *testing.T) {
	t.Parallel()

	const delay = 50 * time.Millisecond
	var out, logs bytes.Buffer

	start := time.Now()
	got, err := display(context.Background(), &out, newLogger(&logs), 3, delay)

	require.NoError(t, err)
	require.Equal(t, uint64(3), got)
	require.GreaterOrEqual(t, time.Since(start), delay)
	require.Equal(t, "3\n", out.String())
}

func TestDisplayNegativeDurationDoesNotWait(t *testing.T) {
	t.Parallel()

	var out, logs bytes.Buffer
	got, err := display(context.Background(), &out, newLogger(&logs), 7, -time.Hour)

	require.NoError(t, err)
	require.Equal(t, uint64(7), got)
}

// math.MaxUint64 would take about 150 years, so we only check that the
// search can be interrupted. Patience is not unit-testable.
func TestDisplayMaxUint64CanBeCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var out, logs bytes.Buffer
	reached, err := display(ctx, &out, newLogger(&logs), math.MaxUint64, 0)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, reached, uint64(math.MaxUint64))
	require.Zero(t, reached&checkEvery, "cancellation is only checked on mask boundaries")
	require.Empty(t, out.String())
	require.Empty(t, logs.String())
}

func TestDisplayAlreadyCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var out, logs bytes.Buffer
	reached, err := display(ctx, &out, newLogger(&logs), 1_000, 0)

	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, uint64(0), reached)
	require.Empty(t, out.String())
}

// Input 0 is found before any context check, so only the wait can notice
// the cancellation.
func TestDisplayCancelledDuringWait(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	var out, logs bytes.Buffer
	reached, err := display(ctx, &out, newLogger(&logs), 0, time.Hour)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, uint64(0), reached)
	require.Empty(t, out.String())
}

func TestDisplayContextPublicAPI(t *testing.T) {
	t.Parallel()

	got, err := DisplayContext(context.Background(), 5, 0)

	require.NoError(t, err)
	require.Equal(t, uint64(5), got)
}

func TestDisplayPublicAPI(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint64(12), Display(12, 0))
}

func BenchmarkDisplay(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, input := range []uint64{1_000, 1_000_000, 1_000_000_000} {
		b.Run(formatValue(input), func(b *testing.B) {
			for b.Loop() {
				_, _ = display(context.Background(), io.Discard, logger, input, 0)
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(input), "ns/iteration")
		})
	}
}

func formatValue(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func formatLine(v uint64) string {
	return formatValue(v) + "\n"
}
