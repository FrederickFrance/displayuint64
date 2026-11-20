package main

import (
	"bytes"
	"context"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSuccess(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(context.Background(), []string{"42"}, &stderr)

	require.Equal(t, 0, code)
	require.Empty(t, stderr.String())
}

func TestRunWithTakeYourTime(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(context.Background(), []string{"-take-your-time", "10ms", "3"}, &stderr)

	require.Equal(t, 0, code)
}

func TestRunUsageErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{name: "no argument", args: nil, message: "usage:"},
		{name: "too many arguments", args: []string{"1", "2"}, message: "usage:"},
		{name: "negative", args: []string{"-1"}, message: "flag provided but not defined"},
		{name: "not a number", args: []string{"forty-two"}, message: `not a uint64: "forty-two"`},
		{name: "overflow", args: []string{"18446744073709551616"}, message: "not a uint64"},
		{name: "unknown flag", args: []string{"-fast", "1"}, message: "flag provided but not defined"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			code := run(context.Background(), tt.args, &stderr)

			require.Equal(t, 2, code)
			require.Contains(t, stderr.String(), tt.message)
		})
	}
}

func TestRunTimeout(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	maxValue := strconv.FormatUint(math.MaxUint64, 10)
	code := run(context.Background(), []string{"-timeout", "50ms", maxValue}, &stderr)

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "timed out after reaching")
}

func TestRunInterrupted(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stderr bytes.Buffer
	code := run(ctx, []string{"1000"}, &stderr)

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "interrupted after reaching 0")
}

func TestProgress(t *testing.T) {
	t.Parallel()

	require.InDelta(t, 100.0, progress(0, 0), 0)
	require.InDelta(t, 50.0, progress(5, 10), 1e-9)
	require.InDelta(t, 0.0, progress(0, math.MaxUint64), 0)
}
