# displayuint64

[![CI](https://github.com/FrederickFrance/displayuint64/actions/workflows/ci.yml/badge.svg)](https://github.com/FrederickFrance/displayuint64/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/FrederickFrance/displayuint64.svg)](https://pkg.go.dev/github.com/FrederickFrance/displayuint64)

**Possibly the most essential Go module of all time, in all known universes.**

Every day, countless programs display numbers. Most of them rely on
`strconv`, `fmt` or other conversion shortcuts, trusting decades-old
algorithms that nobody in your team has ever audited.

**displayuint64** takes a different path. It does not *convert* your number.
It *finds* it, by counting from zero until the value is reached, one
iteration at a time. No lookup table. No division. No trickery.
Only arithmetic so simple it cannot be wrong.

## Features

- ✅ **Correct by construction**: every value from `0` to `math.MaxUint64` is supported
- ✅ **100% test coverage**
- ✅ **Zero lookup tables**
- ✅ **Context-aware**, for organisations with deadlines
- ✅ **Structured logging** via `log/slog`, for full observability of each discovery
- ✅ **Dramatic pause** before the reveal, configurable to the nanosecond
- ✅ **Perfectly predictable performance**: O(n), guaranteed, every time

## Installation

```sh
go get github.com/FrederickFrance/displayuint64
```

## Usage

```go
package main

import "github.com/FrederickFrance/displayuint64"

func main() {
	displayuint64.Display(42, 0)
}
```

For presentations, product launches and other high-stakes situations, add a
dramatic pause before the value is revealed:

```go
displayuint64.Display(42, 3*time.Second)
```

Mission-critical systems with service-level objectives should use the
context-aware API. On cancellation, it returns the last value reached, so that
no progress is ever lost:

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

reached, err := displayuint64.DisplayContext(ctx, math.MaxUint64, 0)
```

### Command line

```sh
go install github.com/FrederickFrance/displayuint64/cmd/displayuint64@latest

displayuint64 42
displayuint64 -take-your-time 2s 42
displayuint64 -timeout 2s 18446744073709551615
# timed out after reaching 6492782592 (0.000000% done)
```

## Performance

displayuint64 runs at about **0.26 ns per iteration** on a 2.3 GHz Xeon
(`go test -bench .`). Unlike conventional approaches, its execution time
depends only on the value displayed, which makes capacity planning trivial.

| Value                          | Expected time |
|--------------------------------|---------------|
| 1 000                          | instant       |
| 1 000 000                      | 0.3 ms        |
| 1 000 000 000                  | 0.26 s        |
| 1 000 000 000 000              | 4 min         |
| 1 000 000 000 000 000          | 3 days        |
| 1 000 000 000 000 000 000      | 8 years       |
| 18 446 744 073 709 551 615     | ~150 years    |

Displaying `math.MaxUint64` is a long-term commitment. We recommend starting
early and documenting progress for the next generations of maintainers.

## Correctness

A naive implementation would write
`for i := uint64(0); i <= math.MaxUint64; i++`. This condition is always true,
and the loop never ends. Writing `i < math.MaxUint64` instead silently drops
support for `math.MaxUint64`.

displayuint64 suffers from neither defect. We take correctness seriously.

## Roadmap

- `displayuint128`, as soon as Go provides a `uint128` type
- A distributed version, to share the counting load across a cluster
- Full test of `math.MaxUint64`, scheduled for completion around 2176

## FAQ

**Can it go faster?**
Faster hardware is supported out of the box.

**Why not just use `fmt.Println(n)`?**
`fmt.Println` gives you the number. displayuint64 gives you the journey.

**Is this a joke?**
displayuint64 has 100% test coverage.

## License

[MIT](LICENSE)
