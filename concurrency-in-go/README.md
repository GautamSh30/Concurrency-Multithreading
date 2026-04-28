# Concurrency in Go - Complete Course Notes

> Comprehensive notes from the Udemy course **"Concurrency in Go (Golang)"** by Deepak Kumar Gunjetti.
> Each section is self-contained with theory notes, working Go code, and practical exercises.

## Course Structure

| # | Section | Topics | Key Concepts |
|---|---------|--------|--------------|
| 01 | [Introduction](./01-introduction/) | Overview, Processes, Threads | Concurrency vs Parallelism, C10K Problem |
| 02 | [Goroutines](./02-goroutines/) | CSP, Goroutine basics | Lightweight threads, 2KB stack |
| 03 | [WaitGroups & Closures](./03-waitgroups-and-closures/) | sync.WaitGroup, Closures | Join points, variable capture |
| 04 | [Go Scheduler Deep Dive](./04-go-scheduler-deep-dive/) | M:N Scheduler, G-M-P model | Work stealing, netpoller |
| 05 | [Channels](./05-channels/) | Buffered, Unbuffered, Direction | hchan internals, ownership |
| 06 | [Select](./06-select/) | Multiplexing, Timeouts | Non-blocking, time.After |
| 07 | [Mutex & Atomic](./07-mutex-and-atomic/) | sync.Mutex, sync.RWMutex | Atomic operations |
| 08 | [Sync Primitives](./08-sync-primitives/) | Cond, Once, Pool | Signal, Broadcast, object reuse |
| 09 | [Race Detector](./09-race-detector/) | go run -race | Finding data races |
| 10 | [Web Crawler](./10-web-crawler/) | Sequential vs Concurrent | Real-world concurrency |
| 11 | [Pipelines & Fan-out/Fan-in](./11-pipelines-and-fanout-fanin/) | Patterns, Cancellation | Done channel, goroutine leaks |
| 12 | [Context Package](./12-context-package/) | WithCancel, Deadline, Timeout | Request-scoped data, HTTP timeouts |
| 13 | [Interfaces](./13-interfaces/) | Implicit satisfaction, Type assertion | io.Writer, Stringer, empty interface |

## Quick Start

Each section is a standalone Go module. To run any example:

```bash
cd 02-goroutines/01-hello
go run main.go
```

To run benchmarks (e.g., sequential vs concurrent addition):

```bash
cd 02-goroutines/03-add
go test -bench=. ./counting/
```

To detect race conditions:

```bash
cd 09-race-detector/01-race-example
go run -race main.go
```

## Prerequisites

- Go 1.21 or later
- Basic Go knowledge (variables, functions, structs)

## Go Concurrency Mental Model

```
          ┌─────────────────────────────────────────┐
          │           Go Runtime                     │
          │                                          │
          │  ┌──────┐   ┌──────┐   ┌──────┐        │
          │  │  P1  │   │  P2  │   │  Pn  │  ...   │
          │  │ ┌──┐ │   │ ┌──┐ │   │ ┌──┐ │        │
          │  │ │G1│ │   │ │G4│ │   │ │G7│ │        │
          │  │ └──┘ │   │ └──┘ │   │ └──┘ │        │
          │  │LRQ:  │   │LRQ:  │   │LRQ:  │        │
          │  │G2,G3 │   │G5,G6 │   │G8,G9 │        │
          │  └──┬───┘   └──┬───┘   └──┬───┘        │
          │     │          │          │              │
          │  ┌──┴───┐   ┌──┴───┐   ┌──┴───┐        │
          │  │  M1  │   │  M2  │   │  Mn  │        │
          │  └──┬───┘   └──┴───┘   └──┬───┘        │
          └─────┼──────────┼──────────┼─────────────┘
             ┌──┴──┐   ┌──┴──┐   ┌──┴──┐
             │Core1│   │Core2│   │CoreN│
             └─────┘   └─────┘   └─────┘

G = Goroutine, M = OS Thread, P = Logical Processor
LRQ = Local Run Queue
```

## Key Principles

1. **"Do not communicate by sharing memory; instead, share memory by communicating."**
2. Goroutines are cheap - use them liberally
3. Channels orchestrate; mutexes serialize
4. Always clean up goroutines (avoid leaks)
5. Use `context` for cancellation propagation

## Bonus

- [Go Concurrency Cheatsheet](./CHEATSHEET.md) - Quick reference card
