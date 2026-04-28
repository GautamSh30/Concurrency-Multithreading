# Section 5: Channels

> Lectures 18–32 — Concurrency in Go (Deepak Kumar Gunjetti)

## What are Channels?

- Channels are the **pipes** that connect concurrent goroutines.
- They allow you to **communicate data** between goroutines in a type-safe, thread-safe manner.
- They **synchronize** goroutines — a send/receive is a point of coordination.
- Go's motto: *"Do not communicate by sharing memory; instead, share memory by communicating."*

---

## Declaring and Initializing

```go
var ch chan T         // declare (nil channel — zero value)
ch = make(chan T)     // unbuffered channel
ch := make(chan T, n) // buffered channel with capacity n
```

- A channel declared but not initialized with `make` is `nil`.
- Channels are reference types — `make()` returns a pointer to the underlying `hchan` struct on the heap.

---

## Sending and Receiving

```go
ch <- value   // send value into channel
v = <-ch      // receive value from channel
v, ok = <-ch  // receive with close check (ok == false means channel is closed and drained)
```

---

## Channels are Blocking

- **Receive** `<-ch`: the goroutine blocks (waits) until another goroutine sends a value.
- **Send** `ch <- value`: the goroutine blocks until another goroutine is ready to receive.

This blocking behavior is what makes channels a synchronization mechanism.

---

## Closing Channels

```go
close(ch)
```

- Signals that **no more values** will be sent on the channel.
- Receiving from a closed channel returns the **zero value** immediately.
- `for value := range ch` iterates over values until the channel is closed and drained.
- **Only the sender should close** a channel, never the receiver.
- Closing an already-closed channel **panics**.
- Closing a nil channel **panics**.

---

## Unbuffered vs Buffered Channels

| Property | Unbuffered `make(chan T)` | Buffered `make(chan T, n)` |
|---|---|---|
| Synchronization | Synchronous — sender blocks until receiver ready | Asynchronous — sender blocks only when buffer full |
| Internal storage | None | In-memory FIFO circular queue |
| Capacity | 0 | n |
| Use case | Guaranteed hand-off | Decouple producer/consumer speed |

---

## Channel Direction (in function parameters)

```go
func send(ch chan<- string) {} // send-only channel
func recv(ch <-chan string) {} // receive-only channel
```

- Restricting direction in function signatures **increases type-safety**.
- The compiler enforces direction at compile time.
- A bidirectional channel is implicitly convertible to a directional one, but not vice versa.

---

## Nil Channel Behavior

| Operation | Nil Channel |
|---|---|
| Read `<-ch` | Blocks forever |
| Write `ch <- v` | Blocks forever |
| Close `close(ch)` | **Panics** |

Always initialize channels before use!

---

## Channel Ownership Pattern

Assign channel **ownership** to a single goroutine to avoid hazards:

- **Owner goroutine**: instantiates the channel, writes to it, and closes it.
- **Consumer goroutine**: only receives (has a read-only `<-chan` view).

This pattern avoids:
1. Writing to a nil channel (blocks forever)
2. Writing to a closed channel (panics)
3. Closing a nil channel (panics)
4. Closing a channel more than once (panics)

```go
func owner() <-chan int {
    ch := make(chan int, 5)
    go func() {
        defer close(ch)
        for i := 0; i < 5; i++ {
            ch <- i
        }
    }()
    return ch // consumer only gets read-only access
}
```

---

## Deep Dive: The `hchan` Struct

When you call `make(chan T)`, Go allocates an `hchan` struct on the **heap** and returns a pointer.

Key fields of `hchan`:

| Field | Purpose |
|---|---|
| `buf` | Circular ring buffer (for buffered channels) |
| `sendx` | Send index into the ring buffer |
| `recvx` | Receive index into the ring buffer |
| `lock` | Mutex protecting the struct |
| `sendq` | Wait queue of goroutines blocked on send (`sudog` linked list) |
| `recvq` | Wait queue of goroutines blocked on receive (`sudog` linked list) |

A `sudog` struct represents a goroutine waiting on a channel. It stores a pointer to the goroutine (`g`) and the element being sent/received (`elem`).

---

## Buffer Full Scenario (Buffered Channel)

1. **G1** fills the buffer, tries to send one more → G1 is **blocked**.
2. G1 is parked in the `sendq` wait queue. Its value is stored in `sudog.elem`. Runtime calls `gopark()` — G1 is moved off its processor (P).
3. **G2** comes to receive: dequeues a value from the ring buffer into its variable.
4. G2 pops G1 from `sendq`, enqueues G1's `sudog.elem` value into the buffer.
5. G2 calls `goready(G1)` → G1 is placed back on the **Local Run Queue (LRQ)**.

---

## Buffer Empty Scenario (Buffered Channel)

1. **G2** tries to receive from an empty buffer → G2 is **blocked**.
2. G2 is parked in the `recvq` wait queue. The address of G2's target variable is stored in `sudog.elem`. Runtime calls `gopark()`.
3. **G1** comes to send: finds G2 waiting in `recvq`.
4. G1 writes the value **directly into G2's stack variable** (no buffer copy needed!).
5. G1 calls `goready(G2)` → G2 is placed back on the LRQ.

This direct write optimization avoids the extra copy through the buffer.

---

## Unbuffered Channel Internals

An unbuffered channel has **no ring buffer** — it's a pure synchronization point.

- **Send**: if a receiver is waiting in `recvq` → write the value directly into the receiver's stack variable. Otherwise, the sender is parked in `sendq`.
- **Receive**: if a sender is waiting in `sendq` → copy the value from the sender's `sudog.elem`. Otherwise, the receiver is parked in `recvq`.

---

## Key Takeaways

1. Channels are Go's primary mechanism for goroutine communication and synchronization.
2. Unbuffered channels guarantee a hand-off; buffered channels decouple sender and receiver timing.
3. Use channel direction in function signatures to enforce correct usage at compile time.
4. Follow the **ownership pattern**: one goroutine owns (writes + closes), consumers only read.
5. Understand the `hchan` internals — ring buffer, wait queues, and the direct-write optimization — to reason about performance and behavior.
6. Never send on a closed channel, never close a nil channel, and always initialize channels before use.
