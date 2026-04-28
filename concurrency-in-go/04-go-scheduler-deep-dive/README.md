# Section 4: Go Scheduler Deep Dive

> Lectures 14–17 — Concurrency in Go (Deepak Kumar Gunjetti, Udemy)

---

## M:N Scheduler

- Go's scheduler is part of the **Go runtime**, known as an **M:N scheduler**
- Runs entirely in **user space** (not kernel space)
- **N** goroutines are scheduled on **M** OS threads running on at most **GOMAXPROCS** logical processors
- Go runtime creates worker OS threads equal to `GOMAXPROCS` (default = number of CPU cores)
- The scheduler distributes runnable goroutines across multiple worker OS threads

---

## Scheduler Components — The G-M-P Model

| Component | Name | Description |
|---|---|---|
| **G** | Goroutine | Contains scheduling info: stack, instruction pointer, channel blocking info |
| **M** | Machine | Represents an OS thread managed by the kernel |
| **P** | Processor | Logical processor that manages scheduling of goroutines (at most GOMAXPROCS) |
| **LRQ** | Local Run Queue | Each P has one; holds goroutines assigned to that processor |
| **GRQ** | Global Run Queue | Holds newly created goroutines not yet assigned to any P |

### Architecture Diagram

```
┌──────────────────────────────────────────────────────┐
│                      Go Runtime                      │
│                                                      │
│   ┌──────────────┐        ┌──────────────┐           │
│   │      P1      │        │      P2      │   ...     │
│   │  ┌────────┐  │        │  ┌────────┐  │           │
│   │  │   G1   │◄─┤        │  │   G5   │◄─┤           │
│   │  │(running)│  │        │  │(running)│  │           │
│   │  └────────┘  │        │  └────────┘  │           │
│   │              │        │              │           │
│   │  LRQ:        │        │  LRQ:        │           │
│   │  G2, G3, G4  │        │  G6, G7, G8  │           │
│   └──────┬───────┘        └──────┬───────┘           │
│          │                       │                   │
│   ┌──────┴───────┐        ┌──────┴───────┐           │
│   │      M1      │        │      M2      │           │
│   │  (OS Thread) │        │  (OS Thread) │           │
│   └──────┬───────┘        └──────┬───────┘           │
│          │                       │                   │
│          │    ┌──────────────┐   │                   │
│          │    │     GRQ      │   │                   │
│          │    │  G9, G10 ... │   │                   │
│          │    └──────────────┘   │                   │
└──────────┼───────────────────────┼───────────────────┘
      ┌────┴─────┐           ┌────┴─────┐
      │  Core 1  │           │  Core 2  │
      └──────────┘           └──────────┘
```

---

## Asynchronous Preemption

- As of **Go 1.14**, the scheduler implements **asynchronous preemption**
- Prevents long-running goroutines from hogging the CPU indefinitely
- A goroutine running for more than **10ms** is signaled to yield
- Before Go 1.14, preemption only happened at function call boundaries (cooperative preemption)

---

## Goroutine States

A goroutine transitions through three states:

```
              Preempted / Yield
                    │
                    ▼
  ┌──────────┐   ┌──────────┐
  │ Runnable │──▶│Executing │
  └──────────┘   └────┬─────┘
       ▲              │
       │        I/O or event
       │           wait
       │              │
       │              ▼
       │         ┌─────────┐
       └─────────│ Waiting │
      (I/O done) └─────────┘
```

| State | Description |
|---|---|
| **Runnable** | Ready to run, waiting in a run queue |
| **Executing** | Currently running on an OS thread (M) |
| **Waiting** | Blocked on I/O, channel, mutex, sleep, etc. |

---

## Context Switch: Synchronous System Call

When a goroutine makes a **synchronous system call** (e.g., file-based I/O), the OS thread blocks. The scheduler handles this transparently:

### Step 1 — Before the syscall

```
    P1
    ├── G1 (running on M1) ← about to make syscall
    └── LRQ: G2, G3
```

### Step 2 — G1 enters blocking syscall, M1 blocks

```
    G1 ←──── M1 (blocked in syscall)

    P1 detaches from M1
```

### Step 3 — Scheduler brings M2 from thread pool, attaches P1

```
    P1 ──── M2 (from thread pool)
    ├── G2 (now running)
    └── LRQ: G3

    G1 ←──── M1 (still blocked)
```

### Step 4 — Syscall returns

```
    G1 completes syscall
    G1 → placed back into P1's LRQ (or GRQ)
    M1 → returned to thread pool (put to sleep)
```

**Key insight:** The logical processor P is never idle — it detaches from the blocked thread and continues running goroutines on a new thread.

---

## Context Switch: Asynchronous System Call

For **network I/O**, Go uses the **netpoller** — an abstraction over OS async I/O facilities:

| OS | Interface |
|---|---|
| macOS | `kqueue` |
| Linux | `epoll` |
| Windows | `iocp` (I/O completion ports) |

### Step 1 — G1 makes a network call

```
    P1
    ├── G1 (running on M1) ← makes network I/O call
    └── LRQ: G2, G3
```

### Step 2 — G1 parked at netpoller, M1 stays with P1

```
    P1 ──── M1
    ├── G2 (now running)        Netpoller
    └── LRQ: G3                 └── G1 (waiting for fd ready)
```

Unlike synchronous syscalls, **M1 is NOT blocked**. The goroutine is parked at the netpoller while M1 continues executing other goroutines.

### Step 3 — File descriptor ready, G1 returns to run queue

```
    P1 ──── M1                  Netpoller
    ├── G2 (running)            └── (empty)
    └── LRQ: G3, G1
```

**Key insight:** The complexity of async I/O multiplexing is moved from the application into the Go runtime. You write straightforward blocking-style code; the runtime handles the rest.

---

## Work Stealing

When a processor (P) has an **empty local run queue**, it doesn't sit idle. The scheduler performs **work stealing** in this order:

1. **Check the Global Run Queue (GRQ)** for a runnable G
2. **Steal half the goroutines** from another P's Local Run Queue
3. **Check the netpoller** for goroutines whose I/O is ready

### Example

```
  Before work stealing:

    P1 ──── M1                    P2 ──── M2
    └── LRQ: (empty!)            └── LRQ: G5, G6, G7, G8

  After P1 steals from P2:

    P1 ──── M1                    P2 ──── M2
    └── LRQ: G7, G8              └── LRQ: G5, G6
```

Work stealing ensures **balanced distribution** of goroutines across all logical processors, leading to better CPU utilization and faster execution.

---

## Key Takeaways

1. Go's M:N scheduler maps N goroutines onto M OS threads across GOMAXPROCS logical processors
2. The **G-M-P model** (Goroutine, Machine, Processor) is the foundation of Go's scheduling
3. **Synchronous syscalls** cause M to block — scheduler detaches P and assigns a new M
4. **Asynchronous syscalls** (network I/O) use the **netpoller** — neither M nor P blocks
5. **Work stealing** keeps all processors busy by redistributing goroutines
6. **Asynchronous preemption** (Go 1.14+) ensures no goroutine can monopolize the CPU beyond ~10ms
7. The programmer writes simple, blocking-style code; the runtime handles the concurrency complexity
