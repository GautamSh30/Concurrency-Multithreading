# Section 1: Introduction to Concurrency

> **Course**: Concurrency in Go (Udemy) — Deepak Kumar Gunjetti
> **Lectures**: 1–4

---

## Overview of Concurrency

- Concurrency is about **multiple things happening at the same time in random order**
- Go provides built-in support for concurrency as a first-class citizen
- To run faster, applications need to be divided into **multiple independent units** and run in parallel

---

## Concurrency vs Parallelism

| Aspect | Concurrency | Parallelism |
|--------|-------------|-------------|
| Definition | Composition of independent execution computations, which may or may not run in parallel | Ability to execute multiple computations simultaneously |
| CPU Cores | Works on single-core (time-slicing) | Requires multi-core (truly simultaneous) |
| Focus | **Dealing** with lots of things at once | **Doing** lots of things at once |

> **Key insight**: Concurrency *enables* Parallelism. Concurrency is about structure; parallelism is about execution.

```
Concurrency (single core — time slicing):

  Core 0:  |--Task A--|--Task B--|--Task A--|--Task B--|-->
            ↑ tasks interleave on one core

Parallelism (multi-core — simultaneous):

  Core 0:  |--Task A--|--Task A--|--Task A--|-->
  Core 1:  |--Task B--|--Task B--|--Task B--|-->
            ↑ tasks truly run at the same time
```

---

## Processes

A **process** is an instance of a running program. The OS provides each process with its own isolated environment.

### Process Memory Layout

```
+---------------------------+
|         Stack             |  ← Local variables, function calls
|            ↓              |
|         (grows down)      |
|                           |
|         (grows up)        |
|            ↑              |
|          Heap             |  ← Dynamic memory allocation (malloc, new)
+---------------------------+
|          Data             |  ← Global / static variables
+---------------------------+
|          Code             |  ← Machine instructions (text segment)
+---------------------------+
```

- **Code**: Machine instructions of the compiled program
- **Data**: Global and static data
- **Heap**: Dynamically allocated memory
- **Stack**: Local variables, function call frames

---

## Threads

- The **smallest unit of execution** that the CPU accepts
- Every process has at least one thread — the **main thread**
- A process can create multiple threads
- Threads **share the same address space** (code, data, heap) but have their own stack
- Threads run **independent** of each other
- The **OS scheduler** makes scheduling decisions at the thread level

### Thread Lifecycle / States

```
                ┌──────────────┐
                │   Runnable   │
                └──────┬───────┘
                       │ scheduled
                       ▼
                ┌──────────────┐
         ┌──── │  Executing   │ ────┐
         │     └──────────────┘     │
         │ I/O request         preempted (time slice expired)
         ▼                          │
  ┌──────────────┐                  │
  │   Waiting    │                  │
  └──────┬───────┘                  │
         │ I/O complete             │
         └──────────┐    ┌─────────┘
                    ▼    ▼
                ┌──────────────┐
                │   Runnable   │  (back in the queue)
                └──────────────┘
```

### Thread vs Process Memory

```
Process
┌─────────────────────────────────────┐
│  Code  │  Data  │       Heap        │  ← shared by all threads
├────────┴────────┴───────────────────┤
│  Thread 1 Stack │  Thread 2 Stack   │  ← each thread gets its own
└─────────────────┴───────────────────┘
```

---

## Limitations of Threads

### 1. Context Switches are Expensive

When the OS switches between threads, it must save and restore:

**Thread context:**
- Program counter (current instruction)
- CPU registers
- Stack pointer

**Process context (when switching across processes):**
- Process state
- CPU scheduling information
- Memory management information (page tables)
- I/O status information

### 2. The C10K Problem

- The scheduler period is divided equally among all threads
- With a minimum time-slice of ~2ms:
  - **10 threads** → 20ms scheduler cycle (responsive)
  - **100 threads** → 200ms scheduler cycle
  - **1,000 threads** → 2s scheduler cycle
  - **10,000 threads** → 20s scheduler cycle (unacceptable!)

```
  10,000 threads × 2ms minimum time-slice = 20 second scheduler period
  ───────────────────────────────────────────────────────────────────
  Each thread gets CPU attention only once every 20 seconds!
```

### 3. Fixed Stack Size

- Threads are typically allocated a **fixed stack size of ~8MB**
- With 1,000 threads: 1,000 × 8MB = **8GB** of stack memory alone
- This limits the number of threads you can create

---

## Why Concurrency is Hard

### Shared Memory

Threads share the same address space. This is both powerful and dangerous — concurrent access to shared memory creates complexity.

### Data Race

A **data race** occurs when two or more threads access the same memory location concurrently, and at least one is a write.

```go
// DATA RACE EXAMPLE
var data int

// Thread 1          // Thread 2
data++               data++

// Expected: data == 2
// Possible: data == 1 (lost update!)
```

### Atomicity

Operations that appear to be single steps are often **not atomic**. For example, `i++` involves three steps:

```
i++  is actually:
  1. RETRIEVE value of i from memory
  2. INCREMENT the value
  3. STORE the new value back to memory

Thread 1:  RETRIEVE(0) → INCREMENT(1) → STORE(1)
Thread 2:       RETRIEVE(0) → INCREMENT(1) → STORE(1)

Result: i == 1 (expected 2!)
```

### Memory Access Synchronization (Mutex)

**Locks** (mutexes) provide exclusive access to shared memory, but they come with trade-offs:

```go
var mu sync.Mutex
var count int

// Thread 1           // Thread 2
mu.Lock()             mu.Lock()        // blocks until Thread 1 unlocks
count++               count++
mu.Unlock()           mu.Unlock()
```

- Locks reduce the degree of parallelism (critical sections are sequential)
- Developer must follow the convention to Lock/Unlock — **no compiler enforcement**
- Forgetting to unlock → deadlock; incorrect lock granularity → poor performance

### Deadlocks

A **deadlock** occurs when threads are waiting on each other in a circular chain:

```
Thread 1:  holds Lock A  →  waiting for Lock B
Thread 2:  holds Lock B  →  waiting for Lock A

Both threads are stuck forever!
```

**Four conditions for deadlock (Coffman conditions):**
1. **Mutual Exclusion** — resources cannot be shared
2. **Hold and Wait** — holding one resource while waiting for another
3. **No Preemption** — resources cannot be forcibly taken
4. **Circular Wait** — circular chain of dependencies

---

## Key Takeaways

1. **Concurrency ≠ Parallelism** — concurrency is about structure, parallelism is about execution
2. **OS threads are expensive** — context switches, fixed stack size, and the C10K problem limit scalability
3. **Shared memory is the root of concurrency bugs** — data races, deadlocks, and non-atomic operations
4. **Mutexes help but don't solve everything** — they reduce parallelism and rely on developer discipline
5. **Go was designed to address these problems** — with goroutines, channels, and a user-space scheduler
