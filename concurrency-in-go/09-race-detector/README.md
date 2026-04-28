# Section 9: Go Race Detector

## What is the Race Detector?
- Built-in tool for finding race conditions in Go code
- Uses ThreadSanitizer (TSan) technology at runtime
- Commands: `go run -race`, `go test -race`, `go build -race`

## Characteristics
- Binary must be race-enabled
- Detects races at runtime (warning printed)
- Race-enabled binary is ~10x slower and uses ~10x more memory
- Good candidates: integration tests, load tests

## Example: Data Race
Show the timer race and fix.
