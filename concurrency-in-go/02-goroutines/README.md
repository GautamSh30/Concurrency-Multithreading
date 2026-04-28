# Section 2: Goroutines

> **Course**: Concurrency in Go (Udemy) — Deepak Kumar Gunjetti
> **Lectures**: 5–8

---

## CSP — Communicating Sequential Processes

- Proposed by **Tony Hoare (1978)**
- Each process is built for **sequential execution**
- Data is **communicated between processes** — no shared memory
- Scale by **adding more of the same**

```
CSP Model:

  ┌───────────┐    channel    ┌───────────┐
  │ Process A │ ────────────► │ Process B │
  │(sequential│               │(sequential│
  │ execution)│ ◄──────────── │ execution)│
  └───────────┘    channel    └───────────┘

  No shared memory — data flows through channels
```

> Go's concurrency model is built on CSP: goroutines are the processes, channels are the communication mechanism.

---

## Go's Concurrency Tool Set

| Tool | Purpose |
|------|---------|
| `goroutines` | Lightweight concurrent execution |
| `channels` | Communication between goroutines |
| `select` | Multiplex channel operations |
| `sync` package | Low-level primitives (Mutex, WaitGroup, etc.) |

---

## What are Goroutines?

Goroutines are **user-space threads** managed by the **Go runtime**, not the OS.

- Extremely **lightweight** — starts with only **2KB** of stack (grows and shrinks dynamically as needed)
- **Low CPU overhead** — approximately three instructions per function call
- Can create **hundreds of thousands** of goroutines in the same address space
- **Channels** are used for communication — sharing memory is avoided
- **Context switching** is much cheaper than OS thread context switching (happens in user space)

```
OS Thread:                          Goroutine:
┌──────────────────────┐            ┌──────────────────────┐
│  Stack: 8MB (fixed)  │            │  Stack: 2KB (dynamic)│
│  Managed by: OS      │            │  Managed by: Go      │
│  Context switch:     │            │  Context switch:     │
│    Kernel space       │            │    User space         │
│    (expensive)       │            │    (cheap)           │
│  Creation: limited   │            │  Creation: hundreds  │
│                      │            │    of thousands      │
└──────────────────────┘            └──────────────────────┘
```

---

## Goroutines vs OS Threads

| Feature | Goroutine | OS Thread |
|---------|-----------|-----------|
| **Stack Size** | 2KB (dynamic, grows/shrinks) | ~8MB (fixed) |
| **Context Switch** | User space (cheap, ~200ns) | Kernel space (expensive, ~1-2μs) |
| **Creation** | Hundreds of thousands | Limited by OS resources |
| **Management** | Go runtime scheduler | OS kernel scheduler |
| **Identity** | No thread-local storage | Has thread ID |
| **Communication** | Channels (CSP model) | Shared memory + locks |

---

## How Goroutines Work

- The Go runtime creates a set of **worker OS threads**
- Goroutines are **multiplexed** onto these OS threads
- Many goroutines execute in the context of a **single OS thread**

```
Go Runtime Scheduler (M:N scheduling):

  Goroutines:    G1  G2  G3  G4  G5  G6  G7  G8  G9

                  ↓   ↓   ↓   ↓   ↓   ↓   ↓   ↓   ↓
                ┌─────────┐ ┌─────────┐ ┌─────────┐
  OS Threads:   │ Thread 1│ │ Thread 2│ │ Thread 3│
                └─────────┘ └─────────┘ └─────────┘
                      ↓           ↓           ↓
                ┌─────────┐ ┌─────────┐ ┌─────────┐
  CPU Cores:    │  Core 0 │ │  Core 1 │ │  Core 2 │
                └─────────┘ └─────────┘ └─────────┘

  Many goroutines → few OS threads → available CPU cores
```

- When a goroutine blocks (e.g., on I/O), the Go scheduler moves other goroutines to a different OS thread
- This is called **M:N scheduling** — M goroutines on N OS threads

---

## Starting a Goroutine

There are three ways to start a goroutine using the `go` keyword:

### 1. Direct Function Call

```go
go functionName(arg1, arg2)
```

### 2. Anonymous Function (Closure)

```go
go func(s string) {
    fmt.Println(s)
}("hello")
```

### 3. Function Value

```go
f := functionName
go f(arg1, arg2)
```

> **Important**: When `main()` returns, the program exits — it does NOT wait for goroutines to finish. Use `sync.WaitGroup` or channels to coordinate goroutine completion.

---

## Examples

### Example 1: Hello Goroutines (`01-hello/`)

Demonstrates all three ways of launching goroutines:

```go
// Direct function call
go fun("direct goroutine")

// Anonymous function
go func(s string) {
    for i := 0; i < 3; i++ {
        fmt.Println(s)
        time.Sleep(100 * time.Millisecond)
    }
}("anonymous goroutine")

// Function value
fv := fun
go fv("function value goroutine")
```

Run: `go run main.go`

### Example 2: Concurrent TCP Server (`02-client-server/`)

Demonstrates goroutines for handling multiple client connections concurrently:

```go
for {
    conn, err := li.Accept()
    if err != nil {
        log.Fatal(err)
    }
    go handleConn(conn)  // each client handled in its own goroutine
}
```

Run:
```bash
# Terminal 1: Start server
cd server && go run main.go

# Terminal 2: Connect client
cd client && go run main.go

# Terminal 3: Connect another client
cd client && go run main.go
```

### Example 3: Sequential vs Concurrent Addition (`03-add/`)

Benchmarks sequential vs concurrent summation of 10 million numbers:

```go
// Sequential — single goroutine sums all numbers
func Add(numbers []int) int64 { ... }

// Concurrent — splits work across NumCPU goroutines
func AddConcurrent(numbers []int) int64 { ... }
```

Key techniques used:
- `runtime.NumCPU()` — detect available cores
- `sync.WaitGroup` — wait for all goroutines to complete
- `atomic.AddInt64()` — safely accumulate partial sums

Run benchmarks:
```bash
cd counting && go test -bench .
```

---

## Key Takeaways

1. **Goroutines are not threads** — they are much lighter (2KB vs 8MB stack) and managed by the Go runtime
2. **Use the `go` keyword** to launch a goroutine — it's that simple
3. **`main()` won't wait** for goroutines — use `sync.WaitGroup` or channels for synchronization
4. **M:N scheduling** — the Go runtime multiplexes many goroutines onto few OS threads
5. **CSP model** — prefer communicating via channels over sharing memory
6. **Concurrency can improve performance** — the `03-add` benchmark demonstrates real speedup from parallelizing work
