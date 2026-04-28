# Section 3: WaitGroups & Closures

> Lectures 9–13 — Concurrency in Go (Deepak Kumar Gunjetti, Udemy)

---

## Race Condition

- A **race condition** occurs when the order of execution is **not guaranteed**
- Concurrent programs do not execute in the order they are coded
- Compiler and CPU optimizations can reorder instructions
- The outcome depends on the non-deterministic ordering of goroutine execution

---

## The Problem with Goroutines

Consider this code:

```go
var data int

go func() { data++ }()

if data == 0 {
    fmt.Printf("the value of data is %v\n", data)
}
```

There are **three possible outputs** depending on scheduling:

| Output | Execution Sequence |
|---|---|
| Nothing printed | Goroutine runs first: `data++` → check `data == 0` is false |
| `"the value of data is 0"` | Main checks before goroutine runs: check → print → `data++` |
| `"the value of data is 1"` | Goroutine runs between check and print: check → `data++` → print |

This is a **data race** — the goroutine and main goroutine access `data` concurrently without synchronization.

---

## sync.WaitGroup

`sync.WaitGroup` deterministically blocks the main goroutine until all child goroutines have completed. It creates a **join point**.

### Methods

| Method | Description |
|---|---|
| `wg.Add(n)` | Increment counter — indicates `n` goroutines to wait for |
| `wg.Done()` | Decrement counter — indicates a goroutine is exiting (usually via `defer`) |
| `wg.Wait()` | Block until counter reaches zero |

### Join Point Diagram

```
main goroutine ──────┬─────────────────── wg.Wait() ──── continues ──
                     │                        ↑
                     └── go func() ── work ── wg.Done()
```

### Example

```go
var data int
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    data++
}()

wg.Wait()
fmt.Printf("the value of data is %v\n", data) // always prints 1
```

> See [`01-waitgroup/main.go`](01-waitgroup/main.go) for the full runnable example.

---

## Goroutines & Closures

- Goroutines execute within the **same address space** they are created in
- A closure can directly read and modify variables in its enclosing lexical block
- This is powerful but dangerous when combined with loops

---

## Common Closure Bug: Loop Variable Capture

### The Bug

```go
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(i) // captures the variable, not the value
    }()
}
wg.Wait()
```

**Problem:** All goroutines share the **same** `i` variable. By the time any goroutine executes, the loop has likely finished and `i == 3`. You may see output like:

```
3
3
3
```

> See [`02-closure/main.go`](02-closure/main.go) to reproduce the bug.

### The Fix — Pass as Argument

```go
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(val int) {
        defer wg.Done()
        fmt.Println(val) // each goroutine gets its own copy
    }(i)
}
wg.Wait()
```

Each goroutine receives its **own copy** of `i` via the function parameter. Output will contain `0`, `1`, and `2` (in any order).

> See [`03-closure-fix/main.go`](03-closure-fix/main.go) for the corrected version.

---

## Key Takeaways

1. Never assume goroutine execution order — use synchronization primitives
2. `sync.WaitGroup` is the simplest way to wait for goroutines to finish
3. Always `defer wg.Done()` to guarantee the counter is decremented
4. When launching goroutines inside loops, pass the loop variable as an argument to avoid the closure capture bug
